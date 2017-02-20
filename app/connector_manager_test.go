package app_test

import (
	"errors"
	"testing"

	"github.com/etcinit/phabulous/app"
	"github.com/etcinit/phabulous/app/interfaces"
	"github.com/etcinit/phabulous/app/messages"
	"github.com/etcinit/phabulous/app/testing/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestInterfaces(t *testing.T) {
	var connector interfaces.Connector = &app.ConnectorManager{}

	connector.Boot()
}

func Test_ConnectorManager_RegisterConnector(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	connector := &app.ConnectorManager{}

	connA := mocks.NewMockConnector(mockCtrl)
	connB := mocks.NewMockConnector(mockCtrl)

	connector.RegisterConnector(connA)
	connector.RegisterConnector(connB)
}

func Test_ConnectorManager_LoadModules(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	connector := &app.ConnectorManager{}

	connA := mocks.NewMockConnector(mockCtrl)
	connB := mocks.NewMockConnector(mockCtrl)

	module := mocks.NewMockModule(mockCtrl)
	modules := []interfaces.Module{module}

	connA.EXPECT().LoadModules(modules).Times(1)
	connB.EXPECT().LoadModules(modules).Times(1)

	connector.RegisterConnector(connA)
	connector.RegisterConnector(connB)

	connector.LoadModules(modules)
}

func Test_ConnectorManager_Boot(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	connector := &app.ConnectorManager{}

	connA := mocks.NewMockConnector(mockCtrl)
	connB := mocks.NewMockConnector(mockCtrl)
	connC := mocks.NewMockConnector(mockCtrl)

	err := errors.New("some error")

	connA.EXPECT().Boot().Times(2)
	connB.EXPECT().Boot().Times(2)
	connC.EXPECT().Boot().Times(1).Return(err)

	connector.RegisterConnector(connA)
	connector.RegisterConnector(connB)

	assert.Nil(t, connector.Boot())

	connector.RegisterConnector(connC)

	assert.Equal(t, err, connector.Boot())
}

func Test_ConnectorManager_Post(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	connector := &app.ConnectorManager{}

	connA := mocks.NewMockConnector(mockCtrl)
	connB := mocks.NewMockConnector(mockCtrl)

	connector.RegisterConnector(connA)
	connector.RegisterConnector(connB)

	connA.EXPECT().Post("somechannel", "awesome", messages.IconDefault, true)
	connB.EXPECT().Post("somechannel", "awesome", messages.IconDefault, true)

	connector.Post("somechannel", "awesome", messages.IconDefault, true)
}

func Test_ConnectorManager_PostImage(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	connector := &app.ConnectorManager{}

	connA := mocks.NewMockConnector(mockCtrl)
	connB := mocks.NewMockConnector(mockCtrl)

	connector.RegisterConnector(connA)
	connector.RegisterConnector(connB)

	connA.EXPECT().PostImage("chan", "awesome", "url", messages.IconDefault, true)
	connB.EXPECT().PostImage("chan", "awesome", "url", messages.IconDefault, true)

	connector.PostImage("chan", "awesome", "url", messages.IconDefault, true)
}

func Test_ConnectorManager_PostFeed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	connector := &app.ConnectorManager{}

	connA := mocks.NewMockConnector(mockCtrl)
	connB := mocks.NewMockConnector(mockCtrl)

	connector.RegisterConnector(connA)
	connector.RegisterConnector(connB)

	connA.EXPECT().PostOnFeed("awesome")
	connB.EXPECT().PostOnFeed("awesome")

	assert.Nil(t, connector.PostOnFeed("awesome"))
}

func Test_ConnectorManager_PostFeed_WithError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	connector := &app.ConnectorManager{}

	connA := mocks.NewMockConnector(mockCtrl)
	connB := mocks.NewMockConnector(mockCtrl)

	err := errors.New("some error")

	connector.RegisterConnector(connA)
	connector.RegisterConnector(connB)

	connA.EXPECT().PostOnFeed("awesome").Return(err)

	assert.Equal(t, err, connector.PostOnFeed("awesome"))
}
