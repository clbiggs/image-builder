package output

import (
	"encoding/json"
	"os"
	"time"
)

type ImageInfoFile struct {
	CreatedUTC string         `json:"createdUtc"`
	Images     []BuiltImage   `json:"images"`
	Metadata   map[string]any `json:"metadata,omitempty"`
}

type BuiltImage struct {
	Repo       string `json:"repo"`
	Image      string `json:"image"`
	OS         string `json:"os,omitempty"`
	Arch       string `json:"arch,omitempty"`
	Dockerfile string `json:"dockerfile,omitempty"`
	Tag        string `json:"tag"`
	FullRef    string `json:"fullRef"`
	Digest     string `json:"digest,omitempty"`
}

func New() *ImageInfoFile {
	return &ImageInfoFile{
		CreatedUTC: time.Now().UTC().Format(time.RFC3339),
		Images:     []BuiltImage{},
		Metadata:   map[string]any{},
	}
}

func (f *ImageInfoFile) Write(path string) error {
	b, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0644)
}

func Read(path string) (*ImageInfoFile, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var f ImageInfoFile
	if err := json.Unmarshal(b, &f); err != nil {
		return nil, err
	}
	if f.Metadata == nil {
		f.Metadata = map[string]any{}
	}
	return &f, nil
}
