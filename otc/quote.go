package otc

import (
	"context"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// QuoteService -- POST /api/v4/otc/quote (private)
//
// Requests a fiat or stablecoin quote. side is the quote direction: "PAY" means
// the user inputs the pay amount, "GET" means the user inputs the receive amount.
type QuoteService struct {
	c    *OTCClient
	body map[string]any
}

func (c *OTCClient) NewQuoteService(side, payCoin, getCoin string) *QuoteService {
	return &QuoteService{c: c, body: map[string]any{
		"side":     side,
		"pay_coin": payCoin,
		"get_coin": getCoin,
	}}
}

// SetPayAmount sets the amount of the currency the user pays.
func (s *QuoteService) SetPayAmount(payAmount decimal.Decimal) *QuoteService {
	s.body["pay_amount"] = payAmount.String()
	return s
}

// SetGetAmount sets the amount of the currency the user receives.
func (s *QuoteService) SetGetAmount(getAmount decimal.Decimal) *QuoteService {
	s.body["get_amount"] = getAmount.String()
	return s
}

// SetCreateQuoteToken controls token generation: "0" previews the quote only,
// "1" also mints a quote token that can be used to place an order.
func (s *QuoteService) SetCreateQuoteToken(createQuoteToken string) *QuoteService {
	s.body["create_quote_token"] = createQuoteToken
	return s
}

// SetPromotionCode attaches an optional promotion code to the quote.
func (s *QuoteService) SetPromotionCode(promotionCode string) *QuoteService {
	s.body["promotion_code"] = promotionCode
	return s
}

func (s *QuoteService) Do(ctx context.Context) (*OTCQuoteResponse, error) {
	req := request.Post(ctx, s.c, "/api/v4/otc/quote", s.body).WithSign()
	return request.Do[OTCQuoteResponse](req)
}

// OTCQuoteResponse is the envelope returned by the OTC quote endpoint.
type OTCQuoteResponse struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Data    OTCQuote `json:"data"`
}

// OTCQuote is a single fiat/stablecoin quote and its exchange terms.
type OTCQuote struct {
	Type          string          `json:"type"`
	PayCoin       string          `json:"pay_coin"`
	GetCoin       string          `json:"get_coin"`
	PayAmount     decimal.Decimal `json:"pay_amount"`
	GetAmount     decimal.Decimal `json:"get_amount"`
	Rate          decimal.Decimal `json:"rate"`
	RateReci      decimal.Decimal `json:"rate_reci"`
	PromotionCode string          `json:"promotion_code"`
	Side          string          `json:"side"`
	OrderType     string          `json:"order_type"`
	QuoteToken    string          `json:"quote_token"`
}
