package client

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/go-resty/resty/v2"
)

func Build(URL, productId string) (*resty.Client, error) {
	u, err := url.ParseRequestURI(URL)

	if err != nil {
		return nil, err
	}

	baseUrl := fmt.Sprintf("%s://%s", u.Scheme, u.Host)

	restyBase := resty.New().
		SetBaseURL(baseUrl).
		SetDebug(strings.ToLower(os.Getenv("TF_LOG")) == "debug").
		OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
			tfLogLevel := strings.ToLower(os.Getenv("TF_LOG"))
			if slices.Contains([]string{"debug", "trace"}, tfLogLevel) {
				r.SetDebug(true)
			}
			return nil
		}).
		OnRequestLog(func(log *resty.RequestLog) error {
			// Never log auth token
			if log.Header.Get("Authorization") != "" {
				log.Header.Set("Authorization", "<REDACTED>")
			}
			return nil
		}).
		SetHeader("content-type", "application/json").
		SetHeader("accept", "*/*").
		SetHeader("user-agent", "jfrog/"+productId).
		SetRetryCount(20)

	restyBase.DisableWarn = true

	return restyBase, nil
}

func AddAuth(client *resty.Client, apiKey, accessToken string) (*resty.Client, error) {
	if accessToken != "" {
		return client.SetAuthToken(accessToken), nil
	}
	if apiKey != "" {
		return client.SetHeader("X-JFrog-Art-Api", apiKey), nil
	}
	return nil, fmt.Errorf("no authentication details supplied")
}

var mergeAndSaveRegex = regexp.MustCompile(".*Could not merge and save new descriptor.*")

func RetryOnMergeError(response *resty.Response, _r error) bool {
	return mergeAndSaveRegex.MatchString(string(response.Body()[:]))
}

func NeverRetry(response *resty.Response, err error) bool {
	return false
}
