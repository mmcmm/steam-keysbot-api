package steam

import (
	"database/sql"
	"net/http"
)

// Render pre-processing after a decoding.
func (to *TradeoffersResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Bind post-processing after a decoding.
func (tr *TradeoffersRequest) Bind(r *http.Request) error {
	return nil
}

// Bind post-processing after a decoding.
func (tr *TradeoffersUpdateRequest) Bind(r *http.Request) error {
	return nil
}

// TradeoffersResponse ...
type TradeoffersResponse struct {
	Type           uint8          `json:"type" validate:"nonzero"`
	Status         uint8          `json:"status" validate:"nonzero"`
	FailureDetails sql.NullString `json:"failure_details"`
	Amount         uint32         `json:"amount" validate:"min=1"`
	CreatedAt      string         `json:"created_at" validate:"nonzero"`
}

// TradeoffersRequest ...
type TradeoffersRequest struct {
	TradeofferID    string `json:"tradeoffer_id" validate:"nonzero"`
	SteamID         string `json:"steam_id" validate:"nonzero"`
	FailureDetails  string `json:"failure_details"`
	Type            uint8  `json:"type" validate:"nonzero"`
	MerchantSteamID string `json:"merchant_steam_id" validate:"nonzero"`
	Amount          uint32 `json:"amount" validate:"min=1"`
	AppID           uint   `json:"app_id" validate:"nonzero"`
}

// TradeoffersUpdateRequest ...
type TradeoffersUpdateRequest struct {
	Status         uint8  `json:"status" validate:"nonzero"`
	FailureDetails string `json:"failure_details"`
}
