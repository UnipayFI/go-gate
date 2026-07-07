package otc

import (
	"context"
	"strconv"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// CreateStableCoinOrderService -- POST /api/v4/otc/stable_coin/order/create (private)
//
// Creates a stablecoin OTC order. Every field is optional on the wire; pay_coin /
// get_coin, one of the amounts, side and quote_token typically come from a prior
// stablecoin quote.
type CreateStableCoinOrderService struct {
	c    *OTCClient
	body map[string]any
}

func (c *OTCClient) NewCreateStableCoinOrderService() *CreateStableCoinOrderService {
	return &CreateStableCoinOrderService{c: c, body: map[string]any{}}
}

// SetPayCoin sets the currency the user pays.
func (s *CreateStableCoinOrderService) SetPayCoin(payCoin string) *CreateStableCoinOrderService {
	s.body["pay_coin"] = payCoin
	return s
}

// SetGetCoin sets the currency the user receives.
func (s *CreateStableCoinOrderService) SetGetCoin(getCoin string) *CreateStableCoinOrderService {
	s.body["get_coin"] = getCoin
	return s
}

// SetPayAmount sets the amount of the currency the user pays.
func (s *CreateStableCoinOrderService) SetPayAmount(payAmount decimal.Decimal) *CreateStableCoinOrderService {
	s.body["pay_amount"] = payAmount.String()
	return s
}

// SetGetAmount sets the amount of the currency the user receives.
func (s *CreateStableCoinOrderService) SetGetAmount(getAmount decimal.Decimal) *CreateStableCoinOrderService {
	s.body["get_amount"] = getAmount.String()
	return s
}

// SetSide sets the quote direction returned by the quote endpoint (used for
// order validation).
func (s *CreateStableCoinOrderService) SetSide(side string) *CreateStableCoinOrderService {
	s.body["side"] = side
	return s
}

// SetPromotionCode attaches an optional promotion code to the order.
func (s *CreateStableCoinOrderService) SetPromotionCode(promotionCode string) *CreateStableCoinOrderService {
	s.body["promotion_code"] = promotionCode
	return s
}

// SetQuoteToken sets the quote token returned by the quote endpoint.
func (s *CreateStableCoinOrderService) SetQuoteToken(quoteToken string) *CreateStableCoinOrderService {
	s.body["quote_token"] = quoteToken
	return s
}

func (s *CreateStableCoinOrderService) Do(ctx context.Context) (*OTCStableCoinOrderCreateResponse, error) {
	req := request.Post(ctx, s.c, "/api/v4/otc/stable_coin/order/create", s.body).WithSign()
	return request.Do[OTCStableCoinOrderCreateResponse](req)
}

// OTCStableCoinOrderCreateResponse is the acknowledgement envelope returned when
// creating a stablecoin order.
type OTCStableCoinOrderCreateResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ListStableCoinOrdersService -- GET /api/v4/otc/stable_coin/order/list (private)
//
// Lists the authenticated user's stablecoin OTC orders, optionally filtered by
// currency, time range or status.
type ListStableCoinOrdersService struct {
	c      *OTCClient
	params map[string]string
}

func (c *OTCClient) NewListStableCoinOrdersService() *ListStableCoinOrdersService {
	return &ListStableCoinOrdersService{c: c, params: map[string]string{}}
}

// SetPageSize sets the number of records per page.
func (s *ListStableCoinOrdersService) SetPageSize(pageSize int) *ListStableCoinOrdersService {
	s.params["page_size"] = strconv.Itoa(pageSize)
	return s
}

// SetPageNumber selects the result page.
func (s *ListStableCoinOrdersService) SetPageNumber(pageNumber int) *ListStableCoinOrdersService {
	s.params["page_number"] = strconv.Itoa(pageNumber)
	return s
}

// SetCoinName narrows the result to a single order currency.
func (s *ListStableCoinOrdersService) SetCoinName(coinName string) *ListStableCoinOrdersService {
	s.params["coin_name"] = coinName
	return s
}

// SetStartTime bounds the result to orders at or after this time.
func (s *ListStableCoinOrdersService) SetStartTime(startTime string) *ListStableCoinOrdersService {
	s.params["start_time"] = startTime
	return s
}

// SetEndTime bounds the result to orders at or before this time.
func (s *ListStableCoinOrdersService) SetEndTime(endTime string) *ListStableCoinOrdersService {
	s.params["end_time"] = endTime
	return s
}

// SetStatus narrows the result to a status ("PROCESSING", "DONE" or "FAILED").
func (s *ListStableCoinOrdersService) SetStatus(status string) *ListStableCoinOrdersService {
	s.params["status"] = status
	return s
}

func (s *ListStableCoinOrdersService) Do(ctx context.Context) (*OTCStableCoinOrderListResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/otc/stable_coin/order/list", s.params).WithSign()
	return request.Do[OTCStableCoinOrderListResponse](req)
}

// OTCStableCoinOrderListResponse is the envelope returned by the stablecoin order
// list endpoint.
type OTCStableCoinOrderListResponse struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    OTCStableCoinOrderList `json:"data"`
}

// OTCStableCoinOrderList is one page of stablecoin OTC orders.
type OTCStableCoinOrderList struct {
	Total      int                          `json:"total"`
	PageSize   int                          `json:"page_size"`
	PageNumber int                          `json:"page_number"`
	TotalPage  int                          `json:"total_page"`
	List       []OTCStableCoinOrderListItem `json:"list"`
}

// OTCStableCoinOrderListItem is a single stablecoin OTC order. create_time is a
// formatted datetime string; create_timest is its integer epoch counterpart.
type OTCStableCoinOrderListItem struct {
	ID           int64           `json:"id"`
	TradeNo      string          `json:"trade_no"`
	PayCoin      string          `json:"pay_coin"`
	PayAmount    decimal.Decimal `json:"pay_amount"`
	GetCoin      string          `json:"get_coin"`
	GetAmount    decimal.Decimal `json:"get_amount"`
	Rate         decimal.Decimal `json:"rate"`
	RateReci     decimal.Decimal `json:"rate_reci"`
	Status       string          `json:"status"`
	CreateTimest int64           `json:"create_timest"`
	CreateTime   string          `json:"create_time"`
}
