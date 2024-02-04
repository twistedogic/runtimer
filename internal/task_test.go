package internal

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type mockTask struct {
	duration    time.Duration
	shouldError bool
}

func (m mockTask) Run(ctx context.Context) error {
	time.Sleep(m.duration)
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if m.shouldError {
			return fmt.Errorf("error")
		}
	}
	return nil
}

func Test_observeTask(t *testing.T) {
	cases := map[string]struct {
		name, result      string
		duration, timeout time.Duration
		shouldError       bool
	}{
		"success": {
			name:    "success",
			result:  "success",
			timeout: 2 * time.Second,
		},
		"error": {
			name:        "error",
			result:      "error",
			timeout:     2 * time.Second,
			shouldError: true,
		},
		"timeout": {
			name:     "timeout",
			result:   "timeout",
			duration: time.Second,
		},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			tCtx, _ := context.WithTimeout(context.TODO(), tc.timeout)
			task := mockTask{duration: tc.duration, shouldError: tc.shouldError}
			labels, _ := observeTask(tCtx, tc.name, task)
			if got := labels[nameLabel]; tc.name != got {
				t.Fatalf("name, want: %s, got: %s", tc.name, got)
			}
			if got := labels[resultLabel]; tc.result != got {
				t.Fatalf("result, want: %s, got: %s", tc.result, got)
			}
		})
	}
}
