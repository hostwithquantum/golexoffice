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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
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

	if isSuccessful(response) {
		// Return data
		return response, nil
	}

	// TODO(till): revisit parsing when we add more API endpoints
	if strings.Contains(url, "invoices") {
		return nil, parseErrorResponse(response)
	}

	return nil, parseLegacyErrorResponse(response)
}

func isSuccessful(response *http.Response) bool {
	return response.StatusCode < 400
}

func parseErrorResponse(response *http.Response) error {
	var errorResp ErrorResponse
	err := json.NewDecoder(response.Body).Decode(&errorResp)
	if err != nil {
		return fmt.Errorf("decoding error while unpacking response: %s", err)
	}
	defer response.Body.Close()

	var keep []error
	for _, detail := range errorResp.Details {
		keep = append(keep, fmt.Errorf(
			"field: %s (%s): %s", detail.Field, detail.Violation, detail.Message,
		))
	}
	if len(keep) == 0 {
		return fmt.Errorf("error: %s (%d %s)", errorResp.Message, errorResp.Status, errorResp.Error)
	}

	return errors.Join(keep...)
}

func parseLegacyErrorResponse(response *http.Response) error {
	var errorResp LegacyErrorResponse
	err := json.NewDecoder(response.Body).Decode(&errorResp)
	if err != nil {
		return fmt.Errorf("decoding error while unpacking response: %s", err)
	}
	defer response.Body.Close()

	// potentially multiple issues returned from the LexOffice API
	var keep []error
	for _, issue := range errorResp.IssueList {
		keep = append(keep, fmt.Errorf("key: %s (%s): %s", issue.Key, issue.Source, issue.Type))
	}
	if len(keep) == 0 {
		return fmt.Errorf("something went wrong but unclear what (empty IssueList)")
	}
	return errors.Join(keep...)
}
