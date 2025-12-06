package infrastructure

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
)

type BackgroundTask struct {
	name       string
	executable func(context.Context) error
	interval   time.Duration
}

type BackgroundTaskRunner struct {
	logger *zap.Logger
	tasks  []BackgroundTask
}

func NewBackgroundTaskRunner(logger *zap.Logger) *BackgroundTaskRunner {
	return &BackgroundTaskRunner{
		logger: logger,
	}
}

func (b *BackgroundTaskRunner) RegisterTask(name string, executable func(context.Context) error, interval time.Duration) {
	b.tasks = append(b.tasks, BackgroundTask{
		name:       name,
		executable: executable,
		interval:   interval,
	})
}

func (b *BackgroundTaskRunner) Start(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(len(b.tasks))

	for i := range b.tasks {
		task := b.tasks[i]
		go func() {
			defer wg.Done()
			b.runTask(ctx, task)
		}()
	}

	wg.Wait()
}

func (b *BackgroundTaskRunner) runTask(ctx context.Context, task BackgroundTask) {
	ticker := time.NewTicker(task.interval)
	defer ticker.Stop()

	b.logger.Info("Started background task",
		zap.String("task", task.name),
		zap.Duration("interval", task.interval))

	for {
		select {
		case <-ctx.Done():
			b.logger.Info("Stopping background task due to context cancellation",
				zap.String("task", task.name),
				zap.Duration("interval", task.interval))
			return
		case <-ticker.C:
			start := time.Now()
			err := task.executable(ctx)
			duration := time.Since(start)

			if err != nil {
				b.logger.Error("Background task failed",
					zap.String("task", task.name),
					zap.Duration("interval", task.interval),
					zap.Duration("execution_time", duration),
					zap.Error(err))
			}
		}
	}
}
