package spot

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// GetFeeService -- GET /api/v4/spot/fee (private)
//
// Returns the personal trading fee rates, optionally scoped to a single pair.
type GetFeeService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetFeeService() *GetFeeService {
	return &GetFeeService{c: c, params: map[string]string{}}
}

// SetCurrencyPair returns the fee rates for a specific trading pair.
func (s *GetFeeService) SetCurrencyPair(currencyPair string) *GetFeeService {
	s.params["currency_pair"] = currencyPair
	return s
}

func (s *GetFeeService) Do(ctx context.Context) (*SpotFee, error) {
	req := request.Get(ctx, s.c, "/api/v4/spot/fee", s.params).WithSign()
	return request.Do[SpotFee](req)
}

// GetBatchSpotFeeService -- GET /api/v4/spot/batch_fee (private)
//
// Returns the personal trading fee rates for multiple pairs at once (max 50).
type GetBatchSpotFeeService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetBatchSpotFeeService(currencyPairs string) *GetBatchSpotFeeService {
	return &GetBatchSpotFeeService{c: c, params: map[string]string{"currency_pairs": currencyPairs}}
}

func (s *GetBatchSpotFeeService) Do(ctx context.Context) (map[string]SpotFee, error) {
	req := request.Get(ctx, s.c, "/api/v4/spot/batch_fee", s.params).WithSign()
	resp, err := request.Do[map[string]SpotFee](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// SpotFee is the personal trading fee schedule for a currency pair.
type SpotFee struct {
	UserID       int64           `json:"user_id"`
	TakerFee     decimal.Decimal `json:"taker_fee"`
	MakerFee     decimal.Decimal `json:"maker_fee"`
	GTDiscount   bool            `json:"gt_discount"`
	GTTakerFee   decimal.Decimal `json:"gt_taker_fee"`
	GTMakerFee   decimal.Decimal `json:"gt_maker_fee"`
	LoanFee      decimal.Decimal `json:"loan_fee"`
	PointType    string          `json:"point_type"`
	CurrencyPair string          `json:"currency_pair"`
	DebitFee     int32           `json:"debit_fee"`
	RPIMakerFee  decimal.Decimal `json:"rpi_maker_fee"`
	RPIMM        decimal.Decimal `json:"rpi_mm"`
}

// ListSpotAccountBookService -- GET /api/v4/spot/account_book (private)
//
// Returns the spot account balance-change history. The queried range may not
// exceed 30 days.
type ListSpotAccountBookService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewListSpotAccountBookService() *ListSpotAccountBookService {
	return &ListSpotAccountBookService{c: c, params: map[string]string{}}
}

// SetCurrency filters the history to a single currency.
func (s *ListSpotAccountBookService) SetCurrency(currency string) *ListSpotAccountBookService {
	s.params["currency"] = currency
	return s
}

// SetFrom sets the start of the query window.
func (s *ListSpotAccountBookService) SetFrom(from time.Time) *ListSpotAccountBookService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo sets the end of the query window (defaults to now when unset).
func (s *ListSpotAccountBookService) SetTo(to time.Time) *ListSpotAccountBookService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

// SetPage selects the page number when paginating.
func (s *ListSpotAccountBookService) SetPage(page int) *ListSpotAccountBookService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit caps the number of records returned in a single response.
func (s *ListSpotAccountBookService) SetLimit(limit int) *ListSpotAccountBookService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetType filters by account-book change type.
func (s *ListSpotAccountBookService) SetType(bookType string) *ListSpotAccountBookService {
	s.params["type"] = bookType
	return s
}

// SetCode filters by account change code. It takes priority over the type filter.
func (s *ListSpotAccountBookService) SetCode(code string) *ListSpotAccountBookService {
	s.params["code"] = code
	return s
}

func (s *ListSpotAccountBookService) Do(ctx context.Context) ([]SpotAccountBook, error) {
	req := request.Get(ctx, s.c, "/api/v4/spot/account_book", s.params).WithSign()
	resp, err := request.Do[[]SpotAccountBook](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// SpotAccountBook is a single spot balance-change record.
type SpotAccountBook struct {
	ID       string          `json:"id"`
	Time     time.Time       `json:"time,format:unixmilli"`
	Currency string          `json:"currency"`
	Change   decimal.Decimal `json:"change"`
	Balance  decimal.Decimal `json:"balance"`
	Type     string          `json:"type"`
	Code     string          `json:"code"`
	Text     string          `json:"text"`
}

// ListMyTradesService -- GET /api/v4/spot/my_trades (private)
//
// Returns the authenticated account's personal trade fills. Without a time
// range only the last 7 days are available, and any range may not exceed 30 days.
type ListMyTradesService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewListMyTradesService() *ListMyTradesService {
	return &ListMyTradesService{c: c, params: map[string]string{}}
}

// SetCurrencyPair filters the fills to a single trading pair.
func (s *ListMyTradesService) SetCurrencyPair(currencyPair string) *ListMyTradesService {
	s.params["currency_pair"] = currencyPair
	return s
}

// SetLimit caps the number of fills returned (default 100, max 1000).
func (s *ListMyTradesService) SetLimit(limit int) *ListMyTradesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetPage selects the page number when paginating.
func (s *ListMyTradesService) SetPage(page int) *ListMyTradesService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetOrderID filters fills to a specific order. SetCurrencyPair is also required
// when this is set.
func (s *ListMyTradesService) SetOrderID(orderID string) *ListMyTradesService {
	s.params["order_id"] = orderID
	return s
}

// SetAccount selects which account's fills to query.
func (s *ListMyTradesService) SetAccount(account Account) *ListMyTradesService {
	s.params["account"] = string(account)
	return s
}

// SetFrom sets the start of the query window.
func (s *ListMyTradesService) SetFrom(from time.Time) *ListMyTradesService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo sets the end of the query window (defaults to now when unset).
func (s *ListMyTradesService) SetTo(to time.Time) *ListMyTradesService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

func (s *ListMyTradesService) Do(ctx context.Context) ([]MyTrade, error) {
	req := request.Get(ctx, s.c, "/api/v4/spot/my_trades", s.params).WithSign()
	resp, err := request.Do[[]MyTrade](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// MyTrade is a single personal spot trade fill.
type MyTrade struct {
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
