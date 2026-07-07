package crossex

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// QueryTransferCoinsService -- GET /api/v4/crossex/transfers/coin (private)
//
// Returns the currencies supported for cross-exchange fund transfers.
type QueryTransferCoinsService struct {
	c      *CrossexClient
	params map[string]string
}

func (c *CrossexClient) NewQueryTransferCoinsService() *QueryTransferCoinsService {
	return &QueryTransferCoinsService{c: c, params: map[string]string{}}
}

// SetCoin narrows the result to a single currency.
func (s *QueryTransferCoinsService) SetCoin(coin string) *QueryTransferCoinsService {
	s.params["coin"] = coin
	return s
}

func (s *QueryTransferCoinsService) Do(ctx context.Context) ([]CrossexTransferCoin, error) {
	req := request.Get(ctx, s.c, "/api/v4/crossex/transfers/coin", s.params).WithSign()
	resp, err := request.Do[[]CrossexTransferCoin](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CrossexTransferCoin is a currency supported for cross-exchange transfers.
type CrossexTransferCoin struct {
	Coin           string          `json:"coin"`
	MinTransAmount decimal.Decimal `json:"min_trans_amount"`
	EstFee         decimal.Decimal `json:"est_fee"`
	Precision      int             `json:"precision"`
	IsDisabled     int             `json:"is_disabled"`
}

// ListTransfersService -- GET /api/v4/crossex/transfers (private)
//
// Returns the account's cross-exchange fund-transfer history.
type ListTransfersService struct {
	c      *CrossexClient
	params map[string]string
}

func (c *CrossexClient) NewListTransfersService() *ListTransfersService {
	return &ListTransfersService{c: c, params: map[string]string{}}
}

// SetCoin narrows the result to a single currency.
func (s *ListTransfersService) SetCoin(coin string) *ListTransfersService {
	s.params["coin"] = coin
	return s
}

// SetOrderID narrows the result to a transfer order id (tx_id) or custom id.
func (s *ListTransfersService) SetOrderID(orderID string) *ListTransfersService {
	s.params["order_id"] = orderID
	return s
}

// SetFrom sets the start time (millisecond precision).
func (s *ListTransfersService) SetFrom(from time.Time) *ListTransfersService {
	s.params["from"] = strconv.FormatInt(from.UnixMilli(), 10)
	return s
}

// SetTo sets the end time (millisecond precision), defaulting to now.
func (s *ListTransfersService) SetTo(to time.Time) *ListTransfersService {
	s.params["to"] = strconv.FormatInt(to.UnixMilli(), 10)
	return s
}

// SetPage selects the result page.
func (s *ListTransfersService) SetPage(page int) *ListTransfersService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit caps the number of records returned (max 1000).
func (s *ListTransfersService) SetLimit(limit int) *ListTransfersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *ListTransfersService) Do(ctx context.Context) ([]CrossexTransferRecord, error) {
	req := request.Get(ctx, s.c, "/api/v4/crossex/transfers", s.params).WithSign()
	resp, err := request.Do[[]CrossexTransferRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CrossexTransferRecord is one cross-exchange fund-transfer record. create_time
// and update_time are millisecond Unix timestamps.
type CrossexTransferRecord struct {
	ID              string          `json:"id"`
	Text            string          `json:"text"`
	FromAccountType string          `json:"from_account_type"`
	ToAccountType   string          `json:"to_account_type"`
	Coin            string          `json:"coin"`
	Amount          decimal.Decimal `json:"amount"`
	ActualReceive   decimal.Decimal `json:"actual_receive"`
	Status          string          `json:"status"`
	FailReason      string          `json:"fail_reason"`
	CreateTime      time.Time       `json:"create_time,format:unixmilli"`
	UpdateTime      time.Time       `json:"update_time,format:unixmilli"`
}

// CreateTransferService -- POST /api/v4/crossex/transfers (private)
//
// Transfers funds between cross-exchange credit accounts (from receiving account
// to debit account).
type CreateTransferService struct {
	c    *CrossexClient
	body map[string]any
}

func (c *CrossexClient) NewCreateTransferService(coin string, amount decimal.Decimal, from, to string) *CreateTransferService {
	return &CreateTransferService{c: c, body: map[string]any{
		"coin":   coin,
		"amount": amount.String(),
		"from":   from,
		"to":     to,
	}}
}

// SetText attaches a user-defined id to the transfer.
func (s *CreateTransferService) SetText(text string) *CreateTransferService {
	s.body["text"] = text
	return s
}

func (s *CreateTransferService) Do(ctx context.Context) (*CrossexTransferResult, error) {
	req := request.Post(ctx, s.c, "/api/v4/crossex/transfers", s.body).WithSign()
	return request.Do[CrossexTransferResult](req)
}

// CrossexTransferResult is the response of a fund transfer.
type CrossexTransferResult struct {
	TxID string `json:"tx_id"`
	Text string `json:"text"`
}
