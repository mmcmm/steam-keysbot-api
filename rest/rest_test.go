package rest

import (
	"database/sql"
	"fmt"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/mtdx/keyc/common"
	"github.com/mtdx/keyc/db"
	"github.com/mtdx/keyc/openid/steamauth"
)

const testSteamID = "11111111111111111"

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
	jwt, err = steamauth.SaveUser(dbconn, testSteamID, "PersonaName", "https://avatar.com/img.jpg")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to generate token: %v\n", err)
		os.Exit(1)
	}

	// run tests & clean & exit
	code := m.Run()
	cleanDb(dbconn)
	os.Exit(code)
}

func TestAccountSummaryAuth(t *testing.T) {
	t.Parallel()

	_, body = common.TestRequest(t, ts, "GET", "/api/v1/account", nil, "jwt")
	expected = `Unauthorized`
	if strings.Compare(strings.TrimSpace(body), expected) != 0 {
		t.Fatalf("expected:%s got:%s", expected, body)
	}
}

func TestAccountSummary(t *testing.T) {
	t.Parallel()

	_, body = common.TestRequest(t, ts, "GET", "/api/v1/account", nil, jwt)
	expected = `{"bitcoin_balance":0,"csgokey_balance":0,"trade_link_url":{"String":"","Valid":false}}`
	if strings.Compare(strings.TrimSpace(body), expected) != 0 {
		t.Fatalf("expected:%s got:%s", expected, body)
	}
}

func cleanDb(dbconn *sql.DB) {
	dbconn.Exec(`DELETE FROM users WHERE steam_id = $1`, testSteamID)
	dbconn.Exec(`DELETE FROM tradeoffers WHERE user_steam_id = $1`, testSteamID)
	dbconn.Exec(`DELETE FROM purchases WHERE user_steam_id = $1`, testSteamID)
	dbconn.Exec(`DELETE FROM withdrawals WHERE user_steam_id = $1`, testSteamID)
	dbconn.Exec(`DELETE FROM sales WHERE user_steam_id = $1`, testSteamID)
}
