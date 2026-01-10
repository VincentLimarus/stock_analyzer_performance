package repository

import (
	"VincentLimarus/stock-analyzer-performance/repository/trades"
)

type IRegistry interface {
	GetTx() *TransactionRunner
	GetTrade() trades.ITrade
}

type Registry struct {
	masterTx *TransactionRunner
	Trade    trades.ITrade
}

func NewRegistry(
	masterTx *TransactionRunner,
	Trade trades.ITrade,
) *Registry {
	return &Registry{
		masterTx: masterTx,
		Trade:    Trade,
	}
}

func (r *Registry) GetTx() *TransactionRunner {
	return r.masterTx
}

func (r *Registry) GetTrade() trades.ITrade {
	return r.Trade
}
