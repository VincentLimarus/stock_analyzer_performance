package service

import "VincentLimarus/stock-analyzer-performance/service/health"

type IRegistry interface {
	GetHealth() health.IHealth
}

type Registry struct {
	health             health.IHealth
}

func NewRegistry(
	health health.IHealth,
) *Registry {
	return &Registry{
		health: health,
	}
}

func (r *Registry) GetHealth() health.IHealth {
	return r.health
}