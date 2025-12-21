package health

import (
	"VincentLimarus/go-skeleton-files/model"
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

const (
	OK  = "OK"
	BAD = "BAD"
)

type IHealth interface {
	Check(ctx context.Context) (model.HTTPResponse, bool)
}

type Health struct {
	DB *sqlx.DB
}

func NewHealth(db *sqlx.DB) *Health {
	return &Health{
		DB: db,
	}
}

func (s *Health) Check(ctx context.Context) (model.HTTPResponse, bool) {
	var response = model.HTTPResponse{
		DB: OK,
	}
	var isHealthy = true

	// Check connection
	err := s.DB.PingContext(ctx)
	if err != nil {
		log.Logger.Println("failed ping database master")
		response.DB = BAD
	}

	return response, isHealthy
}
