package telegram

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Command = cobra.Command{
	Use:     "telegram",
	Aliases: []string{"tg", "t"},
	Short:   "Starts a telegram bot, use /init in-chat to get the chat ID",
	RunE:    runE,
}

func init() {
	Command.Flags().String("bot-token", "t", "Get this token from the @BotFather (https://t.me/BotFather)")

	viper.BindPFlags(Command.Flags())
}

func GetCommand() *cobra.Command {
	return &Command
}

func runE(command *cobra.Command, args []string) error {
	botToken := viper.GetString("bot-token")
	if botToken == "" {
		return fmt.Errorf("failed to receive a bot token")
	}

	botInstance, err := bot.New(botToken)
	if err != nil {
		return fmt.Errorf("failed to create telegram bot instance: %s", err)
	}
	botUser, err := botInstance.GetMe(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get telegram bot details: %s", err)
	}
	logrus.Infof("bot user id: %v", botUser.ID)
	logrus.Infof("bot username: %s", botUser.Username)
	logrus.Infof("bot url: https://t.me/%s", botUser.Username)
	logrus.Infof("add this bot to a chat and talk to it to get the chat id for use in channels")
	logrus.Infof("waiting for pings...")

	startCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	botInstance.RegisterHandler(
		bot.HandlerTypeMessageText,
		"/init",
		bot.MatchTypePrefix,
		func(ctx context.Context, b *bot.Bot, update *models.Update) {
			chatId := update.Message.Chat.ID
			msg, err := b.SendMessage(
				ctx,
				&bot.SendMessageParams{
					ChatID:    chatId,
					Text:      fmt.Sprintf("chat id: `%v`", chatId),
					ParseMode: models.ParseModeMarkdown,
				},
			)
			if err != nil {
				logrus.Errorf("failed to respond to a message with the chat id: %s", err)
			}
			logrus.Infof("successfully responded to chat[%v] with message[%v]", chatId, msg.ID)
		},
	)
	botInstance.Start(startCtx)
	return nil
}
