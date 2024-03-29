package http

import (
	"fmt"
	"net/http"

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
	RunE:    runE,
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

	logrus.Infof("starting server on 0.0.0.0:8080")
	http.ListenAndServe("0.0.0.0:8080", mux)
	return nil
}
