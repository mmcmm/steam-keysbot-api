package vault

import (
	"database/sql"

	"github.com/go-chi/render"
)

func findAllWithdrawals(dbconn *sql.DB, id interface{}) ([]render.Renderer, error) {
	rows, err := dbconn.Query(`SELECT status, payment_address, usd_rate, currency, usd_total, crypto_total,
		txhash, created_at FROM withdrawals WHERE user_steam_id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	withdrawalsresp := []render.Renderer{}
	for rows.Next() {
		resp := &WithdrawalsResponse{}
		if err := rows.Scan(
			&resp.Status,
			&resp.PaymentAddress,
			&resp.UsdRate,
			&resp.Currency,
			&resp.USDTotal,
			&resp.CryptoTotal,
			&resp.Txhash,
			&resp.CreatedAt,
		); err != nil {
			return nil, err
		}
		withdrawalsresp = append(withdrawalsresp, resp)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return withdrawalsresp, nil
}

// func saveWithdrawal(withdrawal interface{}, r *http.Request) {
// 	data := &WithdrawalsRequest{}
// 	if err := render.Bind(r, data); err != nil {
// 		render.Render(w, r, common.ErrInvalidRequest(err))
// 		return
// 	}
// 	withdrawal := data.Withdrawal
// 	dbNewArticle(article)
// }
