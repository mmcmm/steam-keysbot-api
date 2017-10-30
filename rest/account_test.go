package rest

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mtdx/assert"
	"github.com/mtdx/keyc/account"
	"github.com/mtdx/keyc/validator"
)

func assertAuth(t *testing.T, ts *httptest.Server, method string, route string) {
	_, body := callEndpoint(t, ts, method, route, nil, "jwt-test")
	if strings.Compare(strings.TrimSpace(body), `Unauthorized`) != 0 {
		t.Fatalf("got: %s", body)
	}
}

func accountSummaryCheck(t *testing.T, bitcoinBalance float64, tradeLinik string) {
	_, body = callEndpoint(t, ts, "GET", "/api/v1/account", nil, jwt)
	var inforesp = &account.InfoResponse{}
	if err := json.Unmarshal([]byte(body), &inforesp); err != nil {
		t.Fatalf("Failed to Unmarshal, got: %s, error: %s", body, err.Error())
	}

	if err := validator.Validate(inforesp); err != nil {
		t.Fatalf("Error: %s", err.Error())
	}

	assert.Equal(t, bitcoinBalance, inforesp.BitcoinBalance, body)
	assert.Equal(t, tradeLinik, inforesp.TradeLinkURL.String, body)
}
