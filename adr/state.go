package adr

// Directory represents a configured adr directory.
type Directory struct {
	Path string `yaml:"path"`
}

// State repesents the internal state configuration of the adr commands
type State struct {
	Directories []Directory `yaml:"dirs"`
}
