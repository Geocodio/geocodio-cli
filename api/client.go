package api

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

var apiVersion string = "v1.7"

func Request(method string, path string, hostname string, apiKey string) []byte {
	url := fmt.Sprintf("https://%s/%s/%s?api_key=%s", hostname, apiVersion, path, url.QueryEscape(apiKey))

	client := http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", fmt.Sprintf("geocodio-cli"))

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	return body
}

func Upload(file *os.File, direction string, format string, hostname string, apiKey string) []byte {
	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)

	directionWriter, err := writer.CreateFormField("direction")
	if err != nil {
		log.Fatal(err)
	}
	directionWriter.Write([]byte(direction))

	formatWriter, err := writer.CreateFormField("format")
	if err != nil {
		log.Fatal(err)
	}
	formatWriter.Write([]byte(format))

	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))

	if err != nil {
		log.Fatal(err)
	}

	io.Copy(part, file)
	writer.Close()

	url := fmt.Sprintf("https://%s/%s/lists?api_key=%s", hostname, apiVersion, url.QueryEscape(apiKey))

	client := http.Client{
		Timeout: time.Minute * 30,
	}

	req, err := http.NewRequest(http.MethodPost, url, requestBody)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Set("User-Agent", fmt.Sprintf("geocodio-cli"))

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	responseBody, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	return responseBody
}
