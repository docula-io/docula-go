package adr

// Directory represents a configured adr directory.
type Directory struct {
	Path  string `yaml:"path"`
	Name  string `yaml:"name"`
	Index string `yaml:"index"`
}

// State repesents the internal state configuration of the adr commands
type State struct {
	Directories []Directory `yaml:"dirs"`
}
