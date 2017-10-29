package account

import (
	"database/sql"
	"net/http"
)

// Render pre-processing after a decoding.
func (acc *InfoResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Bind post-processing after a decoding.
func (acc *InfoResponse) Bind(r *http.Request) error {
	return nil
}

// InfoResponse ...
type InfoResponse struct {
	BitcoinBalance float64        `json:"bitcoin_balance" validate:"min=0"`
	CsgokeyBalance uint32         `json:"csgokey_balance" validate:"min=0"`
	TradeLinkURL   sql.NullString `json:"trade_link_url"`
}
