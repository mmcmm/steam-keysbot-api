package internal

import (
	"database/sql"
	"testing"
	"time"

	"github.com/mtdx/assert"

	"github.com/mtdx/keyc/labels"
)

func TestBtcPriceUpdated(t *testing.T) {
	t.Parallel()

	InitLiveBtc()
	time.Sleep(2 * time.Second)

	var btcusdrate float64
	var updatedAt time.Time
	err := dbconn.QueryRow("SELECT value, updated_at FROM price_settings WHERE key = $1", labels.BTC_USD_RATE).Scan(
		&btcusdrate,
		&updatedAt,
	)
	if err != nil || err == sql.ErrNoRows {
		t.Fatalf("Failed to get price and updated, error: %s", err.Error())
	}
	assert.NotEqual(t, btcusdrate, 0)
	assert.True(t, time.Now().Unix()-4 < updatedAt.Unix())
}
