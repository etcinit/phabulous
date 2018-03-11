package utilities

import (
	"testing"

	"github.com/etcinit/phabulous/app/interfaces"
	"github.com/etcinit/phabulous/app/testing/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_HasModule(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockBot := mocks.NewMockBot(mockCtrl)
	moduleA := mocks.NewMockModule(mockCtrl)
	moduleB := mocks.NewMockModule(mockCtrl)

	mockBot.EXPECT().GetModules().Times(2).Return([]interfaces.Module{
		moduleA,
		moduleB,
	})

	moduleA.EXPECT().GetName().Times(2).Return("module-a")
	moduleB.EXPECT().GetName().Times(1).Return("module-b")

	assert.True(t, HasModule(mockBot, "module-a"))
	assert.False(t, HasModule(mockBot, "module-c"))
}
