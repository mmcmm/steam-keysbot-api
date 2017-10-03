package rest

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/mtdx/keyc/account"
	"github.com/mtdx/keyc/common"
	"github.com/mtdx/keyc/db"
	"github.com/mtdx/keyc/openid/steamauth"
	"github.com/mtdx/keyc/validator"
)

const testSteamID = "11111111111111111"

var ts *httptest.Server
var body, jwt string
var err error

func TestMain(m *testing.M) {
	dbconn := db.Open()
	defer dbconn.Close()

	r := StartRouter(dbconn)
	ts = httptest.NewServer(r)
	defer ts.Close()

	// setup test user
	jwt = setupTestUserData(dbconn)

	// run tests & clean & exit
	code := m.Run()
	cleanDb(dbconn)
	os.Exit(code)
}

func TestAccountSummaryAuth(t *testing.T) {
	common.AssertAuth(t, ts, "GET", "/api/v1/account")
}

func TestAccountSummary(t *testing.T) {
	t.Parallel()

	_, body = common.TestRequest(t, ts, "GET", "/api/v1/account", nil, jwt)
	var inforesp = &account.InfoResponse{}
	if err := json.Unmarshal([]byte(body), &inforesp); err != nil {
		t.Fatalf("Failed to Unmarshal, got: %s, error: %s", body, err.Error())
	}

	if err := validator.Validate(inforesp); err != nil {
		t.Fatalf("got: %s", err.Error())
	}

	if inforesp.BitcoinBalance != 0 || inforesp.CsgokeyBalance != 0 || inforesp.TradeLinkURL.String != "" {
		t.Fatalf("got: %s", body)
	}
}

func TestTradeoffersAuth(t *testing.T) {
	common.AssertAuth(t, ts, "GET", "/api/v1/account/tradeoffers")
}

func TestTradeoffers(t *testing.T) {
	t.Parallel()

	_, body = common.TestRequest(t, ts, "GET", "/api/v1/account/tradeoffers", nil, jwt)
	tradeoffersresp := make([]account.TradeoffersResponse, 2)
	if err := json.Unmarshal([]byte(body), &tradeoffersresp); err != nil {
		t.Fatalf("Failed to Unmarshal, got: %s, error: %s", body, err.Error())
	}

	for _, tradeoffer := range tradeoffersresp {
		if err := validator.Validate(tradeoffer); err != nil {
			t.Fatalf("got: %s", err.Error())
		}
	}

	if len(tradeoffersresp) != 2 || tradeoffersresp[0].Type != string(account.CSGO_KEY) ||
		tradeoffersresp[1].Type != string(account.CSGO_CASE) {
		t.Fatalf("got: %s", body)
	}
}

func cleanDb(dbconn *sql.DB) {
	dbconn.Exec(`DELETE FROM tradeoffers WHERE user_steam_id = $1`, testSteamID)
	dbconn.Exec(`DELETE FROM purchases WHERE user_steam_id = $1`, testSteamID)
	dbconn.Exec(`DELETE FROM withdrawals WHERE user_steam_id = $1`, testSteamID)
	dbconn.Exec(`DELETE FROM sales WHERE user_steam_id = $1`, testSteamID)
	dbconn.Exec(`DELETE FROM users WHERE steam_id = $1`, testSteamID)
}

func setupTestUserData(dbconn *sql.DB) string {
	jwt, err = steamauth.SaveUser(dbconn, testSteamID, "PersonaName", "https://avatar.com/img.jpg")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to generate token: %v\n", err)
		os.Exit(1)
	}

	dbconn.Exec(`INSERT INTO tradeoffers (id, user_steam_id, type, status, failure_details, amount) 
	 VALUES (1, $1, $2, $3, $4, $5)`, testSteamID, account.CSGO_KEY, account.ACTIVE, "failuredetails1", 1)
	dbconn.Exec(`INSERT INTO tradeoffers (id, user_steam_id, type, status, failure_details, amount) 
	 VALUES (2, $1, $2, $3, $4, $5)`, testSteamID, account.CSGO_CASE, account.ACCEPTED, "failuredetails2", 2)

	dbconn.Exec(`INSERT INTO purchases (id, user_steam_id, tradeoffer_id, status, type, amount, unit_price, 
		payment_address, usd_rate, currency, usd_total, crypto_total) 
	 VALUES (1, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`, testSteamID, 1, account.UNPAID, account.CSGO_KEY,
		1, 1.86, "13XrFK2m8tXvM5srR9tFPYsm2mpmRyAnXb", 4.3343, account.BTC, 200, 0.00425076)
	dbconn.Exec(`INSERT INTO purchases (id, user_steam_id, tradeoffer_id, status, type, amount, unit_price, 
			payment_address, usd_rate, currency, usd_total, crypto_total) 
		 VALUES (2, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`, testSteamID, 2, account.UNPAID, account.CSGO_CASE,
		2, 1.92, "13XrFK2m8tXvM5srR9tFPYsm2mpmRyAnXb", 5.2212, account.BTC, 200, 0.00568021)

	return jwt
}
