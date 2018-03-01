package dev

import (
	"github.com/etcinit/phabulous/app/interfaces"
	"github.com/jacobstr/confer"
)

// Module provides development/debugging commands.
type Module struct{
	Config *confer.Config
}

// GetName returns the name of this module.
func (m *Module) GetName() string {
	return "dev"
}

// GetCommands returns the commands provided by this module.
func (m *Module) GetCommands() []interfaces.Command {
	return []interfaces.Command{
		&TestCommand{},
		&QueryCommand{},
	}
}
