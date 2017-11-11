package keys

import (
	"database/sql"
	"errors"
	"time"

	"github.com/mtdx/keyc/labels"

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

func createTransaction(dbconn *sql.DB, transaction *TransactionsRequest) error {
	var btcrate float64
	var unitprice float64
	var updatedAt time.Time
	err := dbconn.QueryRow("SELECT value, updated_at FROM price_settings WHERE key = $1", labels.BTC_USD_RATE).Scan(
		&btcrate, &updatedAt)
	err = dbconn.QueryRow("SELECT value FROM price_settings WHERE key = $1", transaction.TransactionType).Scan(&unitprice)
	if err != nil {
		return err
	}
	if time.Now().Unix()-2 > updatedAt.Unix() {
		return errors.New("Invalid BTC updated time")
	}
	if unitprice <= 0 {
		return errors.New("Invalid case key price")
	}
	_, err = dbconn.Exec(`INSERT INTO key_transactions (user_steam_id, tradeoffer_id, type, transaction_type, 
		amount, unit_price, payment_address, usd_rate, currency, usd_total, crypto_total, app_id) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		transaction.UserSteamID,
		transaction.TradeofferID,
		transaction.Type,
		transaction.TransactionType,
		transaction.Amount,
		unitprice,
		transaction.PaymentAddress,
		btcrate,
		transaction.Currency,
		unitprice*float64(transaction.Amount),
		(unitprice*float64(transaction.Amount))/btcrate,
		transaction.AppID,
	)

	return err
}
