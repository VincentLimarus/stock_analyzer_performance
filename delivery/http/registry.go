package http

import (
	trade "VincentLimarus/stock-analyzer-performance/delivery/http/trades"
)

type IRegistry interface {
	GetTrade() trade.ITrade
}

type Registry struct {
	trade trade.ITrade
}

func NewRegistry(
	trade trade.ITrade,
) IRegistry {
	return &Registry{
		trade: trade,
	}
}

func (r *Registry) GetTrade() trade.ITrade {
	return r.trade
}
