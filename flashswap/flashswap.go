package flashswap

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// ListFlashSwapCurrencyPairService -- GET /api/v4/flash_swap/currency_pairs
//
// Returns every trading pair supported for flash swap, optionally filtered to a
// single currency.
type ListFlashSwapCurrencyPairService struct {
	c      *FlashSwapClient
	params map[string]string
}

func (c *FlashSwapClient) NewListFlashSwapCurrencyPairService() *ListFlashSwapCurrencyPairService {
	return &ListFlashSwapCurrencyPairService{c: c, params: map[string]string{}}
}

// SetCurrency narrows the result to pairs involving the given currency.
func (s *ListFlashSwapCurrencyPairService) SetCurrency(currency string) *ListFlashSwapCurrencyPairService {
	s.params["currency"] = currency
	return s
}

// SetPage selects the page number.
func (s *ListFlashSwapCurrencyPairService) SetPage(page int) *ListFlashSwapCurrencyPairService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit caps the number of items returned (default 1000, max 1000).
func (s *ListFlashSwapCurrencyPairService) SetLimit(limit int) *ListFlashSwapCurrencyPairService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *ListFlashSwapCurrencyPairService) Do(ctx context.Context) ([]FlashSwapCurrencyPair, error) {
	req := request.Get(ctx, s.c, "/api/v4/flash_swap/currency_pairs", s.params)
	resp, err := request.Do[[]FlashSwapCurrencyPair](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// FlashSwapCurrencyPair is a trading pair supported for flash swap and its
// sell/buy amount limits.
type FlashSwapCurrencyPair struct {
	CurrencyPair  string          `json:"currency_pair"`
	SellCurrency  string          `json:"sell_currency"`
	BuyCurrency   string          `json:"buy_currency"`
	SellMinAmount decimal.Decimal `json:"sell_min_amount"`
	SellMaxAmount decimal.Decimal `json:"sell_max_amount"`
	BuyMinAmount  decimal.Decimal `json:"buy_min_amount"`
	BuyMaxAmount  decimal.Decimal `json:"buy_max_amount"`
}

// ListFlashSwapOrdersService -- GET /api/v4/flash_swap/orders (private)
//
// Returns the authenticated account's flash swap orders.
type ListFlashSwapOrdersService struct {
	c      *FlashSwapClient
	params map[string]string
}

func (c *FlashSwapClient) NewListFlashSwapOrdersService() *ListFlashSwapOrdersService {
	return &ListFlashSwapOrdersService{c: c, params: map[string]string{}}
}

// SetStatus filters by order status (1 = success, 2 = failure).
func (s *ListFlashSwapOrdersService) SetStatus(status int) *ListFlashSwapOrdersService {
	s.params["status"] = strconv.Itoa(status)
	return s
}

// SetSellCurrency filters by the asset sold.
func (s *ListFlashSwapOrdersService) SetSellCurrency(sellCurrency string) *ListFlashSwapOrdersService {
	s.params["sell_currency"] = sellCurrency
	return s
}

// SetBuyCurrency filters by the asset bought.
func (s *ListFlashSwapOrdersService) SetBuyCurrency(buyCurrency string) *ListFlashSwapOrdersService {
	s.params["buy_currency"] = buyCurrency
	return s
}

// SetReverse sorts by ID descending (true, default) or ascending (false).
func (s *ListFlashSwapOrdersService) SetReverse(reverse bool) *ListFlashSwapOrdersService {
	s.params["reverse"] = strconv.FormatBool(reverse)
	return s
}

// SetLimit caps the number of records returned in a single page.
func (s *ListFlashSwapOrdersService) SetLimit(limit int) *ListFlashSwapOrdersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetPage selects the page number.
func (s *ListFlashSwapOrdersService) SetPage(page int) *ListFlashSwapOrdersService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

func (s *ListFlashSwapOrdersService) Do(ctx context.Context) ([]FlashSwapOrder, error) {
	req := request.Get(ctx, s.c, "/api/v4/flash_swap/orders", s.params).WithSign()
	resp, err := request.Do[[]FlashSwapOrder](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CreateFlashSwapOrderService -- POST /api/v4/flash_swap/orders (private)
//
// Creates a flash swap order. A preview must be requested first, and its
// preview_id together with the previewed amounts are supplied here.
type CreateFlashSwapOrderService struct {
	c    *FlashSwapClient
	body map[string]any
}

func (c *FlashSwapClient) NewCreateFlashSwapOrderService(previewID, sellCurrency string, sellAmount decimal.Decimal, buyCurrency string, buyAmount decimal.Decimal) *CreateFlashSwapOrderService {
	return &CreateFlashSwapOrderService{c: c, body: map[string]any{
		"preview_id":    previewID,
		"sell_currency": sellCurrency,
		"sell_amount":   sellAmount.String(),
		"buy_currency":  buyCurrency,
		"buy_amount":    buyAmount.String(),
	}}
}

func (s *CreateFlashSwapOrderService) Do(ctx context.Context) (*FlashSwapOrder, error) {
	req := request.Post(ctx, s.c, "/api/v4/flash_swap/orders", s.body).WithSign()
	return request.Do[FlashSwapOrder](req)
}

// GetFlashSwapOrderService -- GET /api/v4/flash_swap/orders/{order_id} (private)
//
// Returns a single flash swap order by ID.
type GetFlashSwapOrderService struct {
	c       *FlashSwapClient
	orderID int64
}

func (c *FlashSwapClient) NewGetFlashSwapOrderService(orderID int64) *GetFlashSwapOrderService {
	return &GetFlashSwapOrderService{c: c, orderID: orderID}
}

func (s *GetFlashSwapOrderService) Do(ctx context.Context) (*FlashSwapOrder, error) {
	req := request.Get(ctx, s.c, "/api/v4/flash_swap/orders/"+strconv.FormatInt(s.orderID, 10)).WithSign()
	return request.Do[FlashSwapOrder](req)
}

// FlashSwapOrder is a flash swap order.
type FlashSwapOrder struct {
	ID           int64           `json:"id"`
	CreateTime   time.Time       `json:"create_time,format:unixmilli"`
	UserID       int64           `json:"user_id"`
	SellCurrency string          `json:"sell_currency"`
	SellAmount   decimal.Decimal `json:"sell_amount"`
	BuyCurrency  string          `json:"buy_currency"`
	BuyAmount    decimal.Decimal `json:"buy_amount"`
	Price        decimal.Decimal `json:"price"`
	Status       int             `json:"status"`
}

// PreviewFlashSwapOrderService -- POST /api/v4/flash_swap/orders/preview (private)
//
// Previews a flash swap order, returning the preview_id and quote used to create
// the order. Exactly one of sell_amount / buy_amount should be supplied.
type PreviewFlashSwapOrderService struct {
	c    *FlashSwapClient
	body map[string]any
}

func (c *FlashSwapClient) NewPreviewFlashSwapOrderService(sellCurrency, buyCurrency string) *PreviewFlashSwapOrderService {
	return &PreviewFlashSwapOrderService{c: c, body: map[string]any{
		"sell_currency": sellCurrency,
		"buy_currency":  buyCurrency,
	}}
}

// SetSellAmount sets the amount to sell (choose either sell_amount or buy_amount).
func (s *PreviewFlashSwapOrderService) SetSellAmount(sellAmount decimal.Decimal) *PreviewFlashSwapOrderService {
	s.body["sell_amount"] = sellAmount.String()
	return s
}

// SetBuyAmount sets the amount to buy (choose either sell_amount or buy_amount).
func (s *PreviewFlashSwapOrderService) SetBuyAmount(buyAmount decimal.Decimal) *PreviewFlashSwapOrderService {
	s.body["buy_amount"] = buyAmount.String()
	return s
}

func (s *PreviewFlashSwapOrderService) Do(ctx context.Context) (*FlashSwapPreview, error) {
	req := request.Post(ctx, s.c, "/api/v4/flash_swap/orders/preview", s.body).WithSign()
	return request.Do[FlashSwapPreview](req)
}

// FlashSwapPreview is the quote returned by a flash swap order preview.
type FlashSwapPreview struct {
	PreviewID    string          `json:"preview_id"`
	SellCurrency string          `json:"sell_currency"`
	SellAmount   decimal.Decimal `json:"sell_amount"`
	BuyCurrency  string          `json:"buy_currency"`
	BuyAmount    decimal.Decimal `json:"buy_amount"`
	Price        decimal.Decimal `json:"price"`
}
