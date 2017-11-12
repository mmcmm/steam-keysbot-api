package steam

import (
	"database/sql"
	"strings"

	"github.com/go-chi/render"
)

func findAllTradeoffers(dbconn *sql.DB, id interface{}) ([]render.Renderer, error) {
	rows, err := dbconn.Query(`SELECT type, status, failure_details, amount, created_at 
		FROM tradeoffers WHERE user_steam_id = $1 ORDER BY tradeoffer_id DESC`, id)
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

func saveTradeoffer(dbconn *sql.DB, tradeoffer *TradeoffersRequest) error {
	_, err := dbconn.Exec(`INSERT INTO tradeoffers (tradeoffer_id, user_steam_id, type, amount, app_id)
	VALUES ($1, $2, $3, $4, $5)`,
		tradeoffer.TradeofferID,
		tradeoffer.SteamID,
		tradeoffer.Type,
		tradeoffer.Amount,
		tradeoffer.AppID,
	)

	return err
}

func updateStatus(dbconn *sql.DB, tradeoffer *TradeoffersUpdateRequest, tradeofferID string) error {
	_, err := dbconn.Exec(`UPDATE tradeoffers SET status = $2, failure_details = $3 WHERE tradeoffer_id = $1`,
		tradeofferID,
		tradeoffer.Status,
		tradeoffer.FailureDetails,
	)

	return err
}

// IsOurSteamBot ...
func IsOurSteamBot(dbconn *sql.DB, ip string) bool {
	var dbip string
	ip = ip[0:strings.Index(ip, ":")]
	err := dbconn.QueryRow("SELECT ip_address FROM steam_bots WHERE ip_address = $1", ip).Scan(&dbip)
	if err != nil || err == sql.ErrNoRows {
		return false
	}

	return dbip == ip
}
