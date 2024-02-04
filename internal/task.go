package internal

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Task interface {
	Run(context.Context) error
}

type TaskFactory interface {
	Type() string
	Task(Parameters) (Task, error)
}

func observeTask(ctx context.Context, name string, t Task) (prometheus.Labels, time.Duration) {
	labels := map[string]string{
		nameLabel:   name,
		resultLabel: "error",
	}
	start := time.Now()
	err := t.Run(ctx)
	dur := time.Now().Sub(start)
	if err != nil {
		log.Printf("task %q failed: %v", name, err)
	}
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		labels[resultLabel] = "timeout"
	case err == nil:
		labels[resultLabel] = "success"
	}
	return labels, dur
}
