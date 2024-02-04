package internal

import (
	"context"
	"os/exec"
)

type ProcessTaskFactory struct{}

func (p ProcessTaskFactory) Type() string { return "process" }
func (f ProcessTaskFactory) Task(p Parameters) (Task, error) {
	cmd, err := p.String("cmd")
	if err != nil {
		return nil, err
	}
	args, err := p.StringSlice("args")
	if err != nil {
		return nil, err
	}
	return ProcessTask{name: cmd, args: args}, nil
}

type ProcessTask struct {
	name string
	args []string
}

func (p ProcessTask) Run(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, p.name, p.args...)
	return cmd.Run()
}
