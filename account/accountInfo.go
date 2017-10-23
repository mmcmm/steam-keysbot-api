package account

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/mtdx/keyc/common"
)

// InfoHandler rest route handler
func InfoHandler(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	dbconn := r.Context().Value("DBCONN").(*sql.DB)
	inforesp, err := getAccountInfo(dbconn, claims["id"])
	common.ValidateRenderResult(w, r, inforesp, err)
}
