package workload

import "github.com/wernsiet/morchy/pkg/runtime"

type WorkloadStatus string

const (
	NewWorkloadStatus        WorkloadStatus = "new"
	ActiveWorkloadStatus     WorkloadStatus = "active"
	TerminatedWorkloadStatus WorkloadStatus = "terminated"
)

type Workload struct {
	ID     string
	Status WorkloadStatus
	Spec   WorkloadSpec
}

type WorkloadSpec struct {
	Container runtime.Container
}
