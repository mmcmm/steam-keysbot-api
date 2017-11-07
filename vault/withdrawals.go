package vault

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/mtdx/keyc/common"
	"github.com/mtdx/keyc/validator"
)

// WithdrawalsHandler rest route handler
func WithdrawalsHandler(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	dbconn := r.Context().Value("DBCONN").(*sql.DB)
	withdrawalsresp, err := findAllWithdrawals(dbconn, claims["id"])
	common.ValidateRenderResults(w, r, withdrawalsresp, err)
}

// WithdrawalRequestHandler POST /withdrawals
func WithdrawalRequestHandler(w http.ResponseWriter, r *http.Request) {
	withdrawal := &WithdrawalsRequest{}
	if err := render.Bind(r, withdrawal); err != nil {
		render.Render(w, r, common.ErrInvalidRequest(err))
		return
	}
	if err := validator.Validate(withdrawal); err != nil {
		render.Render(w, r, common.ErrInvalidRequest(err))
		return
	}

	_, claims, _ := jwtauth.FromContext(r.Context())
	dbconn := r.Context().Value("DBCONN").(*sql.DB)
	if err := saveWithdrawal(dbconn, withdrawal, claims["id"]); err != nil {
		render.Render(w, r, common.ErrInvalidRequest(err))
		return
	}
	render.Render(w, r, common.SuccessCreatedResponse("Withdrawal has been created"))
}
