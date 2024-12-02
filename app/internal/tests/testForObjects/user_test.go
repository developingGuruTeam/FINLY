package testForObjects

import (
	"cachManagerApp/app/internal/tests/mocks"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"testing"
)

var Update = tgbotapi.Update{
	Message: &tgbotapi.Message{
		Chat: &tgbotapi.Chat{
			ID: 12345,
		},
		From: &tgbotapi.User{
			UserName: "TestUser",
		},
	},
}

func TestUsersHandlers(t *testing.T) {
	mocker := new(mocks.UsersHandlers)

	t.Run("TestPostUser", func(t *testing.T) {
		mocker.On("PostUser", Update).Return(nil)
		err := mocker.PostUser(Update)
		assert.NoError(t, err)
		mocker.AssertCalled(t, "PostUser", Update)
	},
	)

}
