package controllers

import (
	"testing"

	"github.com/jacobstr/confer"
	"github.com/stretchr/testify/assert"
)

func Test_typeAllowed(t *testing.T) {
	config := confer.NewConfig()
	config.Set("channels.feedTypes", []string{})

	assert.True(t, typeAllowed(config, "TASK"))

	config.Set("channels.feedTypes", []string{"DREV", "CMIT"})

	assert.True(t, typeAllowed(config, "DREV"))
	assert.True(t, typeAllowed(config, "CMIT"))
	assert.False(t, typeAllowed(config, "TASK"))
}
