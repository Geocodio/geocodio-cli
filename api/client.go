package api

import (
	"bytes"
	"fmt"
	"github.com/geocodio/geocodio-cli/release"
	"github.com/urfave/cli/v2"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

var apiVersion string = "v1.7"

func Request(method string, path string, c *cli.Context) ([]byte, error) {
	hostname := c.String("hostname")
	apiKey := c.String("apikey")
	url := fmt.Sprintf("https://%s/%s/%s?api_key=%s", hostname, apiVersion, path, url.QueryEscape(apiKey))

	client := http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", buildUserAgent())

	res, getErr := client.Do(req)
	if getErr != nil {
		return nil, getErr
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, readErr
	}

	return body, nil
}

func Upload(file *os.File, direction string, format string, fields string, c *cli.Context) ([]byte, error) {
	hostname := c.String("hostname")
	apiKey := c.String("apikey")
	url := fmt.Sprintf("https://%s/%s/lists?api_key=%s", hostname, apiVersion, url.QueryEscape(apiKey))

	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)

	directionWriter, err := writer.CreateFormField("direction")
	if err != nil {
		return nil, err
	}
	directionWriter.Write([]byte(direction))

	formatWriter, err := writer.CreateFormField("format")
	if err != nil {
		return nil, err
	}
	formatWriter.Write([]byte(format))

	fieldsWriter, err := writer.CreateFormField("fields")
	if err != nil {
		return nil, err
	}
	fieldsWriter.Write([]byte(fields))

	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))

	if err != nil {
		return nil, err
	}

	io.Copy(part, file)
	writer.Close()

	client := http.Client{
		Timeout: time.Minute * 30,
	}

	req, err := http.NewRequest(http.MethodPost, url, requestBody)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Set("User-Agent", buildUserAgent())

	res, getErr := client.Do(req)
	if getErr != nil {
		return nil, getErr
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	responseBody, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, err
	}

	return responseBody, nil
}

func buildUserAgent() string {
	return fmt.Sprintf("geocodio-cli v%s", release.Version())
}
