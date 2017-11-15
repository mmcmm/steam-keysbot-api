package rest

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/mtdx/assert"
	"github.com/mtdx/keyc/common"
	"github.com/mtdx/keyc/config"
	"github.com/mtdx/keyc/keys"
	"github.com/mtdx/keyc/labels"
	"github.com/mtdx/keyc/validator"
)

func transactionsCheck(t *testing.T) {
	transaction1 := &keys.TransactionsRequest{
		UserSteamID:     testSteamID,
		TradeofferID:    testTradeOfferID1,
		BotSteamID:      testBoSteamID,
		Type:            labels.SELLTOUS,
		TransactionType: labels.BUY_CSGOKEY_PRICE,
		Amount:          4,
		PaymentAddress:  "1ombvzJzcYRvNuTyTXe6pzVpcJmWVG4Ma",
		Currency:        labels.BTC,
		AppID:           labels.CSGO_APP_ID,
	}
	transaction2 := &keys.TransactionsRequest{
		UserSteamID:     testSteamID,
		TradeofferID:    testTradeOfferID2,
		BotSteamID:      testBoSteamID,
		Type:            labels.SELLTOUS,
		TransactionType: labels.SELL_CSGOKEY_PRICE,
		Amount:          8,
		PaymentAddress:  "1PUFW7bfU7if63UcLyUN8WsoZkpBuUVtUv",
		Currency:        labels.BTC,
		AppID:           labels.CSGO_APP_ID,
	}

	jsonreq, _ = json.Marshal(transaction1)
	_, body = callEndpoint(t, ts, "POST", "/api/v1/keys-transactions?key="+config.SteamBotsAPIKey(), bytes.NewReader(jsonreq), jwt)
	jsonreq, _ = json.Marshal(transaction2)
	_, body = callEndpoint(t, ts, "POST", "/api/v1/keys-transactions?key="+config.SteamBotsAPIKey(), bytes.NewReader(jsonreq), jwt)

	var successResp common.SuccessResponse
	if err := json.Unmarshal([]byte(body), &successResp); err != nil {
		t.Fatalf("Failed to Unmarshal, got: %s, error: %s", body, err.Error())
	}
	assert.Equal(t, "Transaction has been created", successResp.StatusText, successResp.StatusText+", "+successResp.SuccessText)

	jsonreq, _ = json.Marshal(transaction2)
	_, body = callEndpoint(t, ts, "POST", "/api/v1/keys-transactions", bytes.NewReader(jsonreq), jwt)
	var errResp common.ErrResponse
	if err := json.Unmarshal([]byte(body), &errResp); err != nil {
		t.Fatalf("Failed to Unmarshal, got: %s, error: %s", body, err.Error())
	}
	assert.Equal(t, errResp.StatusText, "Invalid request", errResp.StatusText)

	_, body = callEndpoint(t, ts, "GET", "/api/v1/keys-transactions", nil, jwt)
	transactionresp := make([]keys.TransactionsResponse, 2)
	if err := json.Unmarshal([]byte(body), &transactionresp); err != nil {
		t.Fatalf("Failed to Unmarshal, got: %s, error: %s", body, err.Error())
	}

	for _, transaction := range transactionresp {
		if err := validator.Validate(transaction); err != nil {
			t.Fatalf("Error: %s", err.Error())
		}
	}

	if assert.Equal(t, 2, len(transactionresp), body) {
		assert.Equal(t, labels.PENDING, int(transactionresp[0].Status), body)
		assert.Equal(t, labels.SELLTOUS, int(transactionresp[1].Type), body)
		assert.Equal(t, 8, int(transactionresp[0].Amount), body)
	}
}
