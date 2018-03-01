package core

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
	return "core"
}

// GetCommands returns the commands provided by this module.
func (m *Module) GetCommands() []interfaces.Command {
	var lookupCommandAdditionalMatchers []string
	if m.Config.InConfig("commands.lookup.matchers") {
		lookupCommandAdditionalMatchers = m.Config.GetStringSlice("commands.lookup.matchers")
	}
	return []interfaces.Command{

		&LookupCommand{AdditionalMatchers:lookupCommandAdditionalMatchers},
		&SummonCommand{},
		&MemeCommand{},
		&ModulesCommand{},
		&HelpCommand{},
	}
}
