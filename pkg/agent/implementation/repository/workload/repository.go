package workload

import (
	"github.com/wernsiet/morchy/pkg/agent/domain/workload"
	"github.com/wernsiet/morchy/pkg/runtime"
)

type Repository struct {
	workloadStorage map[string]*workload.Workload
	limitsStorage   *runtime.ResourceLimits
}

func NewRepository() *Repository {
	return &Repository{
		workloadStorage: make(map[string]*workload.Workload),
		limitsStorage: &runtime.ResourceLimits{
			CPU: 0,
			RAM: 0,
		},
	}
}
