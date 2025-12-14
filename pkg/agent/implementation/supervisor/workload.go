package supervisor

import (
	"github.com/wernsiet/morchy/pkg/agent/domain/workload"
	"github.com/wernsiet/morchy/pkg/agent/infrastructure"
)

type WorkloadSupervisor = infrastructure.Supervisor[workload.Workload]

func NewSupervisor() *WorkloadSupervisor {
	return infrastructure.NewSupervisor[workload.Workload]()
}
