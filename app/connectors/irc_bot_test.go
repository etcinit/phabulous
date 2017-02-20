package connectors

import (
	"testing"

	"github.com/etcinit/phabulous/app/interfaces"
	"github.com/etcinit/phabulous/app/modules/core"
	"github.com/fluffle/goirc/client"
	"github.com/jacobstr/confer"
	"github.com/stretchr/testify/assert"
)

func Test_IRCConnector_LoadModules(t *testing.T) {
	connector := IRCConnector{}
	connector.LoadModules([]interfaces.Module{&core.Module{}})

	assert.Equal(t, 0, len(connector.GetHandlers()))
	assert.Equal(t, 0, len(connector.GetIMHandlers()))

	modules := []interfaces.Module{&core.Module{}}

	connector.client = client.SimpleClient("phabulous")
	connector.LoadModules(modules)

	assert.True(t, len(connector.GetHandlers()) > 0)
	assert.True(t, len(connector.GetIMHandlers()) > 0)
	assert.Equal(t, modules, connector.GetModules())
}

func Test_IRCConnector_GetConfig(t *testing.T) {
	connector := IRCConnector{}
	connector.config = confer.NewConfig()

	assert.Equal(t, connector.config, connector.GetConfig())
}

func Test_IRCConnector_GetUsername(t *testing.T) {
	connector := IRCConnector{}

	username, err := connector.GetUsername("same")
	assert.Nil(t, err)
	assert.Equal(t, "same", username)
}

func Test_IRCConnector_GetUsageHandler(t *testing.T) {
	connector := IRCConnector{}

	connector.GetUsageHandler()
}

func Test_IRCConnector_StartTyping(t *testing.T) {
	connector := IRCConnector{}

	connector.StartTyping("somechannel")
}
