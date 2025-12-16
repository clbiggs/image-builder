package manifest

import (
	"encoding/json"
	"fmt"
	"os"
)

func Load(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parse manifest %s: %w", path, err)
	}

	if m.Variables == nil {
		m.Variables = map[string]string{}
	}

	return &m, nil
}
