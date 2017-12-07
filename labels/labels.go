package labels

// type
const (
	_         = iota
	SELLTOUS  // sell to us
	BUYFROMUS // buy from us
)

// games
const (
	CSGO_APP_ID = 730
)

// status
const (
	_ = iota
	INVALID
	ACTIVE
	ACCEPTED
	COUNTERED
	EXPIRED
	CANCELLED
	DECLINED
	INVALIDITEMS
	CREATEDNEEDSCONFIRMATION
	CANCELEDBYSECONDFACTOR
	INESCROW
	COMPLETED
	UNPAID
	PENDING
)

// currency
const (
	_ = iota
	BTC
	BCH
)

// settings
const (
	_ = iota
	BTC_USD_RATE
	BUY_CSGOKEY_PRICE
	SELL_CSGOKEY_PRICE
)
