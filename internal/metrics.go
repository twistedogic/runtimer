package internal

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	nameLabel   = "name"
	resultLabel = "result"
)

func newCounter() *prometheus.CounterVec {
	return promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "task_complete_total",
			Help: "total number of task completed",
		},
		[]string{"name", "result"},
	)
}

func newHistogram(buckets []float64) *prometheus.HistogramVec {
	return promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "task_complete_seconds",
			Help:    "duration of task taken to complete",
			Buckets: buckets,
		},
		[]string{"name", "result"},
	)
}
