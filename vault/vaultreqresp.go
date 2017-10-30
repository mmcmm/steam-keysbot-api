package vault

import (
	"database/sql"
	"net/http"
)

// Render pre-processing after a decoding.
func (wi *WithdrawalsResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Bind post-processing after a decoding.
func (wi *WithdrawalsRequest) Bind(r *http.Request) error {
	return nil
}

// WithdrawalsResponse ...
type WithdrawalsResponse struct {
	Status         uint8          `json:"status" validate:"nonzero"`
	PaymentAddress string         `json:"payment_address" validate:"nonzero"`
	UsdRate        float32        `json:"usd_rate" validate:"nonzero"`
	Currency       uint8          `json:"currency" validate:"nonzero"`
	USDTotal       float64        `json:"usd_total" validate:"nonzero"`
	CryptoTotal    float64        `json:"crypto_total" validate:"nonzero"`
	Txhash         sql.NullString `json:"txhash"`
	CreatedAt      string         `json:"created_at" validate:"nonzero"`
}

// WithdrawalsRequest ...
type WithdrawalsRequest struct {
	PaymentAddress string  `json:"payment_address" validate:"nonzero"`
	CryptoTotal    float64 `json:"crypto_total" min:"0.00000001"`
}
