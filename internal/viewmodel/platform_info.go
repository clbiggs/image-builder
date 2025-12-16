package viewmodel

type PlatformInfo struct {
	OS         string
	Arch       string
	Dockerfile string
	BaseImage  string
	Tags       []TagInfo
}
