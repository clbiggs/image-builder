package viewmodel

import "github.com/clbiggs/image-builder/internal/manifest"

type ManifestInfo struct {
	Manifest *manifest.Manifest
	Repos    []RepoInfo
}

func Build(m *manifest.Manifest, filter ManifestFilter) (*ManifestInfo, error) {
	if err := manifest.Validate(m); err != nil {
		return nil, err
	}

	repos := make([]RepoInfo, 0, len(m.Repos))
	for _, r := range m.Repos {
		if filter.Repo != "" && r.Name != filter.Repo {
			continue
		}
		ri := RepoInfo{Name: r.Name}
		for _, img := range r.Images {
			if filter.Image != "" && img.Name != filter.Image {
				continue
			}
			ii := ImageInfo{Name: img.Name}
			for _, p := range img.Platforms {
				if filter.Dockerfile != "" && p.Dockerfile != filter.Dockerfile {
					continue
				}
				if filter.OS != "" && p.OS != filter.OS {
					continue
				}
				if filter.Arch != "" && p.Arch != filter.Arch {
					continue
				}
				pi := PlatformInfo{
					OS:         p.OS,
					Arch:       p.Arch,
					Dockerfile: p.Dockerfile,
					BaseImage:  p.BaseImage,
				}
				for _, t := range p.Tags {
					if filter.Tag != "" && t.Name != filter.Tag {
						continue
					}
					pi.Tags = append(pi.Tags, TagInfo{
						Name:          t.Name,
						Documentation: t.Documentation,
					})
				}
				if len(pi.Tags) > 0 {
					continue
				}
				ii.Platforms = append(ii.Platforms, pi)
			}
			if len(ii.Platforms) > 0 {
				continue
			}
			ri.Images = append(ri.Images, ii)
		}
		if len(ri.Images) > 0 {
			continue
		}
		repos = append(repos, ri)
	}

	return &ManifestInfo{Manifest: m, Repos: repos}, nil
}
