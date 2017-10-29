package rest

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/mtdx/keyc/common"

	"github.com/mtdx/keyc/account"
	"github.com/mtdx/keyc/db"
	"github.com/mtdx/keyc/keys"
	"github.com/mtdx/keyc/labels"
	"github.com/mtdx/keyc/openid/steamauth"
	"github.com/mtdx/keyc/steam"
	"github.com/mtdx/keyc/validator"
	"github.com/mtdx/keyc/vault"
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

func assertAuthRequired(t *testing.T, ts *httptest.Server, method string, route string) {
	t.Parallel()
	_, body := callEndpoint(t, ts, method, route, nil, "jwt-test")
	if strings.Compare(strings.TrimSpace(body), `Unauthorized`) != 0 {
		t.Fatalf("got: %s", body)
	}
}

func TestMain(m *testing.M) {
	dbconn := db.Open()
	defer dbconn.Close()

	r := StartRouter(dbconn)
	ts = httptest.NewServer(r)
	defer ts.Close()

	// setup test user
	cleanTestUserData(dbconn)
	jwt = setupTestUserData(dbconn)

	// run tests & clean & exit
	code := m.Run()
	os.Exit(code)
}

func AccountSummaryCheck(t *testing.T, bitcoinBalance float64, tradeLinik string) {
	_, body = callEndpoint(t, ts, "GET", "/api/v1/account", nil, jwt)
	var inforesp = &account.InfoResponse{}
	if err := json.Unmarshal([]byte(body), &inforesp); err != nil {
		t.Fatalf("Failed to Unmarshal, got: %s, error: %s", body, err.Error())
	}

	if err := validator.Validate(inforesp); err != nil {
		t.Fatalf("got: %s", err.Error())
	}

	if inforesp.BitcoinBalance != bitcoinBalance || inforesp.TradeLinkURL.String != tradeLinik {
		t.Fatalf("got: %s", body)
	}
}
func TestAccountSummaryAuthRequired(t *testing.T) {
	assertAuthRequired(t, ts, "GET", "/api/v1/account")
}

func TestAccountSummary(t *testing.T) {
	t.Parallel()

	AccountSummaryCheck(t, 1, "")
}

func TestTradeoffersAuthRequired(t *testing.T) {
	assertAuthRequired(t, ts, "GET", "/api/v1/tradeoffers")
}

func TestTradeoffers(t *testing.T) {
	t.Parallel()

	_, body = callEndpoint(t, ts, "GET", "/api/v1/tradeoffers", nil, jwt)
	tradeoffersresp := make([]steam.TradeoffersResponse, 2)
	if err := json.Unmarshal([]byte(body), &tradeoffersresp); err != nil {
		t.Fatalf("Failed to Unmarshal, got: %s, error: %s", body, err.Error())
	}

	for _, tradeoffer := range tradeoffersresp {
		if err := validator.Validate(tradeoffer); err != nil {
			t.Fatalf("got: %s", err.Error())
		}
	}

	if len(tradeoffersresp) != 2 || tradeoffersresp[0].Status != labels.ACTIVE ||
		tradeoffersresp[1].Type != labels.CSGO_CASE {
		t.Fatalf("got: %s", body)
	}
}

func TestPurchasesAuthRequired(t *testing.T) {
	assertAuthRequired(t, ts, "GET", "/api/v1/keys-transactions")
}

func TestPurchases(t *testing.T) {
	t.Parallel()

	_, body = callEndpoint(t, ts, "GET", "/api/v1/keys-transactions", nil, jwt)
	purchasesresp := make([]keys.TransactionsResponse, 2)
	if err := json.Unmarshal([]byte(body), &purchasesresp); err != nil {
		t.Fatalf("Failed to Unmarshal, got: %s, error: %s", body, err.Error())
	}

	for _, purchase := range purchasesresp {
		if err := validator.Validate(purchase); err != nil {
			t.Fatalf("got: %s", err.Error())
		}
	}

	if len(purchasesresp) != 2 || purchasesresp[0].Status != labels.UNPAID ||
		purchasesresp[1].Type != labels.CSGO_KEY {
		t.Fatalf("got: %s", body)
	}
}

func TestWithdrawalsAuthRequired(t *testing.T) {
	assertAuthRequired(t, ts, "GET", "/api/v1/withdrawals")
}

func TestWithdrawals(t *testing.T) {
	t.Parallel()

	bb1 := 0.00041839
	bb2 := 0.00086091
	withdrawalreq1 := &vault.WithdrawalsRequest{
		PaymentAddress: "1Gj4mwxWC8W9yhnrK5fVfYC2oV1jNdpiCS",
		CryptoTotal:    bb1,
	}
	withdrawalreq2 := &vault.WithdrawalsRequest{
		PaymentAddress: "1PUFW7bfU7if63UcLyUN8WsoZkpBuUVtUv",
		CryptoTotal:    bb2,
	}
	withdrawalreq3 := &vault.WithdrawalsRequest{
		PaymentAddress: "1PUFW7bfU7if63UcLyUN8WsoZkpBuUVtUv",
		CryptoTotal:    1,
	}
	jsonreq, _ = json.Marshal(withdrawalreq1)
	_, body = callEndpoint(t, ts, "POST", "/api/v1/withdrawals", bytes.NewReader(jsonreq), jwt)
	jsonreq, _ = json.Marshal(withdrawalreq2)
	_, body = callEndpoint(t, ts, "POST", "/api/v1/withdrawals", bytes.NewReader(jsonreq), jwt)
	jsonreq, _ = json.Marshal(withdrawalreq3)
	_, body = callEndpoint(t, ts, "POST", "/api/v1/withdrawals", bytes.NewReader(jsonreq), jwt)

	var errLowBalance common.ErrResponse
	if err := json.Unmarshal([]byte(body), &errLowBalance); err != nil {
		t.Fatalf("Failed to Unmarshal, got: %s, error: %s", body, err.Error())
	}
	if errLowBalance.ErrorText != "Not enough balance" {
		t.Fatalf("Failed to reject withdrawal, got %s", errLowBalance.StatusText)
	}

	_, body = callEndpoint(t, ts, "GET", "/api/v1/withdrawals", nil, jwt)
	withdrawalsresp := make([]vault.WithdrawalsResponse, 2)
	if err := json.Unmarshal([]byte(body), &withdrawalsresp); err != nil {
		t.Fatalf("Failed to Unmarshal, got: %s, error: %s", body, err.Error())
	}

	for _, withdrawal := range withdrawalsresp {
		if err := validator.Validate(withdrawal); err != nil {
			t.Fatalf("got: %s", err.Error())
		}
	}

	if len(withdrawalsresp) != 2 || withdrawalsresp[0].Status != labels.PENDING ||
		withdrawalsresp[1].Currency != labels.BTC || withdrawalsresp[0].CryptoTotal != bb2 ||
		withdrawalsresp[1].CryptoTotal != bb1 {
		t.Fatalf("got: %s", body)
	}

	// test balance change
	AccountSummaryCheck(t, 1-bb1-bb2, "")
}

func setupTestUserData(dbconn *sql.DB) string {
	jwt, err = steamauth.SaveUser(dbconn, testSteamID, "PersonaName", "https://avatar.com/img.jpg")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to generate token: %v\n", err)
		os.Exit(1)
	}

	_, err = dbconn.Exec(`INSERT INTO tradeoffers (id, user_steam_id, type, status, failure_details, amount, app_id)
	 VALUES (1, $1, $2, $3, $4, $5, $6)`, testSteamID, labels.CSGO_CASE, labels.ACTIVE, "failuredetails1", 1, labels.CSGO_APP_ID)
	_, err = dbconn.Exec(`INSERT INTO tradeoffers (id, user_steam_id, type, status, failure_details, amount, app_id)
	 VALUES (2, $1, $2, $3, $4, $5, $6)`, testSteamID, labels.CSGO_CASE, labels.ACTIVE, "failuredetails2", 2, labels.CSGO_APP_ID)

	_, err = dbconn.Exec(`INSERT INTO key_transactions (id, user_steam_id, tradeoffer_id, status, type, amount, unit_price,
		payment_address, usd_rate, currency, usd_total, crypto_total, app_id) VALUES (1, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		testSteamID, 1, labels.UNPAID, labels.CSGO_KEY, 1, 1.86, "13XrFK2m8tXvM5srR9tFPYsm2mpmRyAnXb", 4.3343, labels.BTC, 200, 0.00425076, labels.CSGO_APP_ID)
	_, err = dbconn.Exec(`INSERT INTO key_transactions (id, user_steam_id, tradeoffer_id, status, type, amount, unit_price,
			payment_address, usd_rate, currency, usd_total, crypto_total, app_id) VALUES (2, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		testSteamID, 2, labels.UNPAID, labels.CSGO_KEY, 2, 1.92, "13XrFK2m8tXvM5srR9tFPYsm2mpmRyAnXb", 5.2212, labels.BTC, 200, 0.00568021, labels.CSGO_APP_ID)

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
	_, err := dbconn.Exec(`DELETE FROM tradeoffers WHERE user_steam_id = $1`, testSteamID)
	_, err = dbconn.Exec(`DELETE FROM users WHERE steam_id = $1`, testSteamID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
