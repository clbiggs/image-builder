package app

import (
	"context"
	"errors"

	"github.com/clbiggs/image-builder/internal/viewmodel"
)

var ErrNotImplemented = errors.New("not implemented in this scaffold")

func (a *App) GernerateDockerfiles(ctx context.Context, mi *viewmodel.ManifestInfo, outDir string) error {
	return ErrNotImplemented
}

func (a *App) GenerateReadmes(ctx context.Context, mi *viewmodel.ManifestInfo, outDir string) error {
	return ErrNotImplemented
}
