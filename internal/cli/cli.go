package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/clbiggs/image-builder/internal/app"
	"github.com/clbiggs/image-builder/internal/manifest"
	"github.com/clbiggs/image-builder/internal/output"
	"github.com/clbiggs/image-builder/internal/viewmodel"
)

func Main(argv []string) int {
	if len(argv) < 2 {
		printHelp(os.Stderr)
		return 2
	}

	cmd := argv[1]
	args := argv[2:]

	switch cmd {
	case "--help", "-h", "help":
		printHelp(os.Stdout)
		return 0
	case "build":
		if err := runBuild(args); err != nil {
			fmt.Fprintln(os.Stderr, "ERROR:", err)
			return 1
		}
		return 0
	case "pull-images":
		if err := runPullImages(args); err != nil {
			fmt.Fprintln(os.Stderr, "ERROR:", err)
			return 1
		}
		return 0
	case "generate-build-matrix":
		if err := runGenerateBuildMatrix(args); err != nil {
			fmt.Fprintln(os.Stderr, "ERROR:", err)
			return 1
		}
		return 0
	case "merge-image-info":
		if err := runMergeImageInfo(args); err != nil {
			fmt.Fprintln(os.Stderr, "ERROR:", err)
			return 1
		}
		return 0
	case "publish-image-info":
		if err := runPublishImageInfo(args); err != nil {
			fmt.Fprintln(os.Stderr, "ERROR:", err)
			return 1
		}
		return 0

	// Parity command surface (stubs)
	case "generate-dockerfiles", "generate-readmes", "copy-acr-images", "copy-base-images",
		"clean-acr-images", "get-stale-images", "publish-manifest", "publish-mcr-docs",
		"post-publish-notification", "get-base-image-status", "generate-signing-payloads",
		"generate-eol-annotation-data-for-all-images", "generate-eol-annotation-data-for-publish",
		"annotate-eol-digests", "wait-for-mar-annotation-ingestion",
		"ingest-kusto-image-info", "show-image-stats":
		fmt.Fprintln(os.Stderr, "ERROR:", app.ErrNotImplemented)
		return 1
	default:
		fmt.Fprintln(os.Stderr, "Unknown command:", cmd)
		printHelp(os.Stderr)
		return 2
	}
}

type commonFlags struct {
	manifest         string
	repo             string
	image            string
	os               string
	arch             string
	dockerfile       string
	tag              string
	registryOverride string
}

func (c *commonFlags) add(fs *flag.FlagSet) {
	fs.StringVar(&c.manifest, "manifest", "manifest.json", "Path to manifest JSON")
	fs.StringVar(&c.repo, "repo", "", "Filter by repo name")
	fs.StringVar(&c.image, "image", "", "Filter by image name")
	fs.StringVar(&c.os, "os", "", "Filter by OS")
	fs.StringVar(&c.arch, "arch", "", "Filter by architecture")
	fs.StringVar(&c.dockerfile, "dockerfile", "", "Filter by dockerfile path")
	fs.StringVar(&c.tag, "tag", "", "Filter by tag name")
	fs.StringVar(&c.registryOverride, "registry", "", "Override registry from manifest")
}

func (c *commonFlags) filter() viewmodel.ManifestFilter {
	return viewmodel.ManifestFilter{
		Repo:       c.repo,
		Image:      c.image,
		OS:         c.os,
		Arch:       c.arch,
		Dockerfile: c.dockerfile,
		Tag:        c.tag,
	}
}

func loadVM(path string, f viewmodel.ManifestFilter) (*manifest.Manifest, *viewmodel.ManifestInfo, error) {
	m, err := manifest.Load(path)
	if err != nil {
		return nil, nil, err
	}
	mi, err := viewmodel.Build(m, f)
	if err != nil {
		return nil, nil, err
	}
	return m, mi, nil
}

func runBuild(args []string) error {
	fs := flag.NewFlagSet("build", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var c commonFlags
	c.add(fs)

	var contextDir string
	var outImageInfo string
	var push bool
	var buildArgs multiString

	fs.StringVar(&contextDir, "context", ".", "Docker build context directory")
	fs.StringVar(&outImageInfo, "output-image-info", "", "Write image-info JSON to this path")
	fs.BoolVar(&push, "push", false, "Push built images to registry")
	fs.Var(&buildArgs, "build-arg", "Docker build args (KEY=VALUE), repeatable")

	if err := fs.Parse(args); err != nil {
		return err
	}

	_, mi, err := loadVM(c.manifest, c.filter())
	if err != nil {
		return err
	}

	ba := map[string]string{}
	for _, kv := range buildArgs {
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) == 2 {
			ba[parts[0]] = parts[1]
		}
	}

	a := app.New()
	return a.Build(context.Background(), mi, app.BuildOptions{
		ContextDir:       contextDir,
		OutputImageInfo:  outImageInfo,
		Push:             push,
		RegistryOverride: c.registryOverride,
		BuildArgs:        ba,
	})
}

func runPullImages(args []string) error {
	fs := flag.NewFlagSet("pull-images", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var c commonFlags
	c.add(fs)

	if err := fs.Parse(args); err != nil {
		return err
	}

	_, mi, err := loadVM(c.manifest, c.filter())
	if err != nil {
		return err
	}

	a := app.New()
	return a.PullImages(context.Background(), mi, c.registryOverride)
}

func runGenerateBuildMatrix(args []string) error {
	fs := flag.NewFlagSet("generate-build-matrix", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var c commonFlags
	c.add(fs)

	out := fs.String("o", "build-matrix.json", "Output path for matrix JSON")

	if err := fs.Parse(args); err != nil {
		return err
	}

	_, mi, err := loadVM(c.manifest, c.filter())
	if err != nil {
		return err
	}

	m := app.GenerateBuildMatrix(mi)
	return m.Write(*out)
}

func runMergeImageInfo(args []string) error {
	fs := flag.NewFlagSet("merge-image-info", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var inputs multiString
	out := fs.String("o", "image-info.merged.json", "Output merged file")
	fs.Var(&inputs, "i", "Input image-info file (repeatable)")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if len(inputs) == 0 {
		return errors.New("at least one -i input file is required")
	}

	merged := output.New()
	for _, p := range inputs {
		f, err := output.Read(p)
		if err != nil {
			return err
		}
		merged.Images = append(merged.Images, f.Images...)
		for k, v := range f.Metadata {
			merged.Metadata[k] = v
		}
	}
	return merged.Write(*out)
}

func runPublishImageInfo(args []string) error {
	fs := flag.NewFlagSet("publish-image-info", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	in := fs.String("input", "image-info.json", "Input image-info file")
	out := fs.String("output", "", "Output path (default overwrites input)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	f, err := output.Read(*in)
	if err != nil {
		return err
	}
	if *out == "" {
		*out = *in
	}
	return f.Write(*out)
}

type multiString []string

func (m *multiString) String() string { return strings.Join(*m, ",") }
func (m *multiString) Set(v string) error {
	*m = append(*m, v)
	return nil
}

func printHelp(w *os.File) {
	fmt.Fprintln(w, `image-builder (Go)

Usage:
  image-builder <command> [flags]

Implemented commands:
  build
  pull-images
  generate-build-matrix
  merge-image-info
  publish-image-info

Parity-stub commands (present but not implemented yet):
  generate-dockerfiles
  generate-readmes
  copy-acr-images
  copy-base-images
  clean-acr-images
  get-stale-images
  publish-manifest
  publish-mcr-docs
  post-publish-notification
  get-base-image-status
  generate-signing-payloads
  generate-eol-annotation-data-for-all-images
  generate-eol-annotation-data-for-publish
  annotate-eol-digests
  wait-for-mar-annotation-ingestion
  ingest-kusto-image-info
  show-image-stats

Common flags:
  --manifest <path>
  --repo <name>
  --image <name>
  --os <os>
  --arch <arch>
  --dockerfile <path>
  --tag <tag>
  --registry <override>`)
}
