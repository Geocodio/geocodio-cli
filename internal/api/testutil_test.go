package api_test

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/geocodio/geocodio-cli/internal/api"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"
)

func newTestClient(t *testing.T, cassetteName string) *api.Client {
	t.Helper()

	cassettePath := filepath.Join("testdata", cassetteName)
	mode := recorder.ModeReplayOnly

	if os.Getenv("VCR_MODE") == "record" {
		mode = recorder.ModeRecordOnly
	} else {
		if _, err := os.Stat(cassettePath + ".yaml"); os.IsNotExist(err) {
			t.Skipf("cassette %s not found, run with VCR_MODE=record to create", cassetteName)
		}
	}

	r, err := recorder.New(cassettePath,
		recorder.WithMode(mode),
		recorder.WithHook(redactHook, recorder.BeforeSaveHook),
		recorder.WithMatcher(matcherIgnoringAPIKey),
	)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}

	t.Cleanup(func() {
		if err := r.Stop(); err != nil {
			t.Errorf("failed to stop recorder: %v", err)
		}
	})

	apiKey := os.Getenv("GEOCODIO_API_KEY")
	if apiKey == "" {
		apiKey = "test-api-key"
	}

	return api.NewClient(
		"https://api.geocod.io/v1.9",
		apiKey,
		api.WithHTTPClient(&http.Client{Transport: r}),
	)
}

func redactHook(i *cassette.Interaction) error {
	u, err := url.Parse(i.Request.URL)
	if err != nil {
		return err
	}
	q := u.Query()
	if q.Get("api_key") != "" {
		q.Set("api_key", "REDACTED")
		u.RawQuery = q.Encode()
		i.Request.URL = u.String()
	}
	delete(i.Request.Form, "api_key")
	return nil
}

func matcherIgnoringAPIKey(r *http.Request, i cassette.Request) bool {
	if r.Method != i.Method {
		return false
	}

	reqURL, _ := url.Parse(r.URL.String())
	cassetteURL, _ := url.Parse(i.URL)

	if reqURL.Host != cassetteURL.Host || reqURL.Path != cassetteURL.Path {
		return false
	}

	reqQuery := reqURL.Query()
	cassetteQuery := cassetteURL.Query()

	reqQuery.Del("api_key")
	cassetteQuery.Del("api_key")

	return reqQuery.Encode() == cassetteQuery.Encode()
}
