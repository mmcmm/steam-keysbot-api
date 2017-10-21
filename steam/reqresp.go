package steam

import "database/sql"

// TradeoffersResponse ...
type TradeoffersResponse struct {
	Type           uint8          `json:"type" validate:"nonzero"`
	Status         uint8          `json:"status" validate:"nonzero"`
	FailureDetails sql.NullString `json:"failure_details"`
	Amount         uint32         `json:"amount" validate:"min=1"`
	CreatedAt      string         `json:"created_at" validate:"nonzero"`
}
