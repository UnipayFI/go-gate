package p2p

import (
	"context"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// PlaceBizPushOrderService -- POST /api/v4/p2p/merchant/books/place_biz_push_order (private)
//
// Publishes or edits a P2P advertisement. adType selects the operation
// (0 publish sell, 1 publish buy, 2 edit sell, 3 edit buy).
type PlaceBizPushOrderService struct {
	c    *P2PClient
	body map[string]any
}

func (c *P2PClient) NewPlaceBizPushOrderService(currencyType, exchangeType, adType string, unitPrice, number decimal.Decimal, payType string) *PlaceBizPushOrderService {
	return &PlaceBizPushOrderService{c: c, body: map[string]any{
		"currencyType": currencyType,
		"exchangeType": exchangeType,
		"type":         adType,
		"unitPrice":    unitPrice.String(),
		"number":       number.String(),
		"payType":      payType,
	}}
}

// SetPayTypeJSON sets a JSON map of payment type -> user's payment method id.
func (s *PlaceBizPushOrderService) SetPayTypeJSON(payTypeJSON string) *PlaceBizPushOrderService {
	s.body["pay_type_json"] = payTypeJSON
	return s
}

// SetRateFixed sets the price type: "0" floating, "1" fixed.
func (s *PlaceBizPushOrderService) SetRateFixed(rateFixed string) *PlaceBizPushOrderService {
	s.body["rateFixed"] = rateFixed
	return s
}

// SetOID sets the advertisement id when editing (omit when publishing).
func (s *PlaceBizPushOrderService) SetOID(oid string) *PlaceBizPushOrderService {
	s.body["oid"] = oid
	return s
}

// SetMinAmount sets the minimum quantity per order (in currencyType).
func (s *PlaceBizPushOrderService) SetMinAmount(minAmount decimal.Decimal) *PlaceBizPushOrderService {
	s.body["minAmount"] = minAmount.String()
	return s
}

// SetMaxAmount sets the maximum quantity per order (in currencyType).
func (s *PlaceBizPushOrderService) SetMaxAmount(maxAmount decimal.Decimal) *PlaceBizPushOrderService {
	s.body["maxAmount"] = maxAmount.String()
	return s
}

// SetLimitBasis sets the trading-limit unit: 0 by crypto quantity, 1 by fiat.
func (s *PlaceBizPushOrderService) SetLimitBasis(limitBasis int) *PlaceBizPushOrderService {
	s.body["limitBasis"] = limitBasis
	return s
}

// SetFiatMinAmount sets the minimum amount per order (in exchangeType).
func (s *PlaceBizPushOrderService) SetFiatMinAmount(fiatMinAmount decimal.Decimal) *PlaceBizPushOrderService {
	s.body["fiatMinAmount"] = fiatMinAmount.String()
	return s
}

// SetFiatMaxAmount sets the maximum amount per order (in exchangeType).
func (s *PlaceBizPushOrderService) SetFiatMaxAmount(fiatMaxAmount decimal.Decimal) *PlaceBizPushOrderService {
	s.body["fiatMaxAmount"] = fiatMaxAmount.String()
	return s
}

// SetTierLimit sets the minimum counterparty VIP level (0 means no requirement).
func (s *PlaceBizPushOrderService) SetTierLimit(tierLimit string) *PlaceBizPushOrderService {
	s.body["tierLimit"] = tierLimit
	return s
}

// SetVerifiedLimit sets the minimum counterparty verification level (0 no limit).
func (s *PlaceBizPushOrderService) SetVerifiedLimit(verifiedLimit string) *PlaceBizPushOrderService {
	s.body["verifiedLimit"] = verifiedLimit
	return s
}

// SetRegTimeLimit sets the minimum counterparty account age in days (0 no limit).
func (s *PlaceBizPushOrderService) SetRegTimeLimit(regTimeLimit string) *PlaceBizPushOrderService {
	s.body["regTimeLimit"] = regTimeLimit
	return s
}

// SetAdvertisersLimit restricts trading with the advertiser (0 no, 1 yes).
func (s *PlaceBizPushOrderService) SetAdvertisersLimit(advertisersLimit string) *PlaceBizPushOrderService {
	s.body["advertisersLimit"] = advertisersLimit
	return s
}

// SetPolymarketLimit restricts trading with Polymarket users (0 no, 1 yes).
func (s *PlaceBizPushOrderService) SetPolymarketLimit(polymarketLimit int) *PlaceBizPushOrderService {
	s.body["polymarket_limit"] = polymarketLimit
	return s
}

// SetExpireMin sets the payment timeout in minutes.
func (s *PlaceBizPushOrderService) SetExpireMin(expireMin string) *PlaceBizPushOrderService {
	s.body["expire_min"] = expireMin
	return s
}

// SetTradeTips sets the advertisement trade terms shown to ordering users.
func (s *PlaceBizPushOrderService) SetTradeTips(tradeTips string) *PlaceBizPushOrderService {
	s.body["trade_tips"] = tradeTips
	return s
}

// SetAutoReply sets the auto-reply content sent after an order is created.
func (s *PlaceBizPushOrderService) SetAutoReply(autoReply string) *PlaceBizPushOrderService {
	s.body["auto_reply"] = autoReply
	return s
}

// SetMinCompletedLimit sets the minimum completed orders for a counterparty
// (-1 unlimited).
func (s *PlaceBizPushOrderService) SetMinCompletedLimit(minCompletedLimit string) *PlaceBizPushOrderService {
	s.body["min_completed_limit"] = minCompletedLimit
	return s
}

// SetMaxCompletedLimit sets the maximum completed orders for a counterparty
// (-1 unlimited).
func (s *PlaceBizPushOrderService) SetMaxCompletedLimit(maxCompletedLimit string) *PlaceBizPushOrderService {
	s.body["max_completed_limit"] = maxCompletedLimit
	return s
}

// SetCompletedRateLimit sets the counterparty minimum 30-day completion rate
// (-1 no limit).
func (s *PlaceBizPushOrderService) SetCompletedRateLimit(completedRateLimit string) *PlaceBizPushOrderService {
	s.body["completed_rate_limit"] = completedRateLimit
	return s
}

// SetUserCountryLimit sets the KYC nationality restriction (-1 no restriction).
func (s *PlaceBizPushOrderService) SetUserCountryLimit(userCountryLimit string) *PlaceBizPushOrderService {
	s.body["user_country_limit"] = userCountryLimit
	return s
}

// SetUserOrderLimit sets the maximum concurrent orders for a counterparty
// (-1 unlimited).
func (s *PlaceBizPushOrderService) SetUserOrderLimit(userOrderLimit string) *PlaceBizPushOrderService {
	s.body["user_order_limit"] = userOrderLimit
	return s
}

// SetRateReferenceID sets the floating price reference (1 platform, 2 Gate,
// 3 spot).
func (s *PlaceBizPushOrderService) SetRateReferenceID(rateReferenceID string) *PlaceBizPushOrderService {
	s.body["rateReferenceId"] = rateReferenceID
	return s
}

// SetRateOffset sets the absolute floating offset ratio (e.g. "0.5" for 0.5%).
func (s *PlaceBizPushOrderService) SetRateOffset(rateOffset string) *PlaceBizPushOrderService {
	s.body["rateOffset"] = rateOffset
	return s
}

// SetFloatTrend sets the floating direction: "0" markup, "1" markdown.
func (s *PlaceBizPushOrderService) SetFloatTrend(floatTrend string) *PlaceBizPushOrderService {
	s.body["float_trend"] = floatTrend
	return s
}

// SetTeamPaymentUID sets the team payee UID (optional for non-team merchants).
func (s *PlaceBizPushOrderService) SetTeamPaymentUID(teamPaymentUID string) *PlaceBizPushOrderService {
	s.body["team_payment_uid"] = teamPaymentUID
	return s
}

func (s *PlaceBizPushOrderService) Do(ctx context.Context) (*P2PResponse[P2PPlaceAdResult], error) {
	req := request.Post(ctx, s.c, "/api/v4/p2p/merchant/books/place_biz_push_order", s.body).WithSign()
	return doP2P[P2PPlaceAdResult](req)
}

// P2PPlaceAdResult is the payload of a publish/edit-ad call. It is empty on
// success and carries risk-control details when the ad content is flagged.
type P2PPlaceAdResult struct {
	RiskCode  int          `json:"risk_code"`
	RiskEvent P2PRiskEvent `json:"risk_event"`
}

// P2PRiskEvent describes a risk-control prompt raised for advertisement content.
type P2PRiskEvent struct {
	Type            string          `json:"type"`
	Title           string          `json:"title"`
	Msg             string          `json:"msg"`
	Action          []P2PRiskAction `json:"action"`
	ContentRiskType string          `json:"content_risk_type"`
	TradeTips       string          `json:"trade_tips"`
	AutoReply       string          `json:"auto_reply"`
}

// P2PRiskAction is one available action attached to a risk-control prompt.
type P2PRiskAction struct {
	ActionType string         `json:"action_type"`
	Title      string         `json:"title"`
	Mainly     int            `json:"mainly"`
	ActionData map[string]any `json:"action_data"`
}

// AdsUpdateStatusService -- POST /api/v4/p2p/merchant/books/ads_update_status (private)
//
// Lists, delists or closes one of the merchant's advertisements.
type AdsUpdateStatusService struct {
	c    *P2PClient
	body map[string]any
}

func (c *P2PClient) NewAdsUpdateStatusService(advNo int64, advStatus int) *AdsUpdateStatusService {
	return &AdsUpdateStatusService{c: c, body: map[string]any{
		"adv_no":     advNo,
		"adv_status": advStatus,
	}}
}

func (s *AdsUpdateStatusService) Do(ctx context.Context) (*P2PResponse[P2PAdStatusData], error) {
	req := request.Post(ctx, s.c, "/api/v4/p2p/merchant/books/ads_update_status", s.body).WithSign()
	return doP2P[P2PAdStatusData](req)
}

// P2PAdStatusData reports an advertisement's status after an update (1 listed,
// 3 delisted, 4 closed).
type P2PAdStatusData struct {
	Status int `json:"status"`
}

// AdsDetailService -- POST /api/v4/p2p/merchant/books/ads_detail (private)
//
// Returns the full detail of one of the merchant's advertisements.
type AdsDetailService struct {
	c    *P2PClient
	body map[string]any
}

func (c *P2PClient) NewAdsDetailService(advNo string) *AdsDetailService {
	return &AdsDetailService{c: c, body: map[string]any{
		"adv_no": advNo,
	}}
}

func (s *AdsDetailService) Do(ctx context.Context) (*P2PResponse[P2PAdDetail], error) {
	req := request.Post(ctx, s.c, "/api/v4/p2p/merchant/books/ads_detail", s.body).WithSign()
	return doP2P[P2PAdDetail](req)
}

// P2PAdDetail is the full detail of a single advertisement.
type P2PAdDetail struct {
	Rate               decimal.Decimal `json:"rate"`
	Type               string          `json:"type"`
	Amount             decimal.Decimal `json:"amount"`
	MinAmount          decimal.Decimal `json:"min_amount"`
	MaxAmount          decimal.Decimal `json:"max_amount"`
	FiatMinAmount      decimal.Decimal `json:"fiat_min_amount"`
	FiatMaxAmount      decimal.Decimal `json:"fiat_max_amount"`
	LimitBasis         int             `json:"limit_basis"`
	LimitBasisText     string          `json:"limit_basis_text"`
	Total              decimal.Decimal `json:"total"`
	PayAli             int             `json:"pay_ali"`
	PayBank            int             `json:"pay_bank"`
	PayPaypal          int             `json:"pay_paypal"`
	PayWechat          int             `json:"pay_wechat"`
	PayTypeNum         string          `json:"pay_type_num"`
	PayTypeJSON        string          `json:"pay_type_json"`
	LockedAmount       decimal.Decimal `json:"locked_amount"`
	OrderID            int64           `json:"orderid"`
	Timestamp          time.Time       `json:"timestamp,format:unix"`
	CurrencyType       string          `json:"currency_type"`
	WantType           string          `json:"want_type"`
	HideRate           string          `json:"hide_rate"`
	TradeTips          string          `json:"trade_tips"`
	AutoReply          string          `json:"auto_reply"`
	RateRefID          int             `json:"rate_ref_id"`
	RateOffset         decimal.Decimal `json:"rate_offset"`
	Status             string          `json:"status"`
	RateFixed          int             `json:"rate_fixed"`
	FloatTrend         int             `json:"float_trend"`
	ExpireMin          int             `json:"expire_min"`
	TierLimit          int             `json:"tier_limit"`
	RegTimeLimit       int             `json:"reg_time_limit"`
	AdvertisersLimit   int             `json:"advertisers_limit"`
	PolymarketLimit    int             `json:"polymarket_limit"`
	MinCompletedLimit  int             `json:"min_completed_limit"`
	MaxCompletedLimit  int             `json:"max_completed_limit"`
	UserOrdersLimit    int             `json:"user_orders_limit"`
	CompletedRateLimit decimal.Decimal `json:"completed_rate_limit"`
	LimitCountryCn     string          `json:"limit_country_cn"`
	LimitCountryEn     string          `json:"limit_country_en"`
	IsHedge            int             `json:"is_hedge"`
	HidePayment        int             `json:"hide_payment"`
}

// MyAdsListService -- POST /api/v4/p2p/merchant/books/my_ads_list (private)
//
// Returns the merchant's own advertisement list, optionally filtered by asset,
// fiat and side.
type MyAdsListService struct {
	c    *P2PClient
	body map[string]any
}

func (c *P2PClient) NewMyAdsListService() *MyAdsListService {
	return &MyAdsListService{c: c, body: map[string]any{}}
}

// SetAsset filters the ads to a single crypto asset.
func (s *MyAdsListService) SetAsset(asset string) *MyAdsListService {
	s.body["asset"] = asset
	return s
}

// SetFiatUnit filters the ads to a single fiat currency.
func (s *MyAdsListService) SetFiatUnit(fiatUnit string) *MyAdsListService {
	s.body["fiat_unit"] = fiatUnit
	return s
}

// SetTradeType filters the ads by side ("buy" or "sell").
func (s *MyAdsListService) SetTradeType(tradeType string) *MyAdsListService {
	s.body["trade_type"] = tradeType
	return s
}

func (s *MyAdsListService) Do(ctx context.Context) (*P2PResponse[P2PMyAdsListData], error) {
	req := request.Post(ctx, s.c, "/api/v4/p2p/merchant/books/my_ads_list", s.body).WithSign()
	return doP2P[P2PMyAdsListData](req)
}

// P2PMyAdsListData wraps the merchant's own advertisement list.
type P2PMyAdsListData struct {
	Lists []P2PMyAd `json:"lists"`
}

// P2PMyAd is one of the merchant's own advertisements.
type P2PMyAd struct {
	Type               string          `json:"type"`
	Rate               decimal.Decimal `json:"rate"`
	OriginalRate       decimal.Decimal `json:"original_rate"`
	Amount             decimal.Decimal `json:"amount"`
	Total              decimal.Decimal `json:"total"`
	LimitTotal         string          `json:"limit_total"`
	LimitFiat          string          `json:"limit_fiat"`
	MinAmount          decimal.Decimal `json:"min_amount"`
	MaxAmount          decimal.Decimal `json:"max_amount"`
	PayTypeNum         string          `json:"pay_type_num"`
	PayTypeJSON        string          `json:"pay_type_json"`
	ExpireMin          string          `json:"expire_min"`
	TierLimit          string          `json:"tier_limit"`
	AdvertisersLimit   int             `json:"advertisers_limit"`
	RegTimeLimit       int             `json:"reg_time_limit"`
	VerifiedLimit      int             `json:"verified_limit"`
	MinCompletedLimit  int             `json:"min_completed_limit"`
	MaxCompletedLimit  int             `json:"max_completed_limit"`
	UserCountryLimit   int             `json:"user_country_limit"`
	CompletedRateLimit decimal.Decimal `json:"completed_rate_limit"`
	UserOrdersLimit    int             `json:"user_orders_limit"`
	HidePayment        string          `json:"hide_payment"`
	CurrencyType       string          `json:"currencyType"`
	WantType           string          `json:"want_type"`
	TradeTips          string          `json:"trade_tips"`
	NewHand            int             `json:"new_hand"`
	ID                 string          `json:"id"`
	Status             string          `json:"status"`
	LockedAmount       decimal.Decimal `json:"locked_amount"`
	HideRate           string          `json:"hide_rate"`
	IsOutTime          int             `json:"is_out_time"`
	RateRefID          int             `json:"rate_ref_id"`
	RateOffset         string          `json:"rate_offset"`
	RateFixed          int             `json:"rate_fixed"`
	FloatTrend         int             `json:"float_trend"`
	InDispute          int             `json:"in_dispute"`
	AutoReply          string          `json:"auto_reply"`
	Timestamp          time.Time       `json:"timestamp,format:unix"`
	IsHedge            int             `json:"is_hedge"`
}

// AdsListService -- POST /api/v4/p2p/merchant/books/ads_list (private)
//
// Returns the public advertisement list (other merchants' ads) for a
// crypto/fiat pair and side.
type AdsListService struct {
	c    *P2PClient
	body map[string]any
}

func (c *P2PClient) NewAdsListService(asset, fiatUnit, tradeType string) *AdsListService {
	return &AdsListService{c: c, body: map[string]any{
		"asset":      asset,
		"fiat_unit":  fiatUnit,
		"trade_type": tradeType,
	}}
}

func (s *AdsListService) Do(ctx context.Context) (*P2PResponse[[]P2PAdsListItem], error) {
	req := request.Post(ctx, s.c, "/api/v4/p2p/merchant/books/ads_list", s.body).WithSign()
	return doP2P[[]P2PAdsListItem](req)
}

// P2PAdsListItem is one advertisement in the public ad list.
type P2PAdsListItem struct {
	Index                int                     `json:"index"`
	Asset                string                  `json:"asset"`
	FiatUnit             string                  `json:"fiat_unit"`
	AdvNo                int64                   `json:"adv_no"`
	Price                decimal.Decimal         `json:"price"`
	SurplusAmount        decimal.Decimal         `json:"surplus_amount"`
	MaxSingleTransAmount decimal.Decimal         `json:"max_single_trans_amount"`
	MinSingleTransAmount decimal.Decimal         `json:"min_single_trans_amount"`
	FiatMinAmount        decimal.Decimal         `json:"fiat_min_amount"`
	FiatMaxAmount        decimal.Decimal         `json:"fiat_max_amount"`
	LimitBasis           int                     `json:"limit_basis"`
	LimitBasisText       string                  `json:"limit_basis_text"`
	TradeMethods         []P2PAdsListTradeMethod `json:"trade_methods"`
	NickName             string                  `json:"nick_name"`
}

// P2PAdsListTradeMethod is one supported payment method attached to a public ad.
type P2PAdsListTradeMethod struct {
	IconURLColor    string `json:"icon_url_color"`
	Identifier      string `json:"identifier"`
	PayID           string `json:"pay_id"`
	PayType         string `json:"pay_type"`
	TradeMethodName string `json:"trade_method_name"`
}
