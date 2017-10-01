package common

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestRequest ...
func TestRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader, jwt string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	if jwt != "" {
		req.Header.Add("Authorization", "BEARER "+jwt)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}
	defer resp.Body.Close()

	return resp, string(respBody)
}

// AssertAuth ...
func AssertAuth(t *testing.T, ts *httptest.Server, method string, route string) {
	t.Parallel()
	_, body := TestRequest(t, ts, method, route, nil, "jwt-test")
	if strings.Compare(strings.TrimSpace(body), `Unauthorized`) != 0 {
		t.Fatalf("got: %s", body)
	}
}
