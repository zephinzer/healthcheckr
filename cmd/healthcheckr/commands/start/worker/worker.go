package worker

import (
	"encoding/json"
	"fmt"
	"healthcheckr/internal/alert"
	"healthcheckr/internal/check"
	"healthcheckr/internal/scheduler"
	"healthcheckr/internal/worker"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Command = cobra.Command{
	Use:     "worker",
	Aliases: []string{"w"},
	RunE:    runE,
}

func init() {
	Command.Flags().StringP("config-path", "c", "./config.yaml", "Path to configuration file")
	Command.Flags().StringP("server-addr", "a", "0.0.0.0", "Network interface address for healthcheck server to bind to")
	Command.Flags().IntP("server-port", "p", 8080, "Port to expose for healthcheck server")

	viper.BindPFlags(Command.Flags())
}

func GetCommand() *cobra.Command {
	return &Command
}

func runE(command *cobra.Command, args []string) error {
	serverAddr := viper.GetString("server-addr")
	serverPort := viper.GetInt("server-port")
	serverBindAddr := fmt.Sprintf("%s:%v", serverAddr, serverPort)

	configurationFilePath := viper.GetString("config-path")
	configurationData, err := worker.LoadConfigurationFromPath(configurationFilePath)
	if err != nil {
		return fmt.Errorf("failed to load configuration from path[%s]: %s", err, configurationFilePath)
	}
	o, _ := json.MarshalIndent(configurationData, "", "  ")
	logrus.Debugf("loaded configuration as follows:\n%s", string(o))

	channelMap := map[string]alert.Channel{}

	for _, channel := range configurationData.Channels {
		if err := channel.Validate(); err != nil {
			return fmt.Errorf("failed to create channel: %s", err)
		}
		switch *channel.Type {
		case worker.ChannelTelegram:
			botToken, err := channel.ApiKey.Get()
			if err != nil {
				return fmt.Errorf("failed to get the telegram bot token: %s", err)
			}
			chatId, err := channel.ChatId.Get()
			if err != nil {
				return fmt.Errorf("failed to get telegram chat id: %s", err)
			}
			channelInstance, err := alert.NewTelegramChannel(
				botToken,
				chatId,
			)
			if err != nil {
				return fmt.Errorf("failed to create telegram channel: %s", err)
			}
			channelMap[*channel.Name] = channelInstance
		case worker.ChannelSlack:
			webhookUrl, err := channel.Url.Get()
			if err != nil {
				return fmt.Errorf("failed to get the slack webhook url: %s", err)
			}
			channelInstance, err := alert.NewSlackWebhookChannel(webhookUrl)
			if err != nil {
				return fmt.Errorf("failed to create slack webhook channel: %s", err)
			}
			channelMap[*channel.Name] = channelInstance
		}
	}

	for _, httpCheck := range configurationData.Http {
		channels := []alert.Channel{}
		for _, channel := range httpCheck.Channels {
			logrus.Debugf("mapping channel[%s]", channel)
			if channelInstance, ok := channelMap[channel]; ok {
				channels = append(channels, channelInstance)
			} else {
				return fmt.Errorf("failed to identify channel[%s], has it been added to the root-level channels property?", channel)
			}
		}

		logrus.Debugf("scheduling http check with %v channel(s)...", len(channels))
		scheduler.ScheduleHttp(
			check.HttpBasedUrlOpts{
				Scheme:    httpCheck.Scheme,
				Hostname:  httpCheck.Hostname,
				Method:    httpCheck.Method,
				Path:      httpCheck.Path,
				Queries:   httpCheck.Queries,
				UserAgent: httpCheck.UserAgent,

				ExpectedStatusCode:  int(httpCheck.ExpectStatusCode),
				ExpectedBodyRegexes: httpCheck.ExpectBodyRegexes,

				Timeout: time.Duration(httpCheck.TimeoutMs) * time.Millisecond,
			},
			scheduler.ScheduleOpts{
				Interval:             time.Duration(httpCheck.IntervalMs) * time.Millisecond,
				FailureThreshold:     httpCheck.FailureThreshold,
				AlertMinimumInterval: time.Duration(httpCheck.AlertMinimumIntervalS) * time.Second,
			},
			scheduler.ChannelOpts{
				Channels: channels,
			},
		)
	}

	mux := http.NewServeMux()

	logrus.Infof("starting healthcheckr worker listening on addr[%s]...", serverBindAddr)
	if err := http.ListenAndServe(serverBindAddr, mux); err != nil {
		return fmt.Errorf("failed to listen on interface[%s]: %s", serverBindAddr, err)
	}

	return command.Help()
}
