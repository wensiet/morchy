package infrastructure

import (
	"context"
	"sync"
	"time"
)

type TaskID string

type PeriodicTask[T any] struct {
	ID       TaskID
	Interval time.Duration
	Input    T
	Execute  func(context.Context, T) error
	OnError  func(error)
}

type TaskSupervisor[T any] interface {
	Start(ctx context.Context, task PeriodicTask[T])
	Stop(id TaskID)
}

type Supervisor[T any] struct {
	mu     sync.Mutex
	cancel map[TaskID]context.CancelFunc
}

func NewSupervisor[T any]() *Supervisor[T] {
	return &Supervisor[T]{cancel: make(map[TaskID]context.CancelFunc)}
}

func (s *Supervisor[T]) Start(ctx context.Context, task PeriodicTask[T]) {
	s.mu.Lock()
	if _, ok := s.cancel[task.ID]; ok {
		s.mu.Unlock()
		return
	}

	taskCtx, cancel := context.WithCancel(ctx)
	s.cancel[task.ID] = cancel
	s.mu.Unlock()

	go func() {
		ticker := time.NewTicker(task.Interval)
		defer ticker.Stop()

		for {
			select {
			case <-taskCtx.Done():
				return
			case <-ticker.C:
				if err := task.Execute(taskCtx, task.Input); err != nil && task.OnError != nil {
					task.OnError(err)
				}
			}
		}
	}()
}

func (s *Supervisor[T]) Stop(id TaskID) {
	s.mu.Lock()
	if cancel, ok := s.cancel[id]; ok {
		cancel()
		delete(s.cancel, id)
	}
	s.mu.Unlock()
}
