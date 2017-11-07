package internal

import (
	"encoding/json"
	"strconv"

	"github.com/mtdx/keyc/db"
	"github.com/mtdx/keyc/labels"
	"github.com/mtdx/keyc/pusher"
)

var dbconn = db.Open()

// InitLiveBtc ...
func InitLiveBtc() error {
	err := liveBtcPrice()
	if err != nil { // if we still can't connect
		SaveErr(E{"internal.InitLiveBtc", err.Error()})
		resetBtcPrice()
		return err
	}
	return nil
}

// LiveBtcPrice ...
func liveBtcPrice() error {
	pusherClient, err := pusher.NewClient("de504dc5763aeef9ff52") // https://www.bitstamp.net/websocket/
	if err != nil {
		return err
	}
	err = pusherClient.Subscribe("order_book")
	if err != nil {
		return err
	}
	dataChannelTrade, err := pusherClient.Bind("data")
	if err != nil {
		return err
	}

	go func() {
		for dataEvt := range dataChannelTrade {
			updateBtcPrice(dataEvt.Data)
		}
	}()

	return nil
}

func updateBtcPrice(data string) {
	orders := struct {
		Bids [][]string `json:"bids"`
		Asks [][]string `json:"asks"`
	}{}
	if err := json.Unmarshal([]byte(data), &orders); err != nil {
		SaveErr(E{"internal.updateBtcPrice.1", err.Error()})
		resetBtcPrice()
		return
	}
	a, err := strconv.ParseFloat(orders.Bids[0][0], 64)
	b, err := strconv.ParseFloat(orders.Asks[0][0], 64)
	if err != nil {
		SaveErr(E{"internal.updateBtcPrice.2", err.Error()})
		resetBtcPrice()
		return
	}
	_, err = dbconn.Exec(`UPDATE price_settings SET value = $1, updated_at = NOW() WHERE key = $2`,
		(a+b)/2,
		labels.BTC_USD_RATE,
	)
	if err != nil {
		SaveErr(E{"internal.updateBtcPrice.3", err.Error()})
		resetBtcPrice()
		return
	}
}

func resetBtcPrice() {
	_, err := dbconn.Exec(`UPDATE price_settings SET value = $1 WHERE key = $2`, 0, labels.BTC_USD_RATE)
	if err != nil {
		SaveErr(E{"internal.resetBtcPrice", err.Error()})
	}
}
