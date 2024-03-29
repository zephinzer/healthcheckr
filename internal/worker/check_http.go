package worker

import (
	"fmt"
	"net/http"
)

const (
	DefaultHttpMethod                = http.MethodGet
	DefaultHttpPath                  = "/"
	DefaultHttpUserAgent             = "healthcheckr/1.0"
	DefaultHttpTimeoutMs             = 5000
	DefaultHttpIntervalMs            = 5000
	DefaultHttpExpectedStatusCode    = http.StatusOK
	DefaultHttpAlertMinimumIntervalS = 60
)

type HttpChecks []HttpCheck

type HttpCheck struct {
	Scheme   string   `json:"scheme" yaml:"scheme"`
	Hostname string   `json:"hostname" yaml:"hostname"`
	Path     string   `json:"path" yaml:"path"`
	Queries  []string `json:"queries" yaml:"queries"`
	Method   string   `json:"method" yaml:"method"`

	UserAgent string `json:"userAgent" yaml:"userAgent"`
	TimeoutMs int    `json:"timeoutMs" yaml:"timeoutMs"`

	ExpectStatusCode  int      `json:"expectStatusCode" yaml:"expectStatusCode"`
	ExpectBodyRegexes []string `json:"expectBodyRegexes" yaml:"expectBodyRegexes"`

	IntervalMs            int `json:"intervalMs" yaml:"intervalMs"`
	FailureThreshold      int `json:"failureThreshold" yaml:"failureThreshold"`
	AlertMinimumIntervalS int `json:"alertMinimumIntervalS" yaml:"alertMinimumIntervalS"`

	Channels []string `json:"channels"`
}

func (h *HttpCheck) InitWithDefaults() error {
	h.Scheme = initString(&h.Scheme, "https")
	if h.Hostname == "" {
		return fmt.Errorf("failed to receive a hostname")
	}
	h.Path = initString(&h.Path, DefaultHttpPath)
	h.Method = initString(&h.Method, DefaultHttpMethod)
	h.Queries = initStringSlice(&h.Queries, []string{})
	h.UserAgent = initString(&h.UserAgent, DefaultHttpUserAgent)
	h.TimeoutMs = initInt(&h.TimeoutMs, DefaultHttpTimeoutMs)
	h.ExpectBodyRegexes = initStringSlice(&h.ExpectBodyRegexes, []string{})
	h.ExpectStatusCode = initInt(&h.ExpectStatusCode, DefaultHttpExpectedStatusCode)
	h.IntervalMs = initInt(&h.IntervalMs, DefaultHttpIntervalMs)
	h.AlertMinimumIntervalS = initInt(&h.AlertMinimumIntervalS, DefaultHttpAlertMinimumIntervalS)

	return nil
}
