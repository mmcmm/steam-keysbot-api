package vault

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/mtdx/keyc/common"
)

func (rd *WithdrawalsResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// WithdrawalsHandler rest route handler
func WithdrawalsHandler(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	dbconn := r.Context().Value("DBCONN").(*sql.DB)
	withdrawalsresp, err := findAllWithdrawals(dbconn, claims["id"])
	if err != nil {
		render.Render(w, r, common.ErrInternalServer(err))
		return
	}
	common.RenderResults(w, r, withdrawalsresp, err)
}

// RequestWithdrawalHandler put /withdrawal
func RequestWithdrawalHandler(w http.ResponseWriter, r *http.Request) {

}
