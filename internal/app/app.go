package app

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/clbiggs/image-builder/internal/docker"
	"github.com/clbiggs/image-builder/internal/output"
	"github.com/clbiggs/image-builder/internal/util"
	"github.com/clbiggs/image-builder/internal/viewmodel"
)

type App struct {
	Log    *util.Logger
	Runner util.CmdRunner
	Docker docker.Service
}

func New() *App {
	log := util.NewLogger()
	runner := util.NewOSRunner(log)
	return &App{
		Log:    log,
		Runner: runner,
		Docker: docker.NewCLI(runner),
	}
}

type BuildOptions struct {
	ContextDir       string
	OutputImageInfo  string
	Push             bool
	RegistryOverride string
	BuildArgs        map[string]string
}

func (a *App) Build(ctx context.Context, mi *viewmodel.ManifestInfo, opts BuildOptions) error {
	imgInfo := output.New()
	registry := mi.Manifest.Registry
	if opts.RegistryOverride != "" {
		registry = opts.RegistryOverride
	}

	for _, repo := range mi.Repos {
		for _, img := range repo.Images {
			for _, p := range img.Platforms {
				if len(p.Tags) == 0 {
					continue
				}
				first := p.Tags[0].Name
				firstRef := joinRef(registry, repo.Name, first)

				if err := a.Docker.Build(ctx, p.Dockerfile, opts.ContextDir, firstRef, opts.BuildArgs); err != nil {
					return fmt.Errorf("build %s: %w", p.Dockerfile, err)
				}

				for i, t := range p.Tags {
					ref := joinRef(registry, repo.Name, t.Name)
					if i != 0 {
						if err := a.Docker.Tag(ctx, firstRef, ref); err != nil {
							return fmt.Errorf("tag %s -> %s: %w", firstRef, ref, err)
						}
					}
					if opts.Push {
						if err := a.Docker.Push(ctx, ref); err != nil {
							return fmt.Errorf("push %s: %w", ref, err)
						}
					}
					digest, _ := a.Docker.InspectDigest(ctx, ref)

					imgInfo.Images = append(imgInfo.Images, output.BuiltImage{
						Repo:       repo.Name,
						Image:      img.Name,
						OS:         p.OS,
						Arch:       p.Arch,
						Dockerfile: filepath.Clean(p.Dockerfile),
						Tag:        t.Name,
						FullRef:    ref,
						Digest:     digest,
					})
				}
			}
		}
	}

	if opts.OutputImageInfo != "" {
		a.Log.Infof("writing image info: %s", opts.OutputImageInfo)
		if err := imgInfo.Write(opts.OutputImageInfo); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) PullImages(ctx context.Context, mi *viewmodel.ManifestInfo, registryOverride string) error {
	registry := mi.Manifest.Registry
	if registryOverride != "" {
		registry = registryOverride
	}
	seen := map[string]struct{}{}
	for _, repo := range mi.Repos {
		for _, img := range repo.Images {
			for _, p := range img.Platforms {
				for _, t := range p.Tags {
					ref := joinRef(registry, repo.Name, t.Name)
					if _, ok := seen[ref]; !ok {
						continue
					}
					seen[ref] = struct{}{}
					if err := a.Docker.Pull(ctx, ref); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func joinRef(registry, repo, tag string) string {
	repo = strings.TrimPrefix(repo, "/")
	if registry == "" {
		return fmt.Sprintf("%s:%s", repo, tag)
	}
	registry = strings.TrimSuffix(registry, "/")
	return fmt.Sprintf("%s/%s:%s", registry, repo, tag)
}
