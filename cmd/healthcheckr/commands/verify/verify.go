package verify

import (
	"healthcheckr/cmd/healthcheckr/commands/verify/http"

	"github.com/spf13/cobra"
)

var Command = cobra.Command{
	Use:     "verify",
	Aliases: []string{"check", "v", "c"},
	Short:   "Verifies different types of checks",
	RunE:    runE,
}

func init() {
	Command.AddCommand(http.GetCommand())
}

func GetCommand() *cobra.Command {
	return &Command
}

func runE(command *cobra.Command, args []string) error {
	return command.Help()
}
