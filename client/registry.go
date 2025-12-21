package client

type IRegistry interface {
}

type Registry struct {
}

func NewRegistry() IRegistry {
	return &Registry{}
}