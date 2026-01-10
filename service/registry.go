package service

import (
	"VincentLimarus/stock-analyzer-performance/service/health"
	"VincentLimarus/stock-analyzer-performance/service/trades"
)

type IRegistry interface {
	GetHealth() health.IHealth
	GetTrade() trades.ITrade
}

type Registry struct {
	health       health.IHealth
	tradeService trades.ITrade
}

func NewRegistry(
	health health.IHealth,
	tradeService trades.ITrade,
) *Registry {
	return &Registry{
		health:       health,
		tradeService: tradeService,
	}
}

func (r *Registry) GetHealth() health.IHealth {
	return r.health
}

func (r *Registry) GetTrade() trades.ITrade {
	return r.tradeService
}
