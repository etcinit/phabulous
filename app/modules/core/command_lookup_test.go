package core

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

// Asserts that the default matchers are used, if no custom matchers are provided.
func Test_additionalMatchersNil(t *testing.T) {
	lookupCommand := LookupCommand{CustomMatchers: nil}
	matchers := lookupCommand.GetMatchers()
	assert.Equal(t, lookupCommand.GetDefaultMatchers(), matchers)
}

// Asserts that the custom matchers are used, if custom matchers are provided.
func Test_additionalMatchers(t *testing.T) {
	customMatchers := []string{"([D][0-9]{1,16})", "([T][0-9]{1,16})"}
	lookupCommand := LookupCommand{CustomMatchers: customMatchers}
	assert.Contains(t, lookupCommand.GetMatchers(), "([D][0-9]{1,16})")
	assert.Contains(t, lookupCommand.GetMatchers(), "([T][0-9]{1,16})")
}
