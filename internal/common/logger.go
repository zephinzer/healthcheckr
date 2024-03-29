package common

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	ConfigLogLevelDescription = "Level to log at (0-5)"
	ConfigLogLevelDefault     = 3
	ConfigLogLevelKey         = "log-level"
	ConfigLogLevelKeyAbbrv    = "l"
)

func AddCobraFlag(command *cobra.Command) {
	command.Flags().IntP(ConfigLogLevelKey, ConfigLogLevelKeyAbbrv, ConfigLogLevelDefault, ConfigLogLevelDescription)
}

func AddCobraPersistentFlag(command *cobra.Command) {
	command.PersistentFlags().IntP(ConfigLogLevelKey, ConfigLogLevelKeyAbbrv, ConfigLogLevelDefault, ConfigLogLevelDescription)
}

func InitLogLevel() {
	logLevel := viper.GetInt(ConfigLogLevelKey)
	switch logLevel {
	case 0:
		logrus.SetLevel(logrus.FatalLevel)
	case 1:
		logrus.SetLevel(logrus.ErrorLevel)
	case 2:
		logrus.SetLevel(logrus.WarnLevel)
	case 3:
		logrus.SetLevel(logrus.InfoLevel)
	case 4:
		logrus.SetLevel(logrus.DebugLevel)
	default:
		logrus.SetLevel(logrus.TraceLevel)
	}
}
