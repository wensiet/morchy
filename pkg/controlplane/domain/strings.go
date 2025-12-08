package domain

type ErrorCode string

const (
	InternalServerError ErrorCode = "internal_server_error"
	NotFound            ErrorCode = "not_found"

	SDomain  = "domain"
	SUnknown = "unknown"

	SWorkload   = "workload"
	SWorkloadID = "worklod_id"

	SNode   = "node"
	SNodeID = "node_id"

	SContainer     = "container"
	SContainerName = "container_name"
)
