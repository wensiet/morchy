package workload

type WorkloadStatus string

const (
	NewWorkloadStatus        WorkloadStatus = "new"
	ActiveWorkloadStatus     WorkloadStatus = "active"
	TerminatedWorkloadStatus WorkloadStatus = "terminated"
)

type Workload struct {
	ID     string
	Status WorkloadStatus
}
