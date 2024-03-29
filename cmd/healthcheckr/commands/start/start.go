package start

import (
	"healthcheckr/cmd/healthcheckr/commands/start/worker"

	"github.com/spf13/cobra"
)

var Command = cobra.Command{
	Use:     "start",
	Aliases: []string{"s", "run"},
	RunE:    runE,
}

func init() {
	Command.AddCommand(worker.GetCommand())
}

func GetCommand() *cobra.Command {
	return &Command
}

func runE(command *cobra.Command, args []string) error {
	return command.Help()
}
