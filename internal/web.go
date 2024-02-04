package internal

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type WebTaskFactory struct{}

func (w WebTaskFactory) Type() string { return "web" }

func (w WebTaskFactory) Task(p Parameters) (Task, error) {
	method, err := p.String("method")
	if err != nil {
		return nil, err
	}
	url, err := p.String("url")
	if err != nil {
		return nil, err
	}
	statusCode, err := p.Int("status_code")
	if err != nil {
		return nil, err
	}
	body, err := p.String("url")
	if err != nil {
		return nil, err
	}
	return WebTask{
		client:     &http.Client{},
		url:        url,
		method:     method,
		body:       body,
		statusCode: statusCode,
	}, nil
}

type WebTask struct {
	client            *http.Client
	url, method, body string
	statusCode        int
}

func (w WebTask) Run(ctx context.Context) error {
	r := strings.NewReader(w.body)
	method := strings.ToUpper(w.method)
	req, err := http.NewRequestWithContext(ctx, method, w.url, r)
	if err != nil {
		return err
	}
	res, err := w.client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != w.statusCode {
		return fmt.Errorf("status code want %d, got %d", w.statusCode, res.StatusCode)
	}
	return nil
}
