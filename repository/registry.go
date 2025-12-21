package repository

type IRegistry interface {
	GetTx() *TransactionRunner
}

type Registry struct {
	masterTx *TransactionRunner
}

func NewRegistry(
	masterTx *TransactionRunner,
) *Registry {
	return &Registry{
		masterTx: masterTx,
	}
}

func (r *Registry) GetTx() *TransactionRunner {
	return r.masterTx
}