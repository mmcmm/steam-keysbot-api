package account

import (
	"database/sql"

	"github.com/go-chi/render"
	"github.com/mtdx/keyc/validator"
)

func getAccountInfo(dbconn *sql.DB, id interface{}) (render.Renderer, error) {
	inforesp := &InfoResponse{}
	err := dbconn.QueryRow(`SELECT bitcoin_balance, csgokey_balance, trade_link_url 
		FROM users WHERE steam_id = $1`, id).Scan(
		&inforesp.BitcoinBalance,
		&inforesp.CsgokeyBalance,
		&inforesp.TradeLinkURL,
	)
	if err != nil || err == sql.ErrNoRows {
		return nil, err
	}

	if err := validator.Validate(inforesp); err != nil {
		return nil, err
	}

	return inforesp, nil
}
