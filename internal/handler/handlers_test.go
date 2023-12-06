package handler

import (
	"context"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestHandler_HandleStart(t *testing.T) {
	h := &Handler{
		log: logrus.New(),
		ctx: context.Background(),
	}

	message := tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: 123},
	}

	msg := h.HandleStart(&message)

	assert.Equal(t, "Press menu button to see command list", msg.Text)
}
