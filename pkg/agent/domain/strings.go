package domain

type EventActionStatus string

const (
	EventActionStatusSuccess EventActionStatus = "success"
	EventActionStatusFailed  EventActionStatus = "failed"
)

const (
	SAppName = "morchy"

	SCodeNotFound        = "not_found"
	SInternalServerError = "internal_server_error"

	STerminatedOnControlPlane  = "terminated_on_control_plane"
	SInsufficientResources     = "insufficient_resources"
	SHealthcheckFailed         = "healthcheck_failed"
	SNoWorkloadsToSchedule     = "no_workloads_to_schedule"
	SWorkloadAlreadyInStorage = "workload_already_in_storage"

	SOwnedByAnotherNode = "owned_by_other_node"

	SApplyWorkloadJoin = "ApplyWorkloadJoin"
	SLoadCurrentState  = "LoadCurrentState"

	SUsecase     = "usecase"
	SDomain      = "domain"
	SReason      = "reason"
	SManager     = "manager"
	SDuration    = "duration"
	SHealthcheck = "healthcheck"

	SHealthy = "healthy"
	SRunning = "running"

	SWorkload       = "workload"
	SWorkloadID     = "workload_id"
	SWorkloadStatus = "workload_status"
)
