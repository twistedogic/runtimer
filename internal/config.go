package internal

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/yaml.v3"
)

type Parameters map[string]interface{}

func (p Parameters) Int(field string) (int, error) {
	i, ok := p[field]
	if !ok {
		return 0, fmt.Errorf("field %q not found", field)
	}
	v, ok := i.(int)
	if !ok {
		return 0, fmt.Errorf("field %q is not int", field)
	}
	return v, nil
}

func (p Parameters) String(field string) (string, error) {
	i, ok := p[field]
	if !ok {
		return "", fmt.Errorf("field %q not found", field)
	}
	v, ok := i.(string)
	if !ok {
		return "", fmt.Errorf("field %q is not string", field)
	}
	return v, nil
}

func (p Parameters) StringSlice(field string) ([]string, error) {
	i, ok := p[field]
	if !ok {
		return nil, fmt.Errorf("field %q not found", field)
	}
	slice, ok := i.([]interface{})
	if !ok {
		return nil, fmt.Errorf("field %q is not []string", field)
	}
	v := make([]string, len(slice))
	for i, val := range slice {
		v[i] = fmt.Sprintf("%v", val)
	}
	return v, nil
}

type TaskConfig struct {
	Name       string     `yaml:"name"`
	Type       string     `yaml:"type"`
	Parameters Parameters `yaml:"parameters"`
}

type Config struct {
	Interval string       `yaml:"interval"`
	Buckets  []float64    `yaml:"buckets"`
	Tasks    []TaskConfig `yaml:"tasks"`
}

func parseConfig(b []byte) (Config, error) {
	var c Config
	return c, yaml.Unmarshal(b, &c)
}

func (c Config) IntervalDuration() (time.Duration, error) {
	return time.ParseDuration(c.Interval)
}

func (c Config) HistogramBuckets() ([]float64, error) {
	if len(c.Buckets) != 0 {
		return c.Buckets, nil
	}
	duration, err := c.IntervalDuration()
	if err != nil {
		return nil, err
	}
	return prometheus.ExponentialBucketsRange(0, duration.Seconds(), 5), nil
}

func (c Config) ToTasks(factories ...TaskFactory) (map[string]Task, error) {
	m := make(map[string]TaskFactory)
	for _, f := range factories {
		m[f.Type()] = f
	}
	tasks := make(map[string]Task)
	for _, config := range c.Tasks {
		f, ok := m[config.Type]
		if !ok {
			return nil, fmt.Errorf("no task found with type %q", config.Type)
		}
		task, err := f.Task(config.Parameters)
		if err != nil {
			return nil, err
		}
		tasks[config.Name] = task
	}
	return tasks, nil
}
