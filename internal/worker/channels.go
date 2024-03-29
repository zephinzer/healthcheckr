package worker

import (
	"fmt"
	"os"
	"strings"
)

type Channels []Channel

type Channel struct {
	Name   *string      `json:"name" yaml:"name"`
	Type   *string      `json:"type" yaml:"type"`
	ApiKey ChannelValue `json:"apiKey" yaml:"apiKey"`
	Url    ChannelValue `json:"url" yaml:"url"`
	ChatId ChannelValue `json:"chatId" yaml:"chatId"`
}

func (c *Channel) Validate() error {
	isNameDefined := c.Name != nil
	isTypeDefined := c.Type != nil
	isApiKeyDefined := c.ApiKey.IsDefined()
	isChatIdDefined := c.ChatId.IsDefined()
	isUrlDefined := c.Url.IsDefined()

	channelIdentifier := ""
	channelInfo := []string{}
	if isNameDefined {
		channelInfo = append(channelInfo, fmt.Sprintf("name[%s]", *c.Name))
	}
	if isTypeDefined {
		channelInfo = append(channelInfo, fmt.Sprintf("type[%s]", *c.Type))
	}
	if len(channelInfo) == 0 {
		channelIdentifier = "[unknown]"
	} else {
		channelIdentifier = "{" + strings.Join(channelInfo, ",") + "}"
	}

	if !isNameDefined {
		return fmt.Errorf("failed to receive channel name for channel[%s]", channelIdentifier)
	}
	if !isTypeDefined {
		return fmt.Errorf("failed to receive channel type for channel[%s]", channelIdentifier)
	}

	switch *c.Type {
	case ChannelTelegram:
		if !isApiKeyDefined {
			return fmt.Errorf("failed to receive bot token for channel[%s]", *c.Name)
		}
		if !isChatIdDefined {
			return fmt.Errorf("failed to receive chat id for channel[%s]", *c.Name)
		}
	case ChannelSlack:
		if !isUrlDefined {
			return fmt.Errorf("failed to receive slack webhook url for channel[%s]", *c.Name)
		}
	}
	return nil
}

type ChannelValue struct {
	Value   *string `json:"value" yaml:"value"`
	FromEnv *string `json:"fromEnv" yaml:"fromEnv"`
}

func (cv *ChannelValue) Get() (string, error) {
	if cv.Value != nil {
		return *cv.Value, nil
	}
	if cv.FromEnv != nil {
		value, isDefined := os.LookupEnv(*cv.FromEnv)
		if !isDefined {
			return "", fmt.Errorf("failed to retrieve envvar[%s]", *cv.FromEnv)
		}
		return value, nil
	}
	return "", fmt.Errorf("failed unexpectedly to Get() a ChannelValue")
}

func (cv *ChannelValue) IsDefined() bool {
	return cv.Value != nil || cv.FromEnv != nil
}
