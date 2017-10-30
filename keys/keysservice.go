package keys

import (
	"database/sql"

	"github.com/go-chi/render"
)

func findAllTransactions(dbconn *sql.DB, id interface{}) ([]render.Renderer, error) {
	rows, err := dbconn.Query(`SELECT status, type, amount, unit_price, payment_address, usd_rate, 
		currency, usd_total, crypto_total, created_at FROM key_transactions WHERE user_steam_id = $1 ORDER BY id DESC`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	transactionsresp := []render.Renderer{}
	for rows.Next() {
		resp := &TransactionsResponse{}
		if err := rows.Scan(
			&resp.Status,
			&resp.Type,
			&resp.Amount,
			&resp.UnitPrice,
			&resp.PaymentAddress,
			&resp.USDRate,
			&resp.Currency,
			&resp.USDTotal,
			&resp.CryptoTotal,
			&resp.CreatedAt,
		); err != nil {
			return nil, err
		}
		transactionsresp = append(transactionsresp, resp)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactionsresp, nil
}
