package api

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

var apiVersion string = "v1.7"

func Request(method string, path string, hostname string, apiKey string) []byte {
	url := fmt.Sprintf("https://%s/%s/%s?api_key=%s", hostname, apiVersion, path, url.QueryEscape(apiKey))

	client := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
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
