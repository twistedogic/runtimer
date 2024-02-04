package internal

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Runner struct {
	clk       clock.Clock
	interval  time.Duration
	counter   *prometheus.CounterVec
	histogram *prometheus.HistogramVec
	tasks     map[string]Task
}

func NewRunner(path string) (Runner, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Runner{}, err
	}
	c, err := parseConfig(b)
	if err != nil {
		return Runner{}, err
	}
	return NewRunnerWithConfig(c)
}

func NewRunnerWithConfig(c Config) (Runner, error) {
	interval, err := c.IntervalDuration()
	if err != nil {
		return Runner{}, err
	}
	buckets, err := c.HistogramBuckets()
	if err != nil {
		return Runner{}, err
	}
	tasks, err := c.ToTasks(WebTaskFactory{})
	if err != nil {
		return Runner{}, err
	}
	return newRunner(interval, buckets, tasks), nil
}

func newRunner(interval time.Duration, buckets []float64, tasks map[string]Task) Runner {
	return Runner{
		clk:       clock.New(),
		interval:  interval,
		counter:   newCounter(),
		histogram: newHistogram(buckets),
		tasks:     tasks,
	}
}

func (r Runner) runOnce(ctx context.Context, name string, task Task) {
	tCtx, cancel := context.WithTimeout(ctx, r.interval)
	defer cancel()
	labels, duration := observeTask(tCtx, name, task)
	r.counter.With(labels).Inc()
	r.histogram.With(labels).Observe(duration.Seconds())
}

func (r Runner) loop(ctx context.Context, name string, task Task) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			go r.runOnce(ctx, name, task)
			r.clk.Sleep(r.interval)
		}
	}
}

func (r Runner) run(ctx context.Context) {
	for name, task := range r.tasks {
		go r.loop(ctx, name, task)
	}
}

func (r Runner) Start(ctx context.Context, addr string) error {
	srv := &http.Server{Addr: addr, Handler: promhttp.Handler()}
	r.run(ctx)
	go func() {
		<-ctx.Done()
		srv.Close()
	}()
	return srv.ListenAndServe()
}
