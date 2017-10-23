package vault

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/mtdx/keyc/common"
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

}
