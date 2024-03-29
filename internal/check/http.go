package check

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"healthcheckr/internal/utils"

	"github.com/sirupsen/logrus"
)

type HttpBasedUrlOpts struct {
	Scheme    string
	Hostname  string
	Path      string
	Queries   []string
	Method    string
	UserAgent string

	ExpectedBodyRegexes []string
	ExpectedStatusCode  int

	Timeout time.Duration
}

func HttpBasedUrl(opts HttpBasedUrlOpts) error {
	urlInstance, err := utils.MakeHttpUrl(opts.Scheme, opts.Hostname, opts.Path, opts.Queries)
	if err != nil {
		return fmt.Errorf("failed to make url: %s", err)
	}

	logrus.Infof("running url check on url[%s] using method[%s]...", urlInstance.String(), opts.Method)

	client := http.Client{}
	requestCtx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
	defer cancel()
	requestStartTime := time.Now()
	logrus.Debugf("creating request to url[%s]...", urlInstance.String())
	request, err := http.NewRequestWithContext(requestCtx, opts.Method, urlInstance.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %s", err)
	}
	request.Header.Set("User-Agent", opts.UserAgent)
	logrus.Debugf("executing request to url[%s]...", urlInstance.String())
	response, err := client.Do(request)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("failed to receive a response within timeout[%s]", opts.Timeout.String())
		}
		return fmt.Errorf("failed to execute request: %s", err)
	}
	requestDuration := time.Since(requestStartTime)
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %s", err)
	}
	logrus.Debugf("duration: %vms", requestDuration.Milliseconds())
	logrus.Debugf("response body size: %v", len(responseBody))
	logrus.Debugf("response status code: %v", response.StatusCode)
	logrus.Debugf("response headers:")
	responseHeaders := response.Header.Clone()
	for responseHeaderKey, responseHeaderValue := range responseHeaders {
		logrus.Debugf("  %s: [%v]{'%s'}", responseHeaderKey, len(responseHeaderValue), strings.Join(responseHeaderValue, "', '"))
	}

	if response.StatusCode != opts.ExpectedStatusCode {
		return fmt.Errorf("failed to get expected statusCode[%v], got statusCode[%v] instead", opts.ExpectedStatusCode, response.StatusCode)
	}

	expectedMatchErrors := []string{}
	if len(opts.ExpectedBodyRegexes) > 0 {
		for _, expectedBodyRegex := range opts.ExpectedBodyRegexes {
			logrus.Debugf("verifying regex[%s] in response body...", expectedBodyRegex)
			matcher := regexp.MustCompile(expectedBodyRegex)
			if !matcher.Match(responseBody) {
				expectedMatchErrors = append(
					expectedMatchErrors,
					fmt.Sprintf("failed to match regexp[%s] in response body", expectedBodyRegex),
				)
			}
		}
	}
	if len(expectedMatchErrors) > 0 {
		return fmt.Errorf("failed to match %v expected regular expressions: ['%s']", len(expectedMatchErrors), strings.Join(expectedMatchErrors, "', '"))
	}

	logrus.Infof("successfully received statusCode[%v] from url[%s] after %vms", opts.ExpectedStatusCode, urlInstance.String(), requestDuration.Milliseconds())

	return nil
}
