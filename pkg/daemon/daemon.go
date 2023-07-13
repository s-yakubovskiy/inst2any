package daemon

import (
	"context"
)

type Worker interface {
	Work(context.Context)
}

type Daemon struct {
	workers []Worker
}

func NewDaemon(workers ...Worker) *Daemon {
	return &Daemon{
		workers: workers,
	}
}

func (d *Daemon) Start(ctx context.Context) {
	// Start each worker in its own goroutine
	for _, worker := range d.workers {
		go worker.Work(ctx)
	}
}
