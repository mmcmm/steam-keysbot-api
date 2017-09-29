package handler

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/mtdx/keyc/common"
)

type accountResponse struct {
	BitcoinBalance int64 `json:"bitcoin_balance"`
	CsgokeyBalance int64 `json:"csgokey_balance"`
	TradeLinkURL   int64 `json:"trade_link_url"`
}

func (rd *accountResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Account route handler
func Account(w http.ResponseWriter, r *http.Request) {
	account := &accountResponse{}
	_, claims, _ := jwtauth.FromContext(r.Context())
	dbconn := r.Context().Value("DBCONN").(*sql.DB)
	err := dbconn.QueryRow(`SELECT (trade_link_url, bitcoin_balance, csgokey_balance) 
	FROM users WHERE steam_id = $1`, claims["id"]).Scan(&account.BitcoinBalance,
		&account.CsgokeyBalance, &account.TradeLinkURL)
	if err != nil || err == sql.ErrNoRows {
		render.Render(w, r, common.ErrInternalServer(err))
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, account)
}
