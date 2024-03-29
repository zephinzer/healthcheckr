package debug

import (
	"healthcheckr/cmd/healthcheckr/commands/debug/http"
	"healthcheckr/cmd/healthcheckr/commands/debug/telegram"

	"github.com/spf13/cobra"
)

var Command = cobra.Command{
	Use:     "debug",
	Aliases: []string{"dbg", "d"},
	Short:   "Debug/utilities invocations",
	RunE:    runE,
}

func init() {
	Command.AddCommand(http.GetCommand())
	Command.AddCommand(telegram.GetCommand())
}

func GetCommand() *cobra.Command {
	return &Command
}

func runE(command *cobra.Command, args []string) error {
	return command.Help()
}
