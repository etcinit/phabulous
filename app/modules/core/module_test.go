package core

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/jacobstr/confer"
)

// Asserts that the custom matchers from the configuration are loaded
func Test_LoadCustomMatchersForLookupCommand(t *testing.T) {
	customMatchers := []string{"([D][0-9]{1,16})", "([T][0-9]{1,16})"}

	config := confer.NewConfig()
	config.Set("commands.lookup.matchers", customMatchers)

	module := Module{Config: config}
	commands := module.GetCommands()

	for _, command := range commands {
		switch lookupCommand := command.(type) {
		case *LookupCommand:
			assert.Equal(t, customMatchers, lookupCommand.GetMatchers())
		}
	}
}

// Asserts that the default matchers are used if the config array is empty
func Test_LoadDefaultLookupMatchersEmptyConfig(t *testing.T) {
	config := confer.NewConfig()
	config.Set("commands.lookup.matchers", []string{})

	module := Module{Config: config}
	commands := module.GetCommands()

	for _, command := range commands {
		switch lookupCommand := command.(type) {
		case *LookupCommand:
			assert.Equal(t, lookupCommand.GetDefaultMatchers(), lookupCommand.GetMatchers())
		}
	}
}

// Asserts that the default matchers are used if the config array is missing
func Test_LoadDefaultLookupMatchersMissingConfig(t *testing.T) {
	config := confer.NewConfig()
	module := Module{Config: config}
	commands := module.GetCommands()

	for _, command := range commands {
		switch lookupCommand := command.(type) {
		case *LookupCommand:
			assert.Equal(t, lookupCommand.GetDefaultMatchers(), lookupCommand.GetMatchers())
		}
	}
}