package repository

import (
	"sync"

	"github.com/wernsiet/morchy/pkg/edge/domain"
)

type ProxyPath string

type Repository struct {
	edgeStorage map[ProxyPath]*domain.Edge
	mu          sync.RWMutex
}

func NewRepository() *Repository {
	return &Repository{
		edgeStorage: make(map[ProxyPath]*domain.Edge),
	}
}
