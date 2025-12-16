package manifest

import "encoding/json"

type Manifest struct {
	Registry  string            `json:"registry"`
	Variables map[string]string `json:"variables"`
	Repos     []Repo            `json:"repos"`

	Extra map[string]json.RawMessage `json:"-"`
}

func (m *Manifest) UnmarshalJSON(b []byte) error {
	type alias Manifest
	var a alias
	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}
	*m = Manifest(a)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	delete(raw, "registry")
	delete(raw, "variables")
	delete(raw, "repos")
	m.Extra = raw
	return nil
}

type Repo struct {
	Name   string  `json:"name"`
	Images []Image `json:"images"`

	Extra map[string]json.RawMessage `json:"-"`
}

func (r *Repo) UnmarshalJSON(b []byte) error {
	type alias Repo
	var a alias
	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}
	*r = Repo(a)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	delete(raw, "name")
	delete(raw, "images")
	r.Extra = raw
	return nil
}

type Image struct {
	Name      string     `json:"name"`
	Platforms []Platform `json:"platforms"`

	Extra map[string]json.RawMessage `json:"-"`
}

func (i *Image) UnmarshalJSON(b []byte) error {
	type alias Image
	var a alias
	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}
	*i = Image(a)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	delete(raw, "name")
	delete(raw, "platforms")
	i.Extra = raw
	return nil
}

type Platform struct {
	OS         string `json:"os"`
	Arch       string `json:"arch"`
	Dockerfile string `json:"dockerfile"`
	BaseImage  string `json:"baseImage,omitempty"`
	Tags       []Tag  `json:"tags"`

	Extra map[string]json.RawMessage `json:"-"`
}

func (p *Platform) UnmarshalJSON(b []byte) error {
	type alias Platform
	var a alias
	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}
	*p = Platform(a)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	delete(raw, "os")
	delete(raw, "arch")
	delete(raw, "dockerfile")
	delete(raw, "baseImage")
	delete(raw, "tags")
	p.Extra = raw
	return nil
}

type Tag struct {
	Name          string `json:"name"`
	Documentation string `json:"documentation,omitempty"`

	Extra map[string]json.RawMessage `json:"-"`
}

func (t *Tag) UnmarshalJSON(b []byte) error {
	type alias Tag
	var a alias
	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}
	*t = Tag(a)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	delete(raw, "name")
	delete(raw, "documentation")
	t.Extra = raw
	return nil
}
