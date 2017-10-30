package rest

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/mtdx/assert"
	"github.com/mtdx/keyc/common"
	"github.com/mtdx/keyc/labels"
	"github.com/mtdx/keyc/validator"
	"github.com/mtdx/keyc/vault"
)

func withdrawalsCheck(t *testing.T) {
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
	assert.Equal(t, "Not enough balance", errLowBalance.ErrorText, body)

	_, body = callEndpoint(t, ts, "GET", "/api/v1/withdrawals", nil, jwt)
	withdrawalsresp := make([]vault.WithdrawalsResponse, 2)
	if err := json.Unmarshal([]byte(body), &withdrawalsresp); err != nil {
		t.Fatalf("Failed to Unmarshal, got: %s, error: %s", body, err.Error())
	}

	for _, withdrawal := range withdrawalsresp {
		if err := validator.Validate(withdrawal); err != nil {
			t.Fatalf("Error: %s", err.Error())
		}
	}

	assert.Equal(t, 2, len(withdrawalsresp), body)
	assert.Equal(t, labels.PENDING, int(withdrawalsresp[0].Status), body)
	assert.Equal(t, labels.BTC, int(withdrawalsresp[1].Currency), body)
	assert.Equal(t, bb2, float64(withdrawalsresp[0].CryptoTotal), body)
	assert.Equal(t, bb1, float64(withdrawalsresp[1].CryptoTotal), body)

	// test balance change
	accountSummaryCheck(t, 1-bb1-bb2, "")
}
