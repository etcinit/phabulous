package connectors

import (
	"testing"

	"github.com/fluffle/goirc/client"
	"github.com/stretchr/testify/assert"
)

func Test_NewIRCMessage(t *testing.T) {
	line := &client.Line{}

	NewIRCMessage(line, "phabulous")
}

func Test_IRCMessage_GetChannel(t *testing.T) {
	line := &client.Line{
		Cmd:  client.PRIVMSG,
		Nick: "#somechannel",
		Args: []string{"#somechannel"},
	}

	message := NewIRCMessage(line, "phabulous")

	assert.Equal(t, "#somechannel", message.GetChannel())
}

func Test_IRCMessage_GetUserID(t *testing.T) {
	line := &client.Line{
		Nick: "bob",
	}

	message := NewIRCMessage(line, "phabulous")

	assert.Equal(t, "bob", message.GetUserID())
}

func Test_IRCMessage_GetContent(t *testing.T) {
	line := &client.Line{
		Args: []string{"#room", "Awesome message"},
	}

	message := NewIRCMessage(line, "phabulous")

	assert.Equal(t, "Awesome message", message.GetContent())
}

func Test_IRCMessage_GetProvider(t *testing.T) {
	message := NewIRCMessage(&client.Line{}, "phabulous")

	assert.Equal(t, "irc", message.GetProviderName())
}

func Test_IRCMessage_IsIM(t *testing.T) {
	line := &client.Line{
		Cmd:  client.PRIVMSG,
		Nick: "someone",
		Args: []string{"someone"},
	}

	message := NewIRCMessage(line, "phabulous")

	assert.True(t, message.IsIM())
}

func Test_IRCMessage_IsSelf(t *testing.T) {
	line := &client.Line{
		Cmd:  client.PRIVMSG,
		Nick: "phabulous",
		Args: []string{"#somechannel"},
	}

	message := NewIRCMessage(line, "phabulous")

	assert.True(t, message.IsSelf())
}
