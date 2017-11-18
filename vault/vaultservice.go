package vault

import (
	"database/sql"
	"errors"
	"math"
	"time"

	"github.com/mtdx/keyc/labels"

	"github.com/go-chi/render"
	"github.com/mtdx/keyc/common"
)

func findAllWithdrawals(dbconn *sql.DB, userid interface{}) ([]render.Renderer, error) {
	rows, err := dbconn.Query(`SELECT status, payment_address, usd_rate, currency, usd_total, crypto_total,
		txhash, created_at FROM withdrawals WHERE user_steam_id = $1 ORDER BY id DESC`, userid)
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

func saveWithdrawal(dbconn *sql.DB, withdrawal *WithdrawalsRequest, userid interface{}) error {
	return common.Transact(dbconn, func(tx *sql.Tx) error {
		var balance float64
		err := tx.QueryRow(`SELECT bitcoin_balance FROM users WHERE steam_id = $1 FOR UPDATE`, userid).Scan(&balance)
		if err != nil || err == sql.ErrNoRows {
			return err
		}
		if balance-math.Abs(withdrawal.CryptoTotal) < 0 {
			return errors.New("Not enough balance")
		}
		if _, err := tx.Exec(`UPDATE users SET bitcoin_balance = bitcoin_balance - $1 WHERE steam_id = $2`,
			math.Abs(withdrawal.CryptoTotal), userid); err != nil {
			return err
		}
		var btcrate float64
		var updatedAt time.Time
		err = tx.QueryRow(`SELECT value, updated_at FROM price_settings WHERE key = $1`,
			labels.BTC_USD_RATE).Scan(&btcrate, &updatedAt)
		if err != nil || err == sql.ErrNoRows {
			return err
		}
		if btcrate == 0 || updatedAt.Unix() < time.Now().Unix()-4 {
			return errors.New("Invalid BTC rate or not recently updated")
		}
		if _, err := tx.Exec(`INSERT INTO withdrawals (user_steam_id, payment_address, usd_rate, currency, 
			usd_total, crypto_total) VALUES ($1, $2, $3, $4, $5, $6)`,
			userid,
			withdrawal.PaymentAddress,
			btcrate,
			labels.BTC,
			withdrawal.CryptoTotal*btcrate,
			withdrawal.CryptoTotal,
		); err != nil {
			return err
		}
		return nil
	})
}
