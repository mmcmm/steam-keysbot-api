package steam

import (
	"database/sql"

	"github.com/go-chi/render"
)

func findAllTradeoffers(dbconn *sql.DB, id interface{}) ([]render.Renderer, error) {
	rows, err := dbconn.Query(`SELECT type, status, failure_details, amount, created_at 
		FROM tradeoffers WHERE user_steam_id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tradeoffersresp := []render.Renderer{}
	for rows.Next() {
		resp := &TradeoffersResponse{}
		if err := rows.Scan(
			&resp.Type,
			&resp.Status,
			&resp.FailureDetails,
			&resp.Amount,
			&resp.CreatedAt,
		); err != nil {
			return nil, err
		}
		tradeoffersresp = append(tradeoffersresp, resp)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tradeoffersresp, nil
}
