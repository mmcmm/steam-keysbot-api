package account

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/mtdx/keyc/common"
	"github.com/mtdx/keyc/validator"
)

// PurchasesResponse ...
type PurchasesResponse struct {
	Status         string  `json:"status" validate:"nonzero"`
	Type           string  `json:"type" validate:"nonzero"`
	Amount         uint32  `json:"amount" validate:"min=1"`
	UnitPrice      float64 `json:"unit_price" validate:"nonzero"`
	PaymentAddress string  `json:"payment_address" validate:"nonzero"`
	USDRate        float64 `json:"usd_rate" validate:"nonzero"`
	Currency       string  `json:"currency" validate:"len=3"`
	USDTotal       float64 `json:"usd_total" validate:"nonzero"`
	CryptoTotal    float64 `json:"crypto_total" validate:"nonzero"`
	CreatedAt      string  `json:"created_at" validate:"nonzero"`
}

func (rd *PurchasesResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// PurchasesHandler rest route handler
func PurchasesHandler(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	dbconn := r.Context().Value("DBCONN").(*sql.DB)
	rows, err := dbconn.Query(`SELECT status, type, amount, unit_price, payment_address, usd_rate, 
		currency, usd_total, crypto_total, created_at FROM purchases WHERE user_steam_id = $1`,
		claims["id"])
	if err != nil {
		render.Render(w, r, common.ErrInternalServer(err))
		return
	}
	defer rows.Close()
	purchasesresp := []render.Renderer{}
	for rows.Next() {
		resp := &PurchasesResponse{}
		if err := rows.Scan(&resp.Status, &resp.Type, &resp.Amount, &resp.UnitPrice, &resp.PaymentAddress,
			&resp.USDRate, &resp.Currency, &resp.USDTotal, &resp.CryptoTotal, &resp.CreatedAt); err != nil {
			render.Render(w, r, common.ErrInternalServer(err))
			return
		}
		purchasesresp = append(purchasesresp, resp)
	}
	if err := rows.Err(); err != nil {
		render.Render(w, r, common.ErrInternalServer(err))
		return
	}

	for _, purchase := range purchasesresp {
		if err := validator.Validate(purchase); err != nil {
			render.Render(w, r, common.ErrInternalServer(err))
		}
	}

	render.Status(r, http.StatusOK)
	if err := render.RenderList(w, r, purchasesresp); err != nil {
		render.Render(w, r, common.ErrRender(err))
	}
}
