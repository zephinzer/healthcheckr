package utils

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
)

func MakeHttpUrl(scheme, hostname, path string, queries []string) (*url.URL, error) {
	fullUrl := fmt.Sprintf("%s://%s%s", scheme, hostname, path)
	urlInstance, err := url.Parse(fullUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse composed url[%s]: %s", fullUrl, err)
	}
	if len(queries) > 0 {
		q := urlInstance.Query()
		for _, query := range queries {
			querySlice := strings.SplitN(query, "=", 2)
			if len(querySlice) == 1 {
				logrus.Warnf("found unexpected query parameter ('%s') with no value assignment", querySlice[0])
				q.Add(querySlice[0], "")
			} else {
				q.Add(querySlice[0], querySlice[1])
			}
		}
		urlInstance.RawQuery = q.Encode()
	}
	return urlInstance, nil
}
