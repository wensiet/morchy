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
