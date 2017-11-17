package keys

import "net/http"

// Render pre-processing after a decoding.
func (tr *TransactionsResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Bind post-processing after a decoding.
func (tr *TransactionsRequest) Bind(r *http.Request) error {
	return nil
}

// Bind post-processing after a decoding.
func (tr *TransactionsUpdateRequest) Bind(r *http.Request) error {
	return nil
}

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

// TransactionsRequest ...
type TransactionsRequest struct {
	UserSteamID     string `json:"user_steam_id" validate:"nonzero"`
	TradeofferID    string `json:"tradeoffer_id" validate:"nonzero"`
	BotSteamID      string `json:"bot_steam_id" validate:"nonzero"`
	Type            uint8  `json:"type" validate:"nonzero"`
	TransactionType uint8  `json:"transaction_type" validate:"nonzero"`
	Amount          uint32 `json:"amount" validate:"min=1"`
	PaymentAddress  string `json:"payment_address" validate:"nonzero"`
	Currency        uint8  `json:"currency" validate:"nonzero"`
	AppID           uint32 `json:"app_id" validate:"nonzero"`
}

// TransactionsUpdateRequest ...
type TransactionsUpdateRequest struct {
	TradeofferID string `json:"tradeoffer_id" validate:"nonzero"`
	Status       uint8  `json:"statis" validate:"nonzero"`
}
