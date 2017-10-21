package keys

// TransactionsResponse ...
type TransactionsResponse struct {
	Status         uint8   `json:"status" validate:"nonzero"`
	Type           uint8   `json:"type" validate:"nonzero"`
	Amount         uint32  `json:"amount" validate:"min=1"`
	UnitPrice      float64 `json:"unit_price" validate:"nonzero"`
	PaymentAddress string  `json:"payment_address" validate:"nonzero"`
	USDRate        float64 `json:"usd_rate" validate:"nonzero"`
	Currency       uint8   `json:"currency" validate:"nonzero"`
	USDTotal       float64 `json:"usd_total" validate:"nonzero"`
	CryptoTotal    float64 `json:"crypto_total" validate:"nonzero"`
	CreatedAt      string  `json:"created_at" validate:"nonzero"`
}
