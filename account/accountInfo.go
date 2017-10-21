package account

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/mtdx/mx/common"
)

func (rd *InfoResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// InfoHandler rest route handler
func InfoHandler(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	dbconn := r.Context().Value("DBCONN").(*sql.DB)
	inforesp, err := getAccountInfo(dbconn, claims["id"])
	if err != nil {
		render.Render(w, r, common.ErrInternalServer(err))
		return
	}
	render.Status(r, http.StatusOK)
	render.Render(w, r, inforesp)
}
