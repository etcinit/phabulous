package connectors

import (
	"regexp"
	"testing"

	"github.com/etcinit/phabulous/app/interfaces"
	"github.com/etcinit/phabulous/app/testing/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_processMessage_WithSelf(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockBot := mocks.NewMockBot(mockCtrl)
	mockMsg := mocks.NewMockMessage(mockCtrl)

	mockMsg.EXPECT().IsSelf().Times(1).Return(true)

	processMessage(mockBot, mockMsg)
}

func Test_processMessage_WithPublic(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockBot := mocks.NewMockBot(mockCtrl)
	mockMsg := mocks.NewMockMessage(mockCtrl)
	mockTuple := mocks.NewMockHandlerTuple(mockCtrl)

	mockMsg.EXPECT().IsSelf().Times(1).Return(false)
	mockMsg.EXPECT().GetContent().Times(1).Return("Hello World")
	mockMsg.EXPECT().IsIM().Times(1).Return(false)

	mockBot.EXPECT().GetHandlers().Times(1).Return([]interfaces.HandlerTuple{
		mockTuple,
	})

	ran := make(chan bool)

	pattern, _ := regexp.Compile("Hello World")
	mockTuple.EXPECT().GetPattern().Times(1).Return(pattern)
	mockTuple.EXPECT().GetHandler().Times(1).Return(
		func(bot interfaces.Bot, msg interfaces.Message, result []string) {
			assert.Equal(t, bot, mockBot)
			assert.Equal(t, msg, mockMsg)
			assert.Equal(t, result, []string{"Hello World"})

			ran <- true
		},
	)

	processMessage(mockBot, mockMsg)

	<-ran
}

func Test_processMessage_WithHandledIM(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockBot := mocks.NewMockBot(mockCtrl)
	mockMsg := mocks.NewMockMessage(mockCtrl)
	mockTuple := mocks.NewMockHandlerTuple(mockCtrl)

	mockMsg.EXPECT().IsSelf().Times(1).Return(false)
	mockMsg.EXPECT().GetContent().Times(1).Return("Hello World")
	mockMsg.EXPECT().IsIM().Times(1).Return(true)

	mockBot.EXPECT().GetIMHandlers().Times(1).Return([]interfaces.HandlerTuple{
		mockTuple,
	})

	ran := make(chan bool)

	pattern, _ := regexp.Compile("Hello World")
	mockTuple.EXPECT().GetPattern().Times(1).Return(pattern)
	mockTuple.EXPECT().GetHandler().Times(1).Return(
		func(bot interfaces.Bot, msg interfaces.Message, result []string) {
			assert.Equal(t, bot, mockBot)
			assert.Equal(t, msg, mockMsg)
			assert.Equal(t, result, []string{"Hello World"})

			ran <- true
		},
	)

	processMessage(mockBot, mockMsg)

	<-ran
}

func Test_processMessage_WithUnhandledIM(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockBot := mocks.NewMockBot(mockCtrl)
	mockMsg := mocks.NewMockMessage(mockCtrl)

	mockMsg.EXPECT().IsSelf().Times(1).Return(false)
	mockMsg.EXPECT().GetContent().Times(1).Return("Hello World")
	mockMsg.EXPECT().IsIM().Times(1).Return(true)

	mockBot.EXPECT().GetIMHandlers().Times(1).Return(
		[]interfaces.HandlerTuple{},
	)

	processMessage(mockBot, mockMsg)
}
