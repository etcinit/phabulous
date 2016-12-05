package app_test

import (
	"testing"

	"github.com/etcinit/phabulous/app"
	"github.com/etcinit/phabulous/app/interfaces"
)

func TestInterfaces(t *testing.T) {
	var connector interfaces.Connector = &app.ConnectorManager{}

	connector.Boot()
}
