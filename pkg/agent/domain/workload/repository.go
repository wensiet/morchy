package workload

import "github.com/wernsiet/morchy/pkg/runtime"

type Repository interface {
	GetWorkload(id string) (*Workload, error)
	ListWorkloads() []*Workload
	SaveWorkload(w Workload) (*Workload, error)
	RemoveWorkload(id string)

	GetResourceLimits() *runtime.ResourceLimits
	SetResourceLimits(limits runtime.ResourceLimits)
	DecreaseResourceLimits(limits runtime.ResourceLimits) error
	IncreaseResourceLimits(limits runtime.ResourceLimits)
}
