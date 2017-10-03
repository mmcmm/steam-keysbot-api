package account

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/mtdx/keyc/common"
	"github.com/mtdx/keyc/validator"
)

// TradeoffersResponse ...
type TradeoffersResponse struct {
	Type           string         `json:"type" validate:"nonzero"`
	Status         string         `json:"status" validate:"nonzero"`
	FailureDetails sql.NullString `json:"failure_details"`
	Amount         uint32         `json:"amount" validate:"min=1"`
	CreatedAt      string         `json:"created_at" validate:"nonzero"`
}

func (rd *TradeoffersResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// TradeoffersHandler rest route handler
func TradeoffersHandler(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	dbconn := r.Context().Value("DBCONN").(*sql.DB)
	rows, err := dbconn.Query(`SELECT type, status, failure_details, amount, created_at 
		FROM tradeoffers WHERE user_steam_id = $1`, claims["id"])
	if err != nil {
		render.Render(w, r, common.ErrInternalServer(err))
		return
	}
	defer rows.Close()
	tradeoffersresp := []render.Renderer{}
	for rows.Next() {
		resp := &TradeoffersResponse{}
		if err := rows.Scan(&resp.Type, &resp.Status, &resp.FailureDetails, &resp.Amount, &resp.CreatedAt); err != nil {
			render.Render(w, r, common.ErrInternalServer(err))
			return
		}
		tradeoffersresp = append(tradeoffersresp, resp)
	}
	if err := rows.Err(); err != nil {
		render.Render(w, r, common.ErrInternalServer(err))
		return
	}

	for _, tradeoffer := range tradeoffersresp {
		if err := validator.Validate(tradeoffer); err != nil {
			render.Render(w, r, common.ErrInternalServer(err))
		}
	}

	render.Status(r, http.StatusOK)
	if err := render.RenderList(w, r, tradeoffersresp); err != nil {
		render.Render(w, r, common.ErrRender(err))
	}
}
