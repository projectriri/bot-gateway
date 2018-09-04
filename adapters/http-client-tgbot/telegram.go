package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"github.com/pkg/errors"
	"net/url"
	log "github.com/sirupsen/logrus"
)

// Telegram constants
const (
	// APIEndpoint is the endpoint for all API methods,
	// with formatting for Sprintf.
	APIEndpoint = "https://api.telegram.org/bot%s/%s"
	// FileEndpoint is the endpoint for downloading a file from Telegram.
	FileEndpoint = "https://api.telegram.org/file/bot%s/%s"
)

var client http.Client

func makeRequest(req *http.Request) ([]byte, error) {

	// DEBUG
	fmt.Printf("%+v", req)

	reg := regexp.MustCompile(`^/(file/)?bot[0-9]+:[a-zA-Z\d_]+/(.*)$`)
	match := reg.FindStringSubmatch(req.URL.EscapedPath())
	var endpoint string
	switch len(match) {
	case 3:
		endpoint = fmt.Sprintf(FileEndpoint, config.Token, match[2])
	case 2:
		endpoint = fmt.Sprintf(APIEndpoint, config.Token, match[1])
	default:
		log.Errorf("[http-client-tgbot] bad url endpoint: %v", req.URL.EscapedPath())
		return nil, errors.New("unknown url scheme")
	}

	var err error
	req.Host = "api.telegram.org"
	req.URL, err = url.Parse(endpoint)
	if err != nil {
		log.Errorf("[http-client-tgbot] fail to parse url : %v", endpoint)
		return nil, errors.New("fail to parse url")
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// read response
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
