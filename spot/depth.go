package spot

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/request"
	"github.com/shopspring/decimal"
)

// ListOrderBookService -- GET /api/v4/spot/order_book
//
// Returns the current market depth for a pair: bids sorted high-to-low and asks
// low-to-high.
type ListOrderBookService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewListOrderBookService(currencyPair string) *ListOrderBookService {
	return &ListOrderBookService{c: c, params: map[string]string{"currency_pair": currencyPair}}
}

// SetInterval sets the price precision used to aggregate depth levels
// ("0" means no aggregation).
func (s *ListOrderBookService) SetInterval(interval string) *ListOrderBookService {
	s.params["interval"] = interval
	return s
}

// SetLimit caps the number of depth levels returned per side.
func (s *ListOrderBookService) SetLimit(limit int) *ListOrderBookService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetWithID requests the order book update ID (populates the id/current/update
// fields).
func (s *ListOrderBookService) SetWithID(withID bool) *ListOrderBookService {
	s.params["with_id"] = strconv.FormatBool(withID)
	return s
}

func (s *ListOrderBookService) Do(ctx context.Context) (*OrderBook, error) {
	req := request.Get(ctx, s.c, "/api/v4/spot/order_book", s.params)
	return request.Do[OrderBook](req)
}

// OrderBook is a snapshot of market depth. Asks and Bids are [price, amount]
// pairs. ID is valid only when with_id is set; Current/Update are millisecond
// epochs.
type OrderBook struct {
	ID      int64               `json:"id"`
	Current time.Time           `json:"current,format:unixmilli"`
	Update  time.Time           `json:"update,format:unixmilli"`
	Asks    [][]decimal.Decimal `json:"asks"`
	Bids    [][]decimal.Decimal `json:"bids"`
}

// ListTradesService -- GET /api/v4/spot/trades
//
// Returns recent market trades for a pair, most recent first.
type ListTradesService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewListTradesService(currencyPair string) *ListTradesService {
	return &ListTradesService{c: c, params: map[string]string{"currency_pair": currencyPair}}
}

// SetLimit caps the number of trades returned (default 100, max 1000).
func (s *ListTradesService) SetLimit(limit int) *ListTradesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetLastID uses the id of the last record from the previous page as the cursor
// for the next page.
func (s *ListTradesService) SetLastID(lastID string) *ListTradesService {
	s.params["last_id"] = lastID
	return s
}

// SetReverse walks back to trades with an id less than last_id when true.
func (s *ListTradesService) SetReverse(reverse bool) *ListTradesService {
	s.params["reverse"] = strconv.FormatBool(reverse)
	return s
}

// SetFrom sets the start of the query time range.
func (s *ListTradesService) SetFrom(from time.Time) *ListTradesService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo sets the end of the query time range (defaults to now).
func (s *ListTradesService) SetTo(to time.Time) *ListTradesService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

// SetPage selects the page number for limit/page pagination.
func (s *ListTradesService) SetPage(page int) *ListTradesService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

func (s *ListTradesService) Do(ctx context.Context) ([]MarketTrade, error) {
	req := request.Get(ctx, s.c, "/api/v4/spot/trades", s.params)
	resp, err := request.Do[[]MarketTrade](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// MarketTrade is a single public market trade. The role/text/order_id and
// fee-related fields are only populated on authenticated queries; public
// responses omit them.
type MarketTrade struct {
	ID           string          `json:"id"`
	CreateTime   time.Time       `json:"create_time,string,format:unix"`
	CreateTimeMs time.Time       `json:"create_time_ms,string,format:unixmilli"`
	CurrencyPair string          `json:"currency_pair"`
	Side         Side            `json:"side"`
	Role         string          `json:"role"`
	Amount       decimal.Decimal `json:"amount"`
	Price        decimal.Decimal `json:"price"`
	OrderID      string          `json:"order_id"`
	Fee          decimal.Decimal `json:"fee"`
	FeeCurrency  string          `json:"fee_currency"`
	PointFee     decimal.Decimal `json:"point_fee"`
	GTFee        decimal.Decimal `json:"gt_fee"`
	AmendText    string          `json:"amend_text"`
	SequenceID   string          `json:"sequence_id"`
	Text         string          `json:"text"`
}

// ListCandlesticksService -- GET /api/v4/spot/candlesticks
//
// Returns OHLC candlestick data for a pair. At most 1000 points per query.
type ListCandlesticksService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewListCandlesticksService(currencyPair string) *ListCandlesticksService {
	return &ListCandlesticksService{c: c, params: map[string]string{"currency_pair": currencyPair}}
}

// SetLimit caps the number of recent points returned (conflicts with from/to).
func (s *ListCandlesticksService) SetLimit(limit int) *ListCandlesticksService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetFrom sets the start time of the candlestick range.
func (s *ListCandlesticksService) SetFrom(from time.Time) *ListCandlesticksService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo sets the end time of the candlestick range (defaults to now).
func (s *ListCandlesticksService) SetTo(to time.Time) *ListCandlesticksService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

// SetInterval selects the time interval between points (e.g. "1m", "1h", "1d";
// "30d" is a calendar month).
func (s *ListCandlesticksService) SetInterval(interval string) *ListCandlesticksService {
	s.params["interval"] = interval
	return s
}

func (s *ListCandlesticksService) Do(ctx context.Context) ([]Candlestick, error) {
	req := request.Get(ctx, s.c, "/api/v4/spot/candlesticks", s.params)
	resp, err := request.Do[[]Candlestick](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// Candlestick is a single OHLC point. Gate encodes each point as an array of
// eight strings, decoded here by fixed position:
// [unix_seconds, quote_volume, close, high, low, open, base_volume, window_closed].
type Candlestick struct {
	Timestamp    time.Time
	QuoteVolume  decimal.Decimal
	Close        decimal.Decimal
	High         decimal.Decimal
	Low          decimal.Decimal
	Open         decimal.Decimal
	BaseVolume   decimal.Decimal
	WindowClosed bool
}

func (c *Candlestick) UnmarshalJSON(data []byte) error {
	var cols []string
	if err := json.Unmarshal(data, &cols); err != nil {
		return err
	}
	if len(cols) < 8 {
		return fmt.Errorf("gate: candlestick expects 8 columns, got %d", len(cols))
	}
	sec, err := strconv.ParseInt(cols[0], 10, 64)
	if err != nil {
		return fmt.Errorf("gate: candlestick timestamp %q: %w", cols[0], err)
	}
	c.Timestamp = time.Unix(sec, 0)
	for _, f := range []struct {
		dst *decimal.Decimal
		src string
	}{
		{&c.QuoteVolume, cols[1]},
		{&c.Close, cols[2]},
		{&c.High, cols[3]},
		{&c.Low, cols[4]},
		{&c.Open, cols[5]},
		{&c.BaseVolume, cols[6]},
	} {
		v, err := decimal.NewFromString(f.src)
		if err != nil {
			return fmt.Errorf("gate: candlestick decimal %q: %w", f.src, err)
		}
		*f.dst = v
	}
	closed, err := strconv.ParseBool(cols[7])
	if err != nil {
		return fmt.Errorf("gate: candlestick window_closed %q: %w", cols[7], err)
	}
	c.WindowClosed = closed
	return nil
}
