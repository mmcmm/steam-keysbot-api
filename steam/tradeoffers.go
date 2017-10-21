package steam

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/mtdx/keyc/common"
)

func (rd *TradeoffersResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// TradeoffersHandler rest route handler
func TradeoffersHandler(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	dbconn := r.Context().Value("DBCONN").(*sql.DB)
	tradeoffersresp, err := findAllTradeoffers(dbconn, claims["id"])
	if err != nil {
		render.Render(w, r, common.ErrInternalServer(err))
		return
	}

	common.RenderResults(w, r, tradeoffersresp, err)
}
