package keys

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/mtdx/keyc/common"
	"github.com/mtdx/keyc/config"
	"github.com/mtdx/keyc/steam"
)

// TransactionsHandler rest route handler GET /keys-transactions
func TransactionsHandler(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	dbconn := r.Context().Value("DBCONN").(*sql.DB)
	transactionsresp, err := findAllTransactions(dbconn, claims["id"])
	common.ValidateRenderResults(w, r, transactionsresp, err)
}

// TransactionCreateHandler POST /keys-transactions
func TransactionCreateHandler(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	dbconn := r.Context().Value("DBCONN").(*sql.DB)
	if queryValues.Get("key") != config.SteamBotsAPIKey() || !steam.IsOurSteamBot(dbconn, r.RemoteAddr) {
		render.Render(w, r, common.ErrInvalidRequest(errors.New("Unauthorized")))
		return
	}

	transaction := &TransactionsRequest{}
	if err := render.Bind(r, transaction); err != nil {
		render.Render(w, r, common.ErrInvalidRequest(err))
		return
	}

	if err := createTransaction(dbconn, transaction); err != nil {
		render.Render(w, r, common.ErrInternalServer(err))
		return
	}
	render.Render(w, r, common.SuccessCreatedResponse("Transaction has been created"))
}
