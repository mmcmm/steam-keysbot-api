package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/mtdx/keyc/db"
	"github.com/mtdx/keyc/labels"
	"github.com/mtdx/keyc/pusher"
)

var dbconn = db.Open()

func init() {
	err := liveBtcPrice()
	for i := 0; err != nil && i < 30; i++ {
		time.Sleep(2 * time.Second)
		err = liveBtcPrice() // retry 30 times
	}

	if err != nil { // if we still can't connect
		fmt.Printf("Error whilde decoding %v\n", err) // TODO: log and set price to 0
		// TODO: log error here
		os.Exit(1)
	}
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
		fmt.Printf("Error whilde decoding %v\n", err) // TODO:
		return
	}
	a, err := strconv.ParseFloat(orders.Bids[0][0], 64)
	b, err := strconv.ParseFloat(orders.Asks[0][0], 64)
	if err != nil {
		fmt.Printf("Error whilde decoding %v\n", err) // TODO:
		return
	}
	_, err = dbconn.Exec(`UPDATE price_settings SET value = $1, updated_at = NOW() WHERE key = $2`,
		(a+b)/2,
		labels.BTC_USD_RATE,
	)
	if err != nil {
		fmt.Printf("Error saving btc price %v\n", err) // TODO:
		return
	}
}
