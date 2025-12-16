package app

import (
	"encoding/json"
	"os"
	"time"

	"github.com/clbiggs/image-builder/internal/viewmodel"
)

type BuildMatrix struct {
	CreatedUTC string      `json:"createdUtc"`
	Jobs       []MatrixJob `json:"jobs"`
}

type MatrixJob struct {
	Repo       string `json:"repo"`
	Image      string `json:"image"`
	OS         string `json:"os,omitempty"`
	Arch       string `json:"arch,omitempty"`
	Dockerfile string `json:"dockerfile,omitempty"`
	Tag        string `json:"tag"`
}

func GenerateBuildMatrix(mi *viewmodel.ManifestInfo) *BuildMatrix {
	m := &BuildMatrix{
		CreatedUTC: time.Now().UTC().Format(time.RFC3339),
		Jobs:       []MatrixJob{},
	}
	for _, repo := range mi.Repos {
		for _, img := range repo.Images {
			for _, p := range img.Platforms {
				for _, t := range p.Tags {
					m.Jobs = append(m.Jobs, MatrixJob{
						Repo:       repo.Name,
						Image:      img.Name,
						OS:         p.OS,
						Arch:       p.Arch,
						Dockerfile: p.Dockerfile,
						Tag:        t.Name,
					})
				}
			}
		}
	}
	return m
}

func (m *BuildMatrix) Write(path string) error {
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0644)
}
