package vault

// WithdrawalsResponse ...
type WithdrawalsResponse struct {
	Status         uint8   `json:"status" validate:"nonzero"`
	PaymentAddress string  `json:"payment_address" validate:"nonzero"`
	UsdRate        float32 `json:"usd_rate" validate:"nonzero"`
	Currency       uint8   `json:"currency" validate:"nonzero"`
	USDTotal       float64 `json:"usd_total" validate:"nonzero"`
	CryptoTotal    float64 `json:"crypto_total" validate:"nonzero"`
	Txhash         string  `json:"txhash" validate:"nonzero"`
	CreatedAt      string  `json:"created_at" validate:"nonzero"`
}
