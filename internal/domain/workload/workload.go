package workload

import pkgworkload "github.com/wernsiet/morchy/pkg/workload"

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
	Container pkgworkload.Container
}
