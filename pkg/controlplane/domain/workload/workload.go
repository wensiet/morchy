package workload

type WorkloadStatus string

const (
	NewWorkloadStatus      WorkloadStatus = "new"
	PendingWorkloadStatus  WorkloadStatus = "pending"
	StuckWorkloadStatus    WorkloadStatus = "stuck"
	ActiveWorkloadStatus   WorkloadStatus = "active"
	FailedWorkloadStatus   WorkloadStatus = "failed"
	DegradedWorkloadStatus WorkloadStatus = "degraded"
)

type Workload struct {
	ID     string
	Status WorkloadStatus
	Spec   WorkloadSpec
}

type WorkloadSpec struct {
	Name    string
	Image   string
	CPU     uint
	RAM     uint
	Command []string
	Env     map[string]string
}
