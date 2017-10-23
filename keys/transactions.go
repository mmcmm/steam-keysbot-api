package keys

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/mtdx/keyc/common"
)

func (tr *TransactionsResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// TransactionsHandler rest route handler
func TransactionsHandler(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	dbconn := r.Context().Value("DBCONN").(*sql.DB)
	transactionsresp, err := findAllTransactions(dbconn, claims["id"])
	common.ValidateRenderResults(w, r, transactionsresp, err)
}
