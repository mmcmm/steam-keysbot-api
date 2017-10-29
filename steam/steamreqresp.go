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
func (to *TradeoffersResponse) Bind(r *http.Request) error {
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
