package common

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/render"
	"github.com/mtdx/keyc/validator"
	"github.com/mtdx/mx/common"
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

// AssertAuthRequired ...
func AssertAuthRequired(t *testing.T, ts *httptest.Server, method string, route string) {
	t.Parallel()
	_, body := TestRequest(t, ts, method, route, nil, "jwt-test")
	if strings.Compare(strings.TrimSpace(body), `Unauthorized`) != 0 {
		t.Fatalf("got: %s", body)
	}
}

// RenderResults ... renders `multiple` results
func RenderResults(w http.ResponseWriter, r *http.Request, resp []render.Renderer, err error) {
	for _, entry := range resp {
		if err := validator.Validate(entry); err != nil {
			render.Render(w, r, common.ErrInternalServer(err))
		}
	}

	render.Status(r, http.StatusOK)
	if err := render.RenderList(w, r, resp); err != nil {
		render.Render(w, r, common.ErrRender(err))
	}
}
