package main

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
)

// Telegram constants
const (
	// APIEndpoint is the endpoint for all API methods,
	// with formatting for Sprintf.
	APIEndpoint = "https://api.telegram.org/bot%s/%s"
	// FileEndpoint is the endpoint for downloading a file from Telegram.
	FileEndpoint = "https://api.telegram.org/file/bot%s/%s"
)

func (p *Plugin) makeRequest(req *http.Request) ([]byte, error) {

	log.Debugf("[http-client-tgbot] %+v", req)

	reg := regexp.MustCompile(`^/(file/)?bot[0-9]+:[a-zA-Z\d_]+/(.*)$`)
	match := reg.FindStringSubmatch(req.URL.EscapedPath())
	var endpoint string
	switch len(match) {
	case 3:
		if match[1] == "" {
			endpoint = fmt.Sprintf(APIEndpoint, p.config.Token, match[2])
		} else {
			endpoint = fmt.Sprintf(FileEndpoint, p.config.Token, match[2])
		}
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
	fmt.Println(req.URL)
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}

	// read response
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Debugf("[http-client-tgbot] %s", string(data))

	return data, nil
}
