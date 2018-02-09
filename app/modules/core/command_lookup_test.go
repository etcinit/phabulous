package core

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func Test_additionalMatchersNil(t *testing.T) {
	lookup_command := LookupCommand{AdditionalMatchers: nil}
	assert.NotEmpty(t, lookup_command.GetMatchers())
}

func Test_additionalMatchers(t *testing.T) {
	additionalMatchers := []string{"([D][0-9]{1,16})", "([D][0-1]{1,16})"}
	lookupCommand := LookupCommand{AdditionalMatchers: additionalMatchers}
	assert.Contains(t, lookupCommand.GetMatchers(), "([D][0-9]{1,16})")
	assert.Contains(t, lookupCommand.GetMatchers(), "([D][0-1]{1,16})")
}
