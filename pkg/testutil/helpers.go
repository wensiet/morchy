package testutil

import (
	"context"
	"time"

	"github.com/wernsiet/morchy/pkg/controlplane/domain/workload"
	"github.com/wernsiet/morchy/pkg/runtime"
	"go.uber.org/zap"
)

func TestContext() context.Context {
	return context.Background()
}

func TestLogger() *zap.Logger {
	return zap.NewNop()
}

func NewTestWorkload(id string) *workload.Workload {
	cpu := uint(100)
	ram := uint(256)
	return &workload.Workload{
		ID:     id,
		Status: workload.NewWorkloadStatus,
		Spec: workload.WorkloadSpec{
			Name:    "test-workload",
			Image:   "nginx:latest",
			CPU:     cpu,
			RAM:     ram,
			Command: []string{"nginx", "-g", "daemon off;"},
			Env:     map[string]string{"ENV1": "value1", "ENV2": "value2"},
		},
	}
}

func NewTestWorkloadSpec() workload.WorkloadSpec {
	cpu := uint(100)
	ram := uint(256)
	return workload.WorkloadSpec{
		Name:    "test-workload",
		Image:   "nginx:latest",
		CPU:     cpu,
		RAM:     ram,
		Command: []string{"nginx", "-g", "daemon off;"},
		Env:     map[string]string{"ENV1": "value1", "ENV2": "value2"},
	}
}

func NewTestResourceLimits() *runtime.ResourceLimits {
	return &runtime.ResourceLimits{
		CPU: 1000,
		RAM: 4096,
	}
}

func NewTestLease(nodeID, workloadID string) *workload.Lease {
	return &workload.Lease{
		NodeID:     nodeID,
		WorkloadID: workloadID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

func NewTestEvent(id, sourceID, nodeID string, payload []byte) *workload.Event {
	return &workload.Event{
		ID:         id,
		SourceID:   sourceID,
		NodeID:     nodeID,
		Payload:    payload,
		ProducedAt: time.Now(),
		CreatedAt:  time.Now(),
	}
}

func NewTestEdge(upstreamAddress, proxyPath string) *workload.Edge {
	return &workload.Edge{
		UpstreamAddress: upstreamAddress,
		ProxyPath:       proxyPath,
	}
}
