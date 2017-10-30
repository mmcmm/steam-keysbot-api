package rest

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/mtdx/keyc/db"
	"github.com/mtdx/keyc/openid/steamauth"
)

const testSteamID = "11111111111111111"

var ts *httptest.Server
var body, jwt string
var jsonreq []byte
var err error

func callEndpoint(t *testing.T, ts *httptest.Server, method, path string, body io.Reader, jwt string) (*http.Response, string) {
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

func TestMain(m *testing.M) {
	dbconn := db.Open()
	defer dbconn.Close()

	r := Router(dbconn)
	ts = httptest.NewServer(r)
	defer ts.Close()

	// setup test user
	cleanTestUserData(dbconn)
	jwt = setupTestUserData(dbconn)

	// run tests & clean & exit
	code := m.Run()
	os.Exit(code)
}

// Auth tests

func TestAccountSummaryAuth(t *testing.T) {
	t.Parallel()
	assertAuth(t, ts, "GET", "/api/v1/account")
}

func TestTradeoffersAuth(t *testing.T) {
	t.Parallel()
	assertAuth(t, ts, "GET", "/api/v1/tradeoffers")
}

func TestWithdrawalsAuth(t *testing.T) {
	t.Parallel()
	assertAuth(t, ts, "GET", "/api/v1/withdrawals")
}

// full tests

func TestAccountSummary(t *testing.T) {
	t.Parallel()
	accountSummaryCheck(t, 1, "")
}

func TestTradeoffers(t *testing.T) {
	t.Parallel()
	tradeoffersCheck(t)
}

func TestWithdrawals(t *testing.T) {
	t.Parallel()
	withdrawalsCheck(t)
}

// func TestKeysTransactionsAuthRequired(t *testing.T) {
// 	assertAuthRequired(t, ts, "GET", "/api/v1/keys-transactions")
// }

// func TestKeysTransactions(t *testing.T) {
// 	t.Parallel()

// 	_, body = callEndpoint(t, ts, "GET", "/api/v1/keys-transactions", nil, jwt)
// 	purchasesresp := make([]keys.TransactionsResponse, 2)
// 	if err := json.Unmarshal([]byte(body), &purchasesresp); err != nil {
// 		t.Fatalf("Failed to Unmarshal, got: %s, error: %s", body, err.Error())
// 	}

// 	for _, purchase := range purchasesresp {
// 		if err := validator.Validate(purchase); err != nil {
// 			t.Fatalf("got: %s", err.Error())
// 		}
// 	}

// 	if len(purchasesresp) != 2 || purchasesresp[0].Status != labels.UNPAID ||
// 		purchasesresp[1].Type != labels.CSGO_KEY {
// 		t.Fatalf("got: %s", body)
// 	}
// }

func setupTestUserData(dbconn *sql.DB) string {
	jwt, err = steamauth.SaveUser(dbconn, testSteamID, "PersonaName", "https://avatar.com/img.jpg")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to generate token: %v\n", err)
		os.Exit(1)
	}

	// TODO: implement
	// _, err = dbconn.Exec(`INSERT INTO key_transactions (id, user_steam_id, tradeoffer_id, status, type, amount, unit_price,
	// 	payment_address, usd_rate, currency, usd_total, crypto_total, app_id) VALUES (1, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
	// 	testSteamID, 1, labels.UNPAID, labels.CSGO_KEY, 1, 1.86, "13XrFK2m8tXvM5srR9tFPYsm2mpmRyAnXb", 4.3343, labels.BTC, 200, 0.00425076, labels.CSGO_APP_ID)
	// _, err = dbconn.Exec(`INSERT INTO key_transactions (id, user_steam_id, tradeoffer_id, status, type, amount, unit_price,
	// 		payment_address, usd_rate, currency, usd_total, crypto_total, app_id) VALUES (2, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
	// 	testSteamID, 2, labels.UNPAID, labels.CSGO_KEY, 2, 1.92, "13XrFK2m8tXvM5srR9tFPYsm2mpmRyAnXb", 5.2212, labels.BTC, 200, 0.00568021, labels.CSGO_APP_ID)

	_, err = dbconn.Exec(`INSERT INTO steam_bots (steam_id, ip_address, display_name) VALUES ($1, $2, $3)`, testSteamID, "127.0.0.1", "testuser")
	_, err = dbconn.Exec(`UPDATE users SET bitcoin_balance = 1 WHERE steam_id = $1`, testSteamID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to add test data: %v\n", err)
		os.Exit(1)
	}

	return jwt
}

func cleanTestUserData(dbconn *sql.DB) {
	_, err = dbconn.Exec(`DELETE FROM key_transactions WHERE user_steam_id = $1`, testSteamID)
	_, err = dbconn.Exec(`DELETE FROM withdrawals WHERE user_steam_id = $1`, testSteamID)
	_, err = dbconn.Exec(`DELETE FROM sales WHERE user_steam_id = $1`, testSteamID)
	_, err = dbconn.Exec(`DELETE FROM tradeoffers WHERE user_steam_id = $1`, testSteamID)
	_, err = dbconn.Exec(`DELETE FROM steam_bots WHERE steam_id = $1`, testSteamID)
	_, err = dbconn.Exec(`DELETE FROM users WHERE steam_id = $1`, testSteamID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
