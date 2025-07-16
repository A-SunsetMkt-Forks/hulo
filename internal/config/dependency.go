package config

type Dependency struct {
	Version      string       `yaml:"version"`
	Resolved     string       `yaml:"resolved"`
	Commit       string       `yaml:"commit"`
	Dependencies []Dependency `yaml:"dependencies"`
	Dev          bool         `yaml:"dev"`
}
