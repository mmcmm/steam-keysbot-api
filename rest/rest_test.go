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
	"time"

	"github.com/mtdx/keyc/db"
	"github.com/mtdx/keyc/internal"
	"github.com/mtdx/keyc/openid/steamauth"
)

const testSteamID = "11111111111111111"
const testTradeOfferID1 = "22222222222222222"
const testTradeOfferID2 = "33333333333333333"
const testBoSteamID = "44444444444444444"

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

	internal.InitLiveBtc()

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

func TestWithdrawalsAuth(t *testing.T) {
	t.Parallel()
	assertAuth(t, ts, "GET", "/api/v1/withdrawals")
}

func TestTransactionsAuth(t *testing.T) {
	t.Parallel()
	assertAuth(t, ts, "GET", "/api/v1/keys-transactions")
}

// full tests
func TestAccountSummary(t *testing.T) {
	t.Parallel()
	accountSummaryCheck(t, 1, "")
}

func TestWithdrawals(t *testing.T) {
	t.Parallel()
	time.Sleep(2 * time.Second)
	withdrawalsCheck(t)
}

func TestTransactions(t *testing.T) {
	t.Parallel()
	time.Sleep(2 * time.Second)
	transactionsCheck(t)
}

func setupTestUserData(dbconn *sql.DB) string {
	jwt, err = steamauth.SaveUser(dbconn, testSteamID, "PersonaName", "https://avatar.com/img.jpg")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to generate token: %v\n", err)
		os.Exit(1)
	}

	_, err = dbconn.Exec(`INSERT INTO steam_bots (steam_id, ip_address, display_name) VALUES ($1, $2, $3)`, testBoSteamID, "127.0.0.1", "testuser")
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
