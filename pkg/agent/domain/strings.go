package domain

const (
	SCodeNotFound        = "not_found"
	SInternalServerError = "internal_server_error"

	STerminatedOnControlPlane = "terminated_on_control_plane"
	SInsufficientResources    = "insufficient_resources"
	SHealthcheckFailed        = "healthcheck_failed"
	SNotWorkloadsToSchedule   = "no_workloads_to_schedule"

	SOwnedByAnotherNode = "owned_by_other_node"

	SApplyWorkloadJoin = "ApplyWorkloadJoin"

	SUsecase = "usecase"
	SDomain  = "domain"
	SReason  = "reason"
	SHealthy = "healthy"
	SRunning = "running"

	SWorkload = "workload"
)
