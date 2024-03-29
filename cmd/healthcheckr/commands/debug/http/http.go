package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	StatusFailure = http.StatusInternalServerError
	StatusSuccess = http.StatusOK
)

var Command = cobra.Command{
	Use:     "http",
	Aliases: []string{"server", "h"},
	Short:   "Starts a utility HTTP server for running checks against",
	Long: strings.Trim(`
Starts a utility HTTP server that returns 200 or 500 on the / endpoint
depending on the mode of operations which can be toggled by triggering
a request to the /mode endpoint.
`, "\n "),
	RunE: runE,
}

func GetCommand() *cobra.Command {
	return &Command
}

func runE(command *cobra.Command, args []string) error {
	isFailModeOn := false
	currentStatus := StatusSuccess
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logrus.Infof("returning status code %v", currentStatus)
		w.WriteHeader(currentStatus)
		w.Write([]byte(fmt.Sprintf(`{"status":%v}`, currentStatus)))
	})
	mux.HandleFunc("/mode", func(w http.ResponseWriter, r *http.Request) {
		isFailModeOn = !isFailModeOn
		if isFailModeOn {
			currentStatus = StatusFailure
		} else {
			currentStatus = StatusSuccess
		}
		logrus.Infof("fail mode toggled, / will now return %v", currentStatus)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	serverAddr := "0.0.0.0:8080"
	logrus.Infof("starting server on %s", serverAddr)
	if err := http.ListenAndServe(serverAddr, mux); err != nil {
		return fmt.Errorf("failed to start server: %s", err)
	}
	return nil
}
