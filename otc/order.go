package otc

import (
	"context"
	"strconv"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// CreateOrderService -- POST /api/v4/otc/order/create (private)
//
// Creates a fiat OTC order. type is "BUY" (on-ramp) or "SELL" (off-ramp); side,
// crypto/fiat currencies and amounts must match a quote obtained from the quote
// endpoint, and quote_token / bank_id come from that quote and the bank-card list.
type CreateOrderService struct {
	c    *OTCClient
	body map[string]any
}

func (c *OTCClient) NewCreateOrderService(orderType, side, cryptoCurrency, fiatCurrency string, cryptoAmount, fiatAmount decimal.Decimal, quoteToken, bankID string) *CreateOrderService {
	return &CreateOrderService{c: c, body: map[string]any{
		"type":            orderType,
		"side":            side,
		"crypto_currency": cryptoCurrency,
		"fiat_currency":   fiatCurrency,
		"crypto_amount":   cryptoAmount.String(),
		"fiat_amount":     fiatAmount.String(),
		"quote_token":     quoteToken,
		"bank_id":         bankID,
	}}
}

// SetPromotionCode attaches an optional promotion code to the order.
func (s *CreateOrderService) SetPromotionCode(promotionCode string) *CreateOrderService {
	s.body["promotion_code"] = promotionCode
	return s
}

func (s *CreateOrderService) Do(ctx context.Context) (*OTCAckResponse, error) {
	req := request.Post(ctx, s.c, "/api/v4/otc/order/create", s.body).WithSign()
	return request.Do[OTCAckResponse](req)
}

// MarkOrderPaidService -- POST /api/v4/otc/order/paid (private)
//
// Marks a fiat order as paid (deposit confirmation), attaching the user's payment
// receipt stored as a file key.
type MarkOrderPaidService struct {
	c    *OTCClient
	body map[string]any
}

func (c *OTCClient) NewMarkOrderPaidService(orderID, paymentReceiptFileKey string) *MarkOrderPaidService {
	return &MarkOrderPaidService{c: c, body: map[string]any{
		"order_id":                 orderID,
		"payment_receipt_file_key": paymentReceiptFileKey,
	}}
}

// SetClientOrderID sets the client order id used by some gateway/Inner Pay paths.
func (s *MarkOrderPaidService) SetClientOrderID(clientOrderID string) *MarkOrderPaidService {
	s.body["client_order_id"] = clientOrderID
	return s
}

// SetPaymentReceipt sets the payment receipt alias compatible with
// payment_receipt_file_key (depends on the gateway's external field name).
func (s *MarkOrderPaidService) SetPaymentReceipt(paymentReceipt string) *MarkOrderPaidService {
	s.body["payment_receipt"] = paymentReceipt
	return s
}

func (s *MarkOrderPaidService) Do(ctx context.Context) (*OTCAckResponse, error) {
	req := request.Post(ctx, s.c, "/api/v4/otc/order/paid", s.body).WithSign()
	return request.Do[OTCAckResponse](req)
}

// CancelOrderService -- POST /api/v4/otc/order/cancel (private)
//
// Cancels a fiat OTC order by id. The order id travels in the query string even
// though the request is a POST.
type CancelOrderService struct {
	c       *OTCClient
	orderID string
}

func (c *OTCClient) NewCancelOrderService(orderID string) *CancelOrderService {
	return &CancelOrderService{c: c, orderID: orderID}
}

func (s *CancelOrderService) Do(ctx context.Context) (*OTCAckResponse, error) {
	req := request.Post(ctx, s.c, "/api/v4/otc/order/cancel")
	req.SetQuery("order_id", s.orderID)
	req.WithSign()
	return request.Do[OTCAckResponse](req)
}

// ListOrdersService -- GET /api/v4/otc/order/list (private)
//
// Lists the authenticated user's fiat OTC orders, optionally filtered by type,
// currency, time range or status.
type ListOrdersService struct {
	c      *OTCClient
	params map[string]string
}

func (c *OTCClient) NewListOrdersService() *ListOrdersService {
	return &ListOrdersService{c: c, params: map[string]string{}}
}

// SetType filters by order type: "BUY" (on-ramp) or "SELL" (off-ramp).
func (s *ListOrdersService) SetType(orderType string) *ListOrdersService {
	s.params["type"] = orderType
	return s
}

// SetFiatCurrency narrows the result to a single fiat currency.
func (s *ListOrdersService) SetFiatCurrency(fiatCurrency string) *ListOrdersService {
	s.params["fiat_currency"] = fiatCurrency
	return s
}

// SetCryptoCurrency narrows the result to a single digital currency.
func (s *ListOrdersService) SetCryptoCurrency(cryptoCurrency string) *ListOrdersService {
	s.params["crypto_currency"] = cryptoCurrency
	return s
}

// SetStartTime bounds the result to orders at or after this time.
func (s *ListOrdersService) SetStartTime(startTime string) *ListOrdersService {
	s.params["start_time"] = startTime
	return s
}

// SetEndTime bounds the result to orders at or before this time.
func (s *ListOrdersService) SetEndTime(endTime string) *ListOrdersService {
	s.params["end_time"] = endTime
	return s
}

// SetStatus narrows the result to a single order status ("DONE", "CANCEL",
// "PROCESSING" or "DISBURSED").
func (s *ListOrdersService) SetStatus(status string) *ListOrdersService {
	s.params["status"] = status
	return s
}

// SetPageNumber selects the result page.
func (s *ListOrdersService) SetPageNumber(pn int) *ListOrdersService {
	s.params["pn"] = strconv.Itoa(pn)
	return s
}

// SetPageSize sets the number of items per page.
func (s *ListOrdersService) SetPageSize(ps int) *ListOrdersService {
	s.params["ps"] = strconv.Itoa(ps)
	return s
}

func (s *ListOrdersService) Do(ctx context.Context) (*OTCOrderListResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/otc/order/list", s.params).WithSign()
	return request.Do[OTCOrderListResponse](req)
}

// OTCOrderListResponse is the envelope returned by the fiat order list endpoint.
type OTCOrderListResponse struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    OTCOrderList `json:"data"`
}

// OTCOrderList is one page of fiat OTC orders.
type OTCOrderList struct {
	Pn      int                `json:"pn"`
	Ps      int                `json:"ps"`
	TotalPn int                `json:"total_pn"`
	Count   int                `json:"count"`
	List    []OTCOrderListItem `json:"list"`
}

// OTCOrderListItem is a single fiat OTC order in a list. time is a formatted
// datetime string; timestamp is its integer epoch counterpart.
type OTCOrderListItem struct {
	Time               string          `json:"time"`
	Timestamp          int64           `json:"timestamp"`
	OrderID            string          `json:"order_id"`
	TradeNo            string          `json:"trade_no"`
	Type               string          `json:"type"`
	Status             string          `json:"status"`
	FiatCurrency       string          `json:"fiat_currency"`
	FiatCurrencyInfo   OTCCurrencyInfo `json:"fiat_currency_info"`
	FiatAmount         decimal.Decimal `json:"fiat_amount"`
	CryptoCurrency     string          `json:"crypto_currency"`
	CryptoCurrencyInfo OTCCurrencyInfo `json:"crypto_currency_info"`
	CryptoAmount       decimal.Decimal `json:"crypto_amount"`
	Rate               decimal.Decimal `json:"rate"`
	PromotionCode      string          `json:"promotion_code"`
}

// OTCCurrencyInfo is the display metadata (name and icon) for a currency.
type OTCCurrencyInfo struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
}

// GetOrderDetailService -- GET /api/v4/otc/order/detail (private)
//
// Returns the full detail of a single fiat OTC order by id.
type GetOrderDetailService struct {
	c      *OTCClient
	params map[string]string
}

func (c *OTCClient) NewGetOrderDetailService(orderID string) *GetOrderDetailService {
	return &GetOrderDetailService{c: c, params: map[string]string{
		"order_id": orderID,
	}}
}

func (s *GetOrderDetailService) Do(ctx context.Context) (*OTCOrderDetailResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/otc/order/detail", s.params).WithSign()
	return request.Do[OTCOrderDetailResponse](req)
}

// OTCOrderDetailResponse is the envelope returned by the fiat order detail
// endpoint. Timestamp is the server Unix time in seconds.
type OTCOrderDetailResponse struct {
	Code      int            `json:"code"`
	Message   string         `json:"message"`
	Data      OTCOrderDetail `json:"data"`
	Timestamp int64          `json:"timestamp"`
}

// OTCOrderDetail is the full detail of a fiat OTC order. create_time is a
// formatted datetime string. The bank_* fields carry the user's own bank
// transfer details and the gate_* fields carry Gate's receiving-bank details.
type OTCOrderDetail struct {
	OrderID        string          `json:"order_id"`
	UID            string          `json:"uid"`
	Type           string          `json:"type"`
	FiatCurrency   string          `json:"fiat_currency"`
	FiatAmount     decimal.Decimal `json:"fiat_amount"`
	CryptoCurrency string          `json:"crypto_currency"`
	CryptoAmount   decimal.Decimal `json:"crypto_amount"`
	Rate           decimal.Decimal `json:"rate"`

	BankAccountName           string `json:"bank_account_name"`
	BankName                  string `json:"bank_name"`
	BankCountry               string `json:"bank_country"`
	BankAddress               string `json:"bank_address"`
	BankAccountNumberIBAN     string `json:"bank_account_number_iban"`
	SwiftCode                 string `json:"swift_code"`
	IntermediateBankName      string `json:"intermediate_bank_name"`
	IntermediaryBankSwiftCode string `json:"intermediary_bank_swift_code"`

	GateBankAccountName           string `json:"gate_bank_account_name"`
	GateBankName                  string `json:"gate_bank_name"`
	GateBankCountry               string `json:"gate_bank_country"`
	GateBankAddress               string `json:"gate_bank_address"`
	GateBankAccountNumberIBAN     string `json:"gate_bank_account_number_iban"`
	GateSwiftCode                 string `json:"gate_swift_code"`
	GateIntermediaryBankName      string `json:"gate_intermediary_bank_name"`
	GateIntermediaryBankSwiftCode string `json:"gate_intermediary_bank_swift_code"`
	GateTransferRemark            string `json:"gate_transfer_remark"`
	GateReferenceCode             string `json:"gate_reference_code"`

	TransferRemark string `json:"transfer_remark"`
	ReferenceCode  string `json:"reference_code"`
	Status         string `json:"status"`
	DBStatus       string `json:"db_status"`
	CreateTime     string `json:"create_time"`
	Memo           string `json:"memo"`
	Side           string `json:"side"`
	PromotionCode  string `json:"promotion_code"`
	TradeNo        string `json:"trade_no"`
}
