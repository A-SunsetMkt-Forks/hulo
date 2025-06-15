package release

type Resource struct {
	Path        string
	InstallPath string
	Executable  bool
}

type Package struct {
	OS        string
	Arch      string
	URL       string
	SHA256    string
	Resources []Resource
}

type Food struct {
	Name        string
	Description string
	Homepage    string
	Version     string
	License     string
	Packages    []Package
}

func (Food) Release() {}
