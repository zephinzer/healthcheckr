package http

import (
	"healthcheckr/internal/check"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Command = cobra.Command{
	Use:     "http",
	Aliases: []string{"h", "url", "u"},
	RunE:    runE,
}

func init() {
	Command.Flags().StringArrayP("expect-body-regex", "R", []string{}, "Regular expression to match in the response body")
	Command.Flags().IntP("expect-status-code", "C", 200, "Status code to expect")
	Command.Flags().IntP("expect-response-time-ms", "t", 5000, "Maximum number of milliseconds before target has to respond")
	Command.Flags().StringP("use-user-agent", "a", "healthcheckr/1.0", "User-Agent to use in connection call")
	Command.Flags().StringP("use-scheme", "s", "https", "Schema to use for the connection")
	Command.Flags().StringP("use-hostname", "H", "google.com", "Hostname to hit, include the port number if applicable")
	Command.Flags().StringP("use-method", "X", http.MethodGet, "Method to use for the connection")
	Command.Flags().StringP("use-path", "p", "/", "Path to use")
	Command.Flags().StringArrayP("use-query", "q", []string{}, "Query parameters to use")

	viper.BindPFlags(Command.Flags())
}

func GetCommand() *cobra.Command {
	return &Command
}

func runE(command *cobra.Command, args []string) error {
	scheme := viper.GetString("use-scheme")
	hostname := viper.GetString("use-hostname")
	path := viper.GetString("use-path")
	queries := viper.GetStringSlice("use-query")
	method := viper.GetString("use-method")
	userAgent := viper.GetString("use-user-agent")
	expectedResponseTimeMs := viper.GetInt("expect-response-time-ms")
	timeoutDuration := time.Millisecond * time.Duration(expectedResponseTimeMs)

	expectedBodyRegexes := viper.GetStringSlice("expect-body-regex")
	expectedStatusCode := viper.GetInt("expect-status-code")

	if err := check.HttpBasedUrl(check.HttpBasedUrlOpts{
		Scheme:    scheme,
		Hostname:  hostname,
		Method:    method,
		Path:      path,
		Queries:   queries,
		UserAgent: userAgent,

		ExpectedStatusCode:  expectedStatusCode,
		ExpectedBodyRegexes: expectedBodyRegexes,

		Timeout: timeoutDuration,
	}); err != nil {
		return err
	}

	return nil
}
