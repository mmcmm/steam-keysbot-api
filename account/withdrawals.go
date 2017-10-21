package account

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/mtdx/keyc/common"
)

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

func (rd *WithdrawalsResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// WithdrawalsHandler rest route handler
func WithdrawalsHandler(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	dbconn := r.Context().Value("DBCONN").(*sql.DB)
	rows, err := dbconn.Query(`SELECT status, payment_address, usd_rate, currency, usd_total, crypto_total,
		 txhash, created_at FROM withdrawals WHERE user_steam_id = $1`, claims["id"])
	if err != nil {
		render.Render(w, r, common.ErrInternalServer(err))
		return
	}
	defer rows.Close()
	withdrawalsresp := []render.Renderer{}
	for rows.Next() {
		resp := &WithdrawalsResponse{}
		if err := rows.Scan(
			&resp.Status,
			&resp.PaymentAddress,
			&resp.UsdRate,
			&resp.Currency,
			&resp.USDTotal,
			&resp.CryptoTotal,
			&resp.Txhash,
			&resp.CreatedAt,
		); err != nil {
			render.Render(w, r, common.ErrInternalServer(err))
			return
		}
		withdrawalsresp = append(withdrawalsresp, resp)
	}

	common.RenderResults(w, r, withdrawalsresp, rows, err)
}

// RequestWithdrawalHandler put /withdrawal
func RequestWithdrawalHandler(w http.ResponseWriter, r *http.Request) {

}
