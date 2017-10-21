package account

import (
	"database/sql"
)

// InfoResponse ...
type InfoResponse struct {
	BitcoinBalance float64        `json:"bitcoin_balance" validate:"min=0"`
	CsgokeyBalance uint32         `json:"csgokey_balance" validate:"min=0"`
	TradeLinkURL   sql.NullString `json:"trade_link_url"`
}
