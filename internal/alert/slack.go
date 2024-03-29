package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func NewSlackWebhookChannel(
	webhookUrl string,
) (Channel, error) {
	return &slackWebhookChannel{
		WebhookUrl: webhookUrl,
	}, nil
}

type slackWebhookChannel struct {
	WebhookUrl string
}

func (sc *slackWebhookChannel) GetType() string { return "slack" }

func (sc *slackWebhookChannel) Send(message Message) error {
	client := http.Client{
		Timeout: time.Second * 5,
	}
	// format reference - https://api.slack.com/messaging/webhooks
	messageTitle := message.Title
	messageTitle = addMessageTitleEmoji(messageTitle, message.Type)
	messageBody := message.Body
	messageDetails := message.Details

	messageBlocks := []map[string]any{
		{
			"type": "section",
			"text": map[string]any{
				"type": "mrkdwn",
				"text": "*" + messageTitle + "*",
			},
		},
	}
	if len(messageBody) > 0 {
		messageBlocks = append(
			messageBlocks,
			map[string]any{
				"type": "section",
				"text": map[string]any{
					"type": "mrkdwn",
					"text": messageBody,
				},
			},
		)
	}
	if len(messageDetails) > 0 {
		messageBlocks = append(
			messageBlocks,
			map[string]any{
				"type": "section",
				"text": map[string]any{
					"type": "mrkdwn",
					"text": "```\n" + messageDetails + "\n```",
				},
			},
		)
	}

	messageData := map[string]any{
		"text":   messageTitle,
		"blocks": messageBlocks,
	}
	requestBody, err := json.Marshal(messageData)
	if err != nil {
		return fmt.Errorf("failed to marshal message data: %s", err)
	}

	messageBuffer := bytes.NewBuffer(requestBody)
	request, err := http.NewRequest(
		http.MethodPost,
		sc.WebhookUrl,
		messageBuffer,
	)
	if err != nil {
		return fmt.Errorf("failed to create http webhook request: %s", err)
	}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("failed to execute http webhook request: %s", err)
	}
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %s", err)
	}
	logrus.Debugf("received response: %s", string(responseBody))
	return nil
}
