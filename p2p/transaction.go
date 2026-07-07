package p2p

import (
	"context"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// GetPendingTransactionListService -- POST /api/v4/p2p/merchant/transaction/get_pending_transaction_list (private)
//
// Returns the merchant's pending (in-progress or disputed) P2P orders for a
// crypto/fiat pair.
type GetPendingTransactionListService struct {
	c    *P2PClient
	body map[string]any
}

func (c *P2PClient) NewGetPendingTransactionListService(cryptoCurrency, fiatCurrency string) *GetPendingTransactionListService {
	return &GetPendingTransactionListService{c: c, body: map[string]any{
		"crypto_currency": cryptoCurrency,
		"fiat_currency":   fiatCurrency,
	}}
}

// SetOrderTab selects the order tab: "pending" (OPEN, PAID, LOCKED, TEMP) or
// "dispute".
func (s *GetPendingTransactionListService) SetOrderTab(orderTab string) *GetPendingTransactionListService {
	s.body["order_tab"] = orderTab
	return s
}

// SetSelectType filters by order side: "buy", "sell", or empty for all.
func (s *GetPendingTransactionListService) SetSelectType(selectType string) *GetPendingTransactionListService {
	s.body["select_type"] = selectType
	return s
}

// SetStatus filters by order status ("open", "paid", "locked", "dispute").
func (s *GetPendingTransactionListService) SetStatus(status string) *GetPendingTransactionListService {
	s.body["status"] = status
	return s
}

// SetTxID narrows the result to a single order id.
func (s *GetPendingTransactionListService) SetTxID(txid int64) *GetPendingTransactionListService {
	s.body["txid"] = txid
	return s
}

// SetStartTime bounds the result to orders created at or after this time.
func (s *GetPendingTransactionListService) SetStartTime(startTime time.Time) *GetPendingTransactionListService {
	s.body["start_time"] = startTime.Unix()
	return s
}

// SetEndTime bounds the result to orders created at or before this time.
func (s *GetPendingTransactionListService) SetEndTime(endTime time.Time) *GetPendingTransactionListService {
	s.body["end_time"] = endTime.Unix()
	return s
}

func (s *GetPendingTransactionListService) Do(ctx context.Context) (*P2PResponse[P2PTransactionListData], error) {
	req := request.Post(ctx, s.c, "/api/v4/p2p/merchant/transaction/get_pending_transaction_list", s.body).WithSign()
	return doP2P[P2PTransactionListData](req)
}

// GetCompletedTransactionListService -- POST /api/v4/p2p/merchant/transaction/get_completed_transaction_list (private)
//
// Returns the merchant's completed (filled/canceled) P2P order history for a
// crypto/fiat pair, with pagination.
type GetCompletedTransactionListService struct {
	c    *P2PClient
	body map[string]any
}

func (c *P2PClient) NewGetCompletedTransactionListService(cryptoCurrency, fiatCurrency string) *GetCompletedTransactionListService {
	return &GetCompletedTransactionListService{c: c, body: map[string]any{
		"crypto_currency": cryptoCurrency,
		"fiat_currency":   fiatCurrency,
	}}
}

// SetSelectType filters by order side: "buy", "sell", or empty for all.
func (s *GetCompletedTransactionListService) SetSelectType(selectType string) *GetCompletedTransactionListService {
	s.body["select_type"] = selectType
	return s
}

// SetStatus filters by order status ("closed" for filled, "cancel" for canceled).
func (s *GetCompletedTransactionListService) SetStatus(status string) *GetCompletedTransactionListService {
	s.body["status"] = status
	return s
}

// SetTxID narrows the result to a single order id.
func (s *GetCompletedTransactionListService) SetTxID(txid int64) *GetCompletedTransactionListService {
	s.body["txid"] = txid
	return s
}

// SetStartTime bounds the result to orders at or after this time.
func (s *GetCompletedTransactionListService) SetStartTime(startTime time.Time) *GetCompletedTransactionListService {
	s.body["start_time"] = startTime.Unix()
	return s
}

// SetEndTime bounds the result to orders at or before this time.
func (s *GetCompletedTransactionListService) SetEndTime(endTime time.Time) *GetCompletedTransactionListService {
	s.body["end_time"] = endTime.Unix()
	return s
}

// SetQueryDispute flags dispute status in the response (1 yes, 0 no).
func (s *GetCompletedTransactionListService) SetQueryDispute(queryDispute int) *GetCompletedTransactionListService {
	s.body["query_dispute"] = queryDispute
	return s
}

// SetPage selects the result page (1-based).
func (s *GetCompletedTransactionListService) SetPage(page int) *GetCompletedTransactionListService {
	s.body["page"] = page
	return s
}

// SetPerPage caps the number of orders per page (default 10, max 200).
func (s *GetCompletedTransactionListService) SetPerPage(perPage int) *GetCompletedTransactionListService {
	s.body["per_page"] = perPage
	return s
}

func (s *GetCompletedTransactionListService) Do(ctx context.Context) (*P2PResponse[P2PTransactionListData], error) {
	req := request.Post(ctx, s.c, "/api/v4/p2p/merchant/transaction/get_completed_transaction_list", s.body).WithSign()
	return doP2P[P2PTransactionListData](req)
}

// P2PTransactionListData wraps a page of P2P order list items. Count is the total
// number of matching orders and ExportedNum the export count.
type P2PTransactionListData struct {
	List        []P2PTransactionListItem `json:"list"`
	Count       int                      `json:"count"`
	ExportedNum int                      `json:"exported_num"`
}

// P2PTransactionListItem is one P2P order row shared by the pending and completed
// order lists.
type P2PTransactionListItem struct {
	TypeBuy        int                        `json:"type_buy"`
	Timest         string                     `json:"timest"`
	TimestExpire   string                     `json:"timest_expire"`
	Timestamp      time.Time                  `json:"timestamp,format:unix"`
	Rate           decimal.Decimal            `json:"rate"`
	Amount         decimal.Decimal            `json:"amount"`
	Total          decimal.Decimal            `json:"total"`
	TxID           int64                      `json:"txid"`
	Status         string                     `json:"status"`
	ItsRealname    string                     `json:"its_realname"`
	ItsUID         string                     `json:"its_uid"`
	ItsNick        string                     `json:"its_nick"`
	SellerRealname string                     `json:"seller_realname"`
	BuyerRealname  string                     `json:"buyer_realname"`
	Cancelable     int                        `json:"cancelable"`
	CurrencyType   string                     `json:"currency_type"`
	WantType       string                     `json:"want_type"`
	HidePayment    int                        `json:"hide_payment"`
	SelPaytype     string                     `json:"sel_paytype"`
	PayOthers      []P2PPayMethodBrief        `json:"pay_others"`
	CdTime         int                        `json:"cd_time"`
	OrderType      int                        `json:"order_type"`
	OrderTag       []string                   `json:"order_tag"`
	ConvertInfo    P2PConvertInfo             `json:"convert_info"`
	TransTime      []P2PTransactionTimeMarker `json:"trans_time"`
}

// P2PPayMethodBrief is a short payment-method descriptor (type and display name)
// attached to an order.
type P2PPayMethodBrief struct {
	PayType string `json:"pay_type"`
	PayName string `json:"pay_name"`
}

// P2PConvertInfo holds flash-swap details for orders of that type.
type P2PConvertInfo struct {
	ConvertType   string          `json:"convert_type"`
	ConvertStatus string          `json:"convert_status"`
	PreRate       decimal.Decimal `json:"pre_rate"`
	Rate          decimal.Decimal `json:"rate"`
	PreFiatRate   decimal.Decimal `json:"pre_fiat_rate"`
	FiatRate      decimal.Decimal `json:"fiat_rate"`
	Amount        decimal.Decimal `json:"amount"`
	ConvertAmount decimal.Decimal `json:"convert_amount"`
	Slippage      decimal.Decimal `json:"slippage"`
	Status        string          `json:"status"`
}

// P2PTransactionTimeMarker is one countdown marker attached to an order.
type P2PTransactionTimeMarker struct {
	OdTime int `json:"od_time"`
}

// GetTransactionDetailsService -- POST /api/v4/p2p/merchant/transaction/get_transaction_details (private)
//
// Returns the full detail of a single P2P order, including counterparty and
// payment information.
type GetTransactionDetailsService struct {
	c    *P2PClient
	body map[string]any
}

func (c *P2PClient) NewGetTransactionDetailsService(txid int64) *GetTransactionDetailsService {
	return &GetTransactionDetailsService{c: c, body: map[string]any{
		"txid": txid,
	}}
}

// SetChannel sets the channel tag: empty for normal P2P, "web3" for Web3 orders.
func (s *GetTransactionDetailsService) SetChannel(channel string) *GetTransactionDetailsService {
	s.body["channel"] = channel
	return s
}

func (s *GetTransactionDetailsService) Do(ctx context.Context) (*P2PResponse[P2PTransactionDetail], error) {
	req := request.Post(ctx, s.c, "/api/v4/p2p/merchant/transaction/get_transaction_details", s.body).WithSign()
	return doP2P[P2PTransactionDetail](req)
}

// P2PTransactionDetail is the full detail of a single P2P order.
type P2PTransactionDetail struct {
	IsSell               int                  `json:"is_sell"`
	TxID                 int64                `json:"txid"`
	OrderID              int64                `json:"orderid"`
	Timest               time.Time            `json:"timest,format:unix"`
	LastPayTime          time.Time            `json:"last_pay_time,format:unix"`
	RemainPayTime        int                  `json:"remain_pay_time"`
	CurrencyType         string               `json:"currency_type"`
	WantType             string               `json:"want_type"`
	Symbol               string               `json:"symbol"`
	Rate                 decimal.Decimal      `json:"rate"`
	Amount               decimal.Decimal      `json:"amount"`
	Total                decimal.Decimal      `json:"total"`
	Status               string               `json:"status"`
	ReasonID             string               `json:"reason_id"`
	ReasonDesc           string               `json:"reason_desc"`
	CancelTime           string               `json:"cancel_time"`
	InAppeal             int                  `json:"in_appeal"`
	DisputeTime          time.Time            `json:"dispute_time,format:unix"`
	Cancelable           int                  `json:"cancelable"`
	HidePayment          int                  `json:"hide_payment"`
	TradeTips            string               `json:"trade_tips"`
	ShowBank             string               `json:"show_bank"`
	Bankname             string               `json:"bankname"`
	Bankbranch           string               `json:"bankbranch"`
	BankID               string               `json:"bankid"`
	BankHolderRealname   string               `json:"bank_holder_realname"`
	ShowAli              string               `json:"show_ali"`
	Aliname              string               `json:"aliname"`
	IsAlicode            int                  `json:"is_alicode"`
	ShowWechat           string               `json:"show_wechat"`
	Wename               string               `json:"wename"`
	ShowOthers           string               `json:"show_others"`
	PayOthers            []P2PDetailPayMethod `json:"pay_others"`
	SelPaytype           string               `json:"sel_paytype"`
	ItsUID               string               `json:"its_uid"`
	ItsNickname          string               `json:"its_nickname"`
	ItsRealname          string               `json:"its_realname"`
	HaveTraded           int                  `json:"have_traded"`
	AppealAllowCancel    int                  `json:"appeal_allow_cancel"`
	AppealVerdictHasOpen string               `json:"appeal_verdict_has_open"`
	IMUnread             int                  `json:"im_unread"`
	PaymentVoucherURL    []string             `json:"payment_voucher_url"`
	TimestPaid           time.Time            `json:"timest_paid,format:unix"`
	OwnRealname          string               `json:"own_realname"`
	OrderType            int                  `json:"order_type"`
	IsShowReceive        int                  `json:"is_show_receive"`
	ShowSellerContact    bool                 `json:"show_seller_contact_info"`
	SupportedPayTypes    []string             `json:"supported_pay_types"`
}

// P2PDetailPayMethod is one payment-method entry in an order detail's pay_others.
type P2PDetailPayMethod struct {
	ID         string `json:"id"`
	AccountDes string `json:"account_des"`
	PayType    string `json:"pay_type"`
	Account    string `json:"account"`
	Memo       string `json:"memo"`
	TradeTips  string `json:"trade_tips"`
	PayName    string `json:"pay_name"`
}

// ConfirmPaymentService -- POST /api/v4/p2p/merchant/transaction/confirm-payment (private)
//
// Marks a P2P order as paid by the buyer.
type ConfirmPaymentService struct {
	c    *P2PClient
	body map[string]any
}

func (c *P2PClient) NewConfirmPaymentService(txid string) *ConfirmPaymentService {
	return &ConfirmPaymentService{c: c, body: map[string]any{
		"txid": txid,
	}}
}

// SetPaymentMethod sets the payment type used (must be among the order-supported
// types).
func (s *ConfirmPaymentService) SetPaymentMethod(paymentMethod string) *ConfirmPaymentService {
	s.body["payment_method"] = paymentMethod
	return s
}

func (s *ConfirmPaymentService) Do(ctx context.Context) (*P2PResponse[P2PEmptyData], error) {
	req := request.Post(ctx, s.c, "/api/v4/p2p/merchant/transaction/confirm-payment", s.body).WithSign()
	return doP2P[P2PEmptyData](req)
}

// ConfirmReceiptService -- POST /api/v4/p2p/merchant/transaction/confirm-receipt (private)
//
// Confirms receipt of payment on a P2P order, releasing the crypto to the buyer.
type ConfirmReceiptService struct {
	c    *P2PClient
	body map[string]any
}

func (c *P2PClient) NewConfirmReceiptService(txid string) *ConfirmReceiptService {
	return &ConfirmReceiptService{c: c, body: map[string]any{
		"txid": txid,
	}}
}

func (s *ConfirmReceiptService) Do(ctx context.Context) (*P2PResponse[P2PEmptyData], error) {
	req := request.Post(ctx, s.c, "/api/v4/p2p/merchant/transaction/confirm-receipt", s.body).WithSign()
	return doP2P[P2PEmptyData](req)
}

// CancelOrderService -- POST /api/v4/p2p/merchant/transaction/cancel (private)
//
// Cancels a P2P order, optionally with a cancellation reason.
type CancelOrderService struct {
	c    *P2PClient
	body map[string]any
}

func (c *P2PClient) NewCancelOrderService(txid string) *CancelOrderService {
	return &CancelOrderService{c: c, body: map[string]any{
		"txid": txid,
	}}
}

// SetReasonID sets the cancel reason id (see Gate's reason list).
func (s *CancelOrderService) SetReasonID(reasonID string) *CancelOrderService {
	s.body["reason_id"] = reasonID
	return s
}

// SetReasonMemo attaches extra cancel notes (required when reason_id is 9).
func (s *CancelOrderService) SetReasonMemo(reasonMemo string) *CancelOrderService {
	s.body["reason_memo"] = reasonMemo
	return s
}

func (s *CancelOrderService) Do(ctx context.Context) (*P2PResponse[P2PEmptyData], error) {
	req := request.Post(ctx, s.c, "/api/v4/p2p/merchant/transaction/cancel", s.body).WithSign()
	return doP2P[P2PEmptyData](req)
}
