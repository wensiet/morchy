package workload

import (
	"github.com/wernsiet/morchy/pkg/agent/domain"
	"github.com/wernsiet/morchy/pkg/agent/domain/workload"
	"github.com/wernsiet/morchy/pkg/runtime"
)

func (r *Repository) GetWorkload(id string) (*workload.Workload, error) {
	domainWorkload, ok := r.workloadStorage[id]
	if !ok {
		return nil, domain.ErrorWorkloadNotFound.Errorf("workload with id=%s not found", id)
	}
	return domainWorkload, nil
}

func (r *Repository) SaveWorklod(w workload.Workload) (*workload.Workload, error) {
	_, ok := r.workloadStorage[w.ID]
	if ok {
		return nil, domain.ErrorWorkloadNotFound.Errorf("workload with id=%s already exists", w.ID)
	}
	r.workloadStorage[w.ID] = &w
	return &w, nil
}

func (r *Repository) ListWorkloads() []*workload.Workload {
	var workloads []*workload.Workload
	for _, w := range r.workloadStorage {
		workloads = append(workloads, w)
	}
	return workloads
}

func (r *Repository) RemoveWorkload(id string) {
	delete(r.workloadStorage, id)
}

func (r *Repository) GetResourceLimits() *runtime.ResourceLimits {
	return r.limitsStorage
}

func (r *Repository) SetResourceLimits(limits runtime.ResourceLimits) {
	r.limitsStorage = &limits
}

func (r *Repository) DecreaseResourceLimits(limits runtime.ResourceLimits) error {
	if r.limitsStorage.CPU < limits.CPU || r.limitsStorage.RAM < limits.RAM {
		return domain.ErrorInsufficientResources.Errorf("insufficient resources: requested CPU=%d, RAM=%d; available CPU=%d, RAM=%d",
			limits.CPU, limits.RAM, r.limitsStorage.CPU, r.limitsStorage.RAM)
	}
	r.limitsStorage.CPU -= limits.CPU
	r.limitsStorage.RAM -= limits.RAM
	return nil
}

func (r *Repository) IncreaseResourceLimits(limits runtime.ResourceLimits) {
	r.limitsStorage.CPU += limits.CPU
	r.limitsStorage.RAM += limits.RAM
}
