//********************************************************************************************************************//
//
// This file is part of golexoffice.
// All code may be used. Feel free and maybe code something better.
//
// Author: Jonas Kwiedor (aka gowizzard)
//
//********************************************************************************************************************//

package golexoffice

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"os"
)

// FileReturn is to decode json data
type FileReturn struct {
	Id string `json:"id"`
}

// AddFile is to upload a file
func (c *Config) AddFile(file *os.File, name string) (FileReturn, error) {

	// Create form data
	body := &bytes.Buffer{}

	// Create writer
	writer := multipart.NewWriter(body)

	// Create body data
	filePart, err := writer.CreateFormFile("file", name)
	if err != nil {
		return FileReturn{}, err
	}

	// Copy form & file
	_, err = io.Copy(filePart, file)
	if err != nil {
		return FileReturn{}, err
	}

	// Create text part
	_ = writer.WriteField("type", "voucher")

	// Close writer
	err = writer.Close()
	if err != nil {
		return FileReturn{}, err
	}

	// Set config for new request
	//c := NewConfig(, token, &http.Client{})

	// Send request
	response, err := c.Send("/v1/files/", body, "POST", writer.FormDataContentType())
	if err != nil {
		return FileReturn{}, err
	}

	// Close request
	defer response.Body.Close()

	// Decode data
	var decode FileReturn

	err = json.NewDecoder(response.Body).Decode(&decode)
	if err != nil {
		return FileReturn{}, err
	}

	// Return data
	return decode, nil

}
