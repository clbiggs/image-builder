package viewmodel

type ManifestFilter struct {
	Repo       string
	Image      string
	OS         string
	Arch       string
	Dockerfile string
	Tag        string
}
