package rest

import (
	"fmt"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/mtdx/keyc/common"
	"github.com/mtdx/keyc/db"
	"github.com/mtdx/keyc/openid/steamauth"
)

var ts *httptest.Server
var body, expected, jwt string
var err error

func TestMain(m *testing.M) {
	dbconn := db.Open()
	defer dbconn.Close()
	r := StartRouter(dbconn)
	ts = httptest.NewServer(r)
	defer ts.Close()

	// setup test user
	jwt, err = steamauth.SaveUser(dbconn, "0", "PersonaName", "https://avatar.com/img.jpg")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to generate token: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()
	os.Exit(code)
}

func TestAccount(t *testing.T) {
	t.Parallel()

	_, body = common.TestRequest(t, ts, "GET", "/api/v1/account", nil, jwt)
	expected = `{"bitcoin_balance":0,"csgokey_balance":0,"trade_link_url":{"String":"","Valid":false}}`
	if strings.Compare(strings.TrimSpace(body), expected) != 0 {
		t.Fatalf("expected:%s got:%s", expected, body)
	}
}
