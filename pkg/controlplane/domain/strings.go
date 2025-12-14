package domain

type ErrorCode string

const (
	InternalServerError ErrorCode = "internal_server_error"
	NotFound            ErrorCode = "not_found"
	BadRequest          ErrorCode = "bad_request"
	Conflict            ErrorCode = "conflict"

	SOwnedByAnotherNode = "owned_by_other_node"

	SDomain     = "domain"
	SUnknown    = "unknown"
	SReason     = "reason"
	SValidation = "validation"

	SWorkload   = "workload"
	SWorkloadID = "workload_id"

	SEvent         = "event"
	SEventSourceID = "event_source_id"

	SNode   = "node"
	SNodeID = "node_id"

	SContainer     = "container"
	SContainerName = "container_name"
)
