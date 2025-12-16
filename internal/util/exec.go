package util

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type CmdRunner interface {
	Run(ctx context.Context, cmd string, args ...string) error
	Output(ctx context.Context, cmd string, args ...string) (string, error)
}

type OSRunner struct {
	Logger *Logger
}

func NewOSRunner(logger *Logger) *OSRunner {
	return &OSRunner{Logger: logger}
}

func (r *OSRunner) Run(ctx context.Context, name string, args ...string) error {
	r.Logger.Infof("exec: %s %s", name, strings.Join(args, " "))
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (r *OSRunner) Output(ctx context.Context, name string, args ...string) (string, error) {
	r.Logger.Infof("exec(out): %s %s", name, strings.Join(args, " "))
	cmd := exec.CommandContext(ctx, name, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	out := strings.TrimSpace(stdout.String())
	if err != nil {
		return out, fmt.Errorf("%w: %s", err, strings.TrimSpace(stderr.String()))
	}
	return out, nil
}
