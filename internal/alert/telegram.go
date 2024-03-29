package alert

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/sirupsen/logrus"
)

func NewTelegramChannel(
	botToken string,
	chatId string,
) (Channel, error) {
	botInstance, err := bot.New(botToken)
	if err != nil {
		return nil, fmt.Errorf("failed to initialise a telegram bot: %s", err)
	}
	return &telegramChannel{
		BotInstance: botInstance,
		ChatId:      chatId,
	}, nil
}

type telegramChannel struct {
	BotInstance *bot.Bot
	ChatId      string
}

func (tgc *telegramChannel) GetType() string { return "telegram" }

func (tgc *telegramChannel) Send(message Message) error {
	sendContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	messageTitle := message.Title
	messageTitle = strings.ReplaceAll(messageTitle, "-", "\\-")
	messageTitle = strings.ReplaceAll(messageTitle, ".", "\\.")
	messageTitle = strings.ReplaceAll(messageTitle, "=", "\\=")
	messageTitle = addMessageTitleEmoji(messageTitle, message.Type)
	messageBody := message.Body
	messageBody = strings.ReplaceAll(messageBody, "-", "\\-")
	messageBody = strings.ReplaceAll(messageBody, ".", "\\.")
	messageBody = strings.ReplaceAll(messageBody, "=", "\\=")
	messageDetails := message.Details
	messageDetails = strings.ReplaceAll(messageDetails, "-", "\\-")
	messageDetails = strings.ReplaceAll(messageDetails, ".", "\\.")
	messageDetails = strings.ReplaceAll(messageDetails, "=", "\\=")

	normalisedMessage := fmt.Sprintf("*%s*\n\n%s\n\n```\n%s\n```", messageTitle, messageBody, messageDetails)

	msg, err := tgc.BotInstance.SendMessage(
		sendContext,
		&bot.SendMessageParams{
			ChatID:    tgc.ChatId,
			Text:      normalisedMessage,
			ParseMode: models.ParseModeMarkdown,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to send message to chat[%v]: %s", tgc.ChatId, err)
	}
	logrus.Debugf("sent message[%v] to chat[%s]", msg.ID, tgc.ChatId)
	return nil
}
