package manifest

import "fmt"

func Validate(m *Manifest) error {
	for _, r := range m.Repos {
		if r.Name == "" {
			return fmt.Errorf("repo name is required")
		}
		for _, img := range r.Images {
			if img.Name == "" {
				return fmt.Errorf("image name is required (repo %s)", r.Name)
			}
			for _, p := range img.Platforms {
				if p.Dockerfile == "" {
					return fmt.Errorf("dockerfile is required (repo %s image %s)", r.Name, img.Name)
				}
				for _, t := range p.Tags {
					if t.Name == "" {
						return fmt.Errorf("tag name is required (repo %s image %s dockerfile %s)", r.Name, img.Name, p.Dockerfile)
					}
				}
			}
		}
	}
	return nil
}
