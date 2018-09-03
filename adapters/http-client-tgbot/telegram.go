package main

import (
	"io/ioutil"
	"net/http"
	"encoding/json"
	"fmt"
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

func makeRequest(req *http.Request) []byte {
	//reg := regexp.MustCompile(`(^https://api.telegram.org/)(file/)?bot[0-9]+:[a-zA-Z\d]+(/[^/]+$)`)

	b, _ := json.Marshal(req)
	fmt.Print(string(b))

	resp, err := client.Do(req)
	if err != nil {
		return nil
	}

	// read response
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	return data
}