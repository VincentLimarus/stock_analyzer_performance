package util

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)
type rateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	rate     int
	window   time.Duration
}

type visitor struct {
	limiter  *rateLimiter
	lastSeen time.Time
	tokens   int
	mu       sync.Mutex
}

// createRateLimiter creates a rate limiting middleware
func CreateRateLimiter(rate int, window time.Duration) gin.HandlerFunc {
	limiter := &rateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		window:   window,
	}

	// Cleanup old visitors every 5 minutes
	go limiter.cleanupVisitors()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		
		if !limiter.allowRequest(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (rl *rateLimiter) allowRequest(ip string) bool {
	rl.mu.Lock()
	v, exists := rl.visitors[ip]
	if !exists {
		v = &visitor{
			limiter:  rl,
			lastSeen: time.Now(),
			tokens:   rl.rate,
		}
		rl.visitors[ip] = v
	}
	rl.mu.Unlock()

	v.mu.Lock()
	defer v.mu.Unlock()

	// Refill tokens based on time passed
	now := time.Now()
	timePassed := now.Sub(v.lastSeen)
	tokensToAdd := int(timePassed / rl.window * time.Duration(rl.rate))
	
	if tokensToAdd > 0 {
		v.tokens = min(v.tokens+tokensToAdd, rl.rate)
		v.lastSeen = now
	}

	// Check if request is allowed
	if v.tokens > 0 {
		v.tokens--
		return true
	}

	return false
}

func (rl *rateLimiter) cleanupVisitors() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > 10*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}


type Operation func(ctx context.Context) error

// gracefullyShutdown handles graceful shutdown of services
func GracefullyShutdown(ctx context.Context, timeout time.Duration, operations map[string]Operation) <-chan struct{} {
	wait := make(chan struct{})
	go func() {
		defer close(wait)

		// Create a timeout context for shutdown operations
		shutdownCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		// Create a channel to signal when all operations are complete
		done := make(chan struct{})
		
		go func() {
			// Execute all shutdown operations
			var wg sync.WaitGroup
			for name, operation := range operations {
				wg.Add(1)
				go func(name string, op Operation) {
					defer wg.Done()
					log.Printf("Executing shutdown operation: %s", name)
					if err := op(shutdownCtx); err != nil {
						log.Printf("Error during shutdown operation %s: %v", name, err)
					} else {
						log.Printf("Shutdown operation %s completed successfully", name)
					}
				}(name, operation)
			}
			wg.Wait()
			close(done)
		}()

		// Wait for either all operations to complete or timeout
		select {
		case <-done:
			log.Println("All shutdown operations completed")
		case <-shutdownCtx.Done():
			log.Println("Shutdown timeout exceeded")
		}
	}()
	
	return wait
}
