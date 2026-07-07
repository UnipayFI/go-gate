package crossex

import (
	"context"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// CreateConvertQuoteService -- POST /api/v4/crossex/convert/quote (private)
//
// Requests a flash-swap quote for converting from_coin into to_coin on a venue.
type CreateConvertQuoteService struct {
	c    *CrossexClient
	body map[string]any
}

func (c *CrossexClient) NewCreateConvertQuoteService(exchangeType, fromCoin, toCoin string, fromAmount decimal.Decimal) *CreateConvertQuoteService {
	return &CreateConvertQuoteService{c: c, body: map[string]any{
		"exchange_type": exchangeType,
		"from_coin":     fromCoin,
		"to_coin":       toCoin,
		"from_amount":   fromAmount.String(),
	}}
}

func (s *CreateConvertQuoteService) Do(ctx context.Context) (*CrossexConvertQuote, error) {
	req := request.Post(ctx, s.c, "/api/v4/crossex/convert/quote", s.body).WithSign()
	return request.Do[CrossexConvertQuote](req)
}

// CrossexConvertQuote is a flash-swap quote. valid_ms is a millisecond Unix
// timestamp marking when the quote expires.
type CrossexConvertQuote struct {
	QuoteID    string          `json:"quote_id"`
	ValidMs    time.Time       `json:"valid_ms,string,format:unixmilli"`
	FromCoin   string          `json:"from_coin"`
	ToCoin     string          `json:"to_coin"`
	FromAmount decimal.Decimal `json:"from_amount"`
	ToAmount   decimal.Decimal `json:"to_amount"`
	Price      decimal.Decimal `json:"price"`
}

// CreateConvertOrderService -- POST /api/v4/crossex/convert/orders (private)
//
// Executes a flash-swap using a quote id returned by the quote endpoint.
type CreateConvertOrderService struct {
	c    *CrossexClient
	body map[string]any
}

func (c *CrossexClient) NewCreateConvertOrderService(quoteID string) *CreateConvertOrderService {
	return &CreateConvertOrderService{c: c, body: map[string]any{
		"quote_id": quoteID,
	}}
}

func (s *CreateConvertOrderService) Do(ctx context.Context) (*CrossexConvertOrderResult, error) {
	req := request.Post(ctx, s.c, "/api/v4/crossex/convert/orders", s.body).WithSign()
	return request.Do[CrossexConvertOrderResult](req)
}

// CrossexConvertOrderResult is the acknowledgement of an executed flash swap.
// The order id cannot be customized.
type CrossexConvertOrderResult struct {
	OrderID string `json:"order_id"`
	Text    string `json:"text"`
}
