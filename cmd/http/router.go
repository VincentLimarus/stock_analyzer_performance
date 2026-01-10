package http

import (
	delivery "VincentLimarus/stock-analyzer-performance/delivery/http"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Router interface {
	Register() *gin.Engine
}

type router struct {
	engine            *gin.Engine
	delivery          delivery.IRegistry
	middlewareLimiter gin.HandlerFunc
}

func NewRouter(
	delivery delivery.IRegistry,
	middlewareLimiter gin.HandlerFunc,
) Router {
	return &router{
		engine:            gin.Default(),
		delivery:          delivery,
		middlewareLimiter: middlewareLimiter,
	}
}

func (r *router) Register() *gin.Engine {
	r.engine.Use(
		cors.New(cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: false,
			MaxAge:           12 * time.Hour,
		}),
		gin.Recovery(),
		gin.Logger(),
	)
	
	r.engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := r.engine.Group("/v1")
	r.Trades(v1)

	return r.engine
}

func (r *router) Trades(versionGroup *gin.RouterGroup) {
	tradesGroup := versionGroup.Group("/trades")

	tradesGroup.POST("/upsert", r.delivery.GetTrade().UpsertTradeData)
	tradesGroup.GET("", r.delivery.GetTrade().GetListTradeByYears) 
}