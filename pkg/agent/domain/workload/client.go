package workload

import "context"

type Client interface {
	CreateContainer(context.Context, Workload) (string, error)
	StartContainer(string) error
}
