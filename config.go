//**********************************************************
//
// This file is part of lexoffice.
// All code may be used. Feel free and maybe code something better.
//
// Author: Jonas Kwiedor
//
//**********************************************************

package golexoffice

import (
	"io"
	"net/http"
)

const (
	baseURL = "https://api.lexoffice.io"
)

// Config is to define the request data
type Config struct {
	token   string
	baseUrl string
	client  *http.Client
}

func NewConfig(token string, httpClient *http.Client) *Config {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	return &Config{
		token:  token,
		client: httpClient,
	}
}

func (c *Config) SetBaseUrl(url string) {
	c.baseUrl = url
}

// Send is to send a new request
func (c *Config) Send(path string, body io.Reader, method, contentType string) (*http.Response, error) {

	// Set url
	var url string
	if c.baseUrl != "" {
		url = c.baseUrl + path
	} else {
		url = baseURL + path
	}

	// Request
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	// Define header
	request.Header.Set("Authorization", "Bearer "+c.token)
	request.Header.Set("Content-Type", contentType)
	request.Header.Set("Accept", "application/json")

	// Send request & get response
	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}

	// Return data
	return response, nil

}
