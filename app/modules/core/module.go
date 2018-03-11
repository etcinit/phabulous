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
	lookupCommandCustomMatchers := getLookupCommandCustomMatchers(m.Config)

	return []interfaces.Command{
		&LookupCommand{CustomMatchers: lookupCommandCustomMatchers},
		&SummonCommand{},
		&MemeCommand{},
		&ModulesCommand{},
		&HelpCommand{},
		&TaskCommand{},
	}
}

// getLookupCommandCustomMatchers checks the configuration for custom matchers and loads them.
func getLookupCommandCustomMatchers(config *confer.Config) []string {
	var lookupCommandCustomMatchers []string

	if config.InConfig("commands.lookup.matchers") {
		configuredMatchers := config.GetStringSlice("commands.lookup.matchers")

		if len(configuredMatchers) > 0 {
			lookupCommandCustomMatchers = configuredMatchers
		}
	}

	return lookupCommandCustomMatchers
}