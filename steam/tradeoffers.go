package steam

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/mtdx/keyc/common"
)

// TradeoffersHandler rest route handler
func TradeoffersHandler(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	dbconn := r.Context().Value("DBCONN").(*sql.DB)
	tradeoffersresp, err := findAllTradeoffers(dbconn, claims["id"])
	common.ValidateRenderResults(w, r, tradeoffersresp, err)
}
