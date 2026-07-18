package otc

import (
	"context"
	"strconv"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// CreateStableCoinOrderService -- POST /api/v4/otc/stable_coin/order/create (private)
//
// Creates a stablecoin OTC order. All fields except promotion_code are required;
// pay_coin / get_coin, the amounts, side and quote_token come from a prior
// stablecoin quote. side accepts "PAY"/"GET" for backward compatibility, but new
// integrations should pass the value returned by the quote response.
type CreateStableCoinOrderService struct {
	c    *OTCClient
	body map[string]any
}

func (c *OTCClient) NewCreateStableCoinOrderService(payCoin, getCoin string, payAmount, getAmount decimal.Decimal, side, quoteToken string) *CreateStableCoinOrderService {
	return &CreateStableCoinOrderService{c: c, body: map[string]any{
		"pay_coin":    payCoin,
		"get_coin":    getCoin,
		"pay_amount":  payAmount.String(),
		"get_amount":  getAmount.String(),
		"side":        side,
		"quote_token": quoteToken,
	}}
}

// SetPromotionCode attaches an optional promotion code to the order.
func (s *CreateStableCoinOrderService) SetPromotionCode(promotionCode string) *CreateStableCoinOrderService {
	s.body["promotion_code"] = promotionCode
	return s
}

func (s *CreateStableCoinOrderService) Do(ctx context.Context) (*OTCStableCoinOrderCreateResponse, error) {
	req := request.Post(ctx, s.c, "/api/v4/otc/stable_coin/order/create", s.body).WithSign()
	return request.Do[OTCStableCoinOrderCreateResponse](req)
}

// OTCStableCoinOrderCreateResponse is the acknowledgement envelope returned when
// creating a stablecoin order. Timestamp is the server Unix time in seconds.
type OTCStableCoinOrderCreateResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
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
	PayIcon      string          `json:"pay_icon"`
	PayAmount    decimal.Decimal `json:"pay_amount"`
	GetCoin      string          `json:"get_coin"`
	GetIcon      string          `json:"get_icon"`
	GetAmount    decimal.Decimal `json:"get_amount"`
	Rate         decimal.Decimal `json:"rate"`
	RateReci     decimal.Decimal `json:"rate_reci"`
	Status       string          `json:"status"`
	CreateTimest int64           `json:"create_timest"`
	CreateTime   string          `json:"create_time"`
}
