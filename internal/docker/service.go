package docker

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/clbiggs/image-builder/internal/util"
)

type Service interface {
	Pull(ctx context.Context, ref string) error
	Build(ctx context.Context, dockerfile string, contextDir string, tag string, buildArgs map[string]string) error
	Tag(ctx context.Context, src string, dst string) error
	Push(ctx context.Context, ref string) error
	InspectDigest(ctx context.Context, ref string) (string, error)
}

type CLI struct {
	Runner util.CmdRunner
}

func NewCLI(runner util.CmdRunner) *CLI {
	return &CLI{Runner: runner}
}

func (d *CLI) Pull(ctx context.Context, ref string) error {
	return d.Runner.Run(ctx, "docker", "pull", ref)
}

func (d *CLI) Build(ctx context.Context, dockerfile string, contextDir string, tag string, buildArgs map[string]string) error {
	args := []string{"build", "-t", tag, "-f", dockerfile}
	for k, v := range buildArgs {
		args = append(args, "--build-arg", fmt.Sprintf("%s=%s", k, v))
	}
	if contextDir == "" {
		contextDir = "."
	}
	contextDir = filepath.Clean(contextDir)
	args = append(args, contextDir)
	return d.Runner.Run(ctx, "docker", args...)
}

func (d *CLI) Tag(ctx context.Context, src string, dst string) error {
	return d.Runner.Run(ctx, "docker", "tag", src, dst)
}

func (d *CLI) Push(ctx context.Context, ref string) error {
	return d.Runner.Run(ctx, "docker", "push", ref)
}

func (d *CLI) InspectDigest(ctx context.Context, ref string) (string, error) {
	return d.Runner.Output(ctx, "docker", "inspect", "--format", "{{index .RepoDigests 0}}", ref)
}
