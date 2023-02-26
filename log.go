// Package gt ...
package gt

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"k8s.io/klog/v2"
)

const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

const (
	tolerateTime = 5 * time.Second
)

var (
	color = reset
)

func NewLogTrace() http.RoundTripper {
	return &logTrace{
		delegatedRoundTripper: &http.Transport{},
	}
}

type logTrace struct {
	delegatedRoundTripper http.RoundTripper
}

func (l *logTrace) RoundTrip(request *http.Request) (*http.Response, error) {
	start := time.Now()
	var requestBodyStr string
	var responseBodyStr string
	if request.Body != nil {
		requestBody, _ := ioutil.ReadAll(request.Body)
		requestBodyStr = string(requestBody)
		request.Body = ioutil.NopCloser(bytes.NewReader(requestBody))
	}

	response, err := l.delegatedRoundTripper.RoundTrip(request)

	if err != nil {
		return response, err
	}

	if response != nil {
		responseBody, _ := ioutil.ReadAll(response.Body)
		responseBodyStr = string(responseBody)
		response.Body = ioutil.NopCloser(bytes.NewReader(responseBody))
	}

	log := fmt.Sprintf("\n%s %s %s \n", request.Host, request.Method, request.URL)
	log += fmt.Sprintf("RequestBody: %s \n", requestBodyStr)
	log += fmt.Sprintf("Response: %s \n", responseBodyStr)

	cost := time.Now().Sub(start)
	if cost > tolerateTime {
		color = red
	}

	log += fmt.Sprintf("Cost time: %s%s\033[0m\n", color, time.Now().Sub(start))

	klog.Info(log)

	return response, err
}
