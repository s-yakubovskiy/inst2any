package daemon

import (
	"context"
	"log"
)

type Worker interface {
	Work(context.Context)
	Enabled() bool
	Name() string
	FullName() string
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
		if worker.Enabled() {
			go worker.Work(ctx)
		} else {
			log.Println(worker.FullName(), "is disabled")
		}
	}
}
