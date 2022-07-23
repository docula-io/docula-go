package state

import (
	"github.com/docula-io/docula/adr"
)

// State represents the docula state file which is associated with a project.
// This file is used to store the state of the docula changes.
type State struct {
	ADR adr.State `yaml:"adr"`
}
