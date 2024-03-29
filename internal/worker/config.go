package worker

import (
	"fmt"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
)

func LoadConfigurationFromPath(configurationFilePath string) (*Configuration, error) {
	if !path.IsAbs(configurationFilePath) {
		currentWorkingDirectory, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get working directory: %s", err)
		}
		configurationFilePath = path.Join(currentWorkingDirectory, configurationFilePath)
	}

	configurationFileData, err := os.ReadFile(configurationFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file at path[%s]: %s", configurationFilePath, err)
	}

	configurationData := &Configuration{}
	if err := yaml.Unmarshal(configurationFileData, configurationData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal configuration file from yaml: %s", err)
	}
	if err := configurationData.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate configuration file: %s", err)
	}

	return configurationData, nil
}

type Configuration struct {
	Http     HttpChecks `json:"http" yaml:"http"`
	Channels Channels   `json:"channels" yaml:"channels"`
}

func (c *Configuration) Validate() error {
	// validate and init http check specifications
	for i := 0; i < len(c.Http); i++ { // no range here because we need to update the default values
		if err := c.Http[i].InitWithDefaults(); err != nil {
			return fmt.Errorf("failed to initialise configuration: %s", err)
		}
	}

	// validate and init channel specifications
	channelNames := map[string]bool{}
	for channelIndex, channel := range c.Channels {
		if err := channel.Validate(); err != nil {
			return fmt.Errorf("failed to initialise notification channels: %s", err)
		}
		channelName := *channel.Name
		if defined, exists := channelNames[channelName]; defined && exists {
			return fmt.Errorf("failed to get a unique name for channel[%s] at index[%v]", channelName, channelIndex)
		}
		channelNames[channelName] = true
	}

	// verify that specified channels in http checks exist in the list of channels
	for _, httpCheck := range c.Http {
		for _, channelName := range httpCheck.Channels {
			if defined, exists := channelNames[channelName]; !exists || !defined {
				existingChannelNames := []string{}
				for existingChannelName := range channelNames {
					existingChannelNames = append(existingChannelNames, existingChannelName)
				}
				return fmt.Errorf("failed to find channel[%s] from available channels['%s']", channelName, strings.Join(existingChannelNames, "', '"))
			}
		}
	}
	return nil
}
