package connectors

import (
	"regexp"
	"sync"
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
	mockMsg.EXPECT().HasUser().Times(1).Return(true)
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

// Test_processMessage_SkipDuplicates asserts that multiple matches of the same element are only handled once.
func Test_processMessage_SkipDuplicates(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockBot := mocks.NewMockBot(mockCtrl)
	mockMsg := mocks.NewMockMessage(mockCtrl)
	mockTuple := mocks.NewMockHandlerTuple(mockCtrl)

	mockMsg.EXPECT().IsSelf().Times(1).Return(false)
	mockMsg.EXPECT().HasUser().Times(1).Return(true)
	mockMsg.EXPECT().GetContent().Times(1).Return("We have T1 and T100 and T100 again.")
	mockMsg.EXPECT().IsIM().Times(1).Return(false)

	mockBot.EXPECT().GetHandlers().Times(1).Return([]interfaces.HandlerTuple{
		mockTuple,
	})

	ran := make(chan bool)

	pattern, _ := regexp.Compile("([T|D][0-9]{1,16})")
	mockTuple.EXPECT().GetPattern().Times(1).Return(pattern)
	// we expect only two calls to GetHandler, even though the regex matches three times, T100 exists twice.
	mockTuple.EXPECT().GetHandler().Times(2).Return(
		func(bot interfaces.Bot, msg interfaces.Message, result []string) {
			assert.Equal(t, bot, mockBot)
			assert.Equal(t, msg, mockMsg)
			assert.Contains(t, [][]string{{"T1"}, {"T100"}}, result)

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
	mockMsg.EXPECT().HasUser().Times(1).Return(true)
	mockMsg.EXPECT().GetContent().Times(1).Return("Hello World")
	mockMsg.EXPECT().IsIM().Times(1).Return(true)

	mockBot.EXPECT().GetIMHandlers().Times(1).Return([]interfaces.HandlerTuple{
		mockTuple,
	})

	var wg sync.WaitGroup

	pattern, _ := regexp.Compile("Hello World")
	mockTuple.EXPECT().GetPattern().Times(1).Return(pattern)
	mockTuple.EXPECT().GetHandler().Times(1).Return(
		func(bot interfaces.Bot, msg interfaces.Message, result []string) {
			assert.Equal(t, bot, mockBot)
			assert.Equal(t, msg, mockMsg)
			assert.Equal(t, result, []string{"Hello World"})

			wg.Done()
		},
	)

	wg.Add(1)

	processMessage(mockBot, mockMsg)

	wg.Wait()
}

func Test_processMessage_WithUnhandledIM(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	var wg sync.WaitGroup

	mockBot := mocks.NewMockBot(mockCtrl)
	mockMsg := mocks.NewMockMessage(mockCtrl)

	mockMsg.EXPECT().IsSelf().Times(1).Return(false)
	mockMsg.EXPECT().HasUser().Times(1).Return(true)
	mockMsg.EXPECT().GetContent().Times(1).Return("Hello World")
	mockMsg.EXPECT().IsIM().Times(1).Return(true)

	mockBot.EXPECT().GetIMHandlers().Times(1).Return(
		[]interfaces.HandlerTuple{},
	)

	wg.Add(1)
	mockBot.EXPECT().GetUsageHandler().Times(1).Return(
		func(b interfaces.Bot, m interfaces.Message, matches []string) {
			wg.Done()
		},
	)

	processMessage(mockBot, mockMsg)

	wg.Wait()
}
