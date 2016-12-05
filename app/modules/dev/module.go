package dev

import "github.com/etcinit/phabulous/app/interfaces"

// Module provides development/debugging commands.
type Module struct{}

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
