package healthcheckr

import (
	"healthcheckr/cmd/healthcheckr/commands/debug"
	"healthcheckr/cmd/healthcheckr/commands/start"
	"healthcheckr/cmd/healthcheckr/commands/verify"
	"healthcheckr/internal/common"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Command = cobra.Command{
	Use:               "healthcheckr",
	PersistentPreRunE: persistentPreRunE,
	RunE:              runE,
}

func init() {
	Command.AddCommand(debug.GetCommand())
	Command.AddCommand(start.GetCommand())
	Command.AddCommand(verify.GetCommand())

	common.AddCobraPersistentFlag(&Command)
	viper.BindPFlags(Command.PersistentFlags())
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
}

func GetCommand() *cobra.Command {
	return &Command
}

func persistentPreRunE(command *cobra.Command, args []string) error {
	common.InitLogLevel()
	return nil
}

func runE(command *cobra.Command, args []string) error {
	return command.Help()
}
