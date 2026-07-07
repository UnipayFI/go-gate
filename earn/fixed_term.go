package earn

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// ListFixedTermProductsService -- GET /api/v4/earn/fixed-term/product (private)
//
// Returns the paginated fixed-term earn product list, optionally filtered by
// currency and product type.
type ListFixedTermProductsService struct {
	c      *EarnClient
	params map[string]string
}

// NewListFixedTermProductsService lists products for the given page and page
// size (limit).
func (c *EarnClient) NewListFixedTermProductsService(page, limit int) *ListFixedTermProductsService {
	return &ListFixedTermProductsService{c: c, params: map[string]string{
		"page":  strconv.Itoa(page),
		"limit": strconv.Itoa(limit),
	}}
}

// SetAsset filters by currency.
func (s *ListFixedTermProductsService) SetAsset(asset string) *ListFixedTermProductsService {
	s.params["asset"] = asset
	return s
}

// SetType filters by product type: 1 for regular, 2 for VIP.
func (s *ListFixedTermProductsService) SetType(productType int) *ListFixedTermProductsService {
	s.params["type"] = strconv.Itoa(productType)
	return s
}

func (s *ListFixedTermProductsService) Do(ctx context.Context) (*FixedTermProductResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/earn/fixed-term/product", s.params).WithSign()
	return request.Do[FixedTermProductResponse](req)
}

// FixedTermProductResponse is the envelope of the fixed-term product list.
type FixedTermProductResponse struct {
	Code      int                  `json:"code"`
	Message   string               `json:"message"`
	Data      FixedTermProductData `json:"data"`
	Timestamp time.Time            `json:"timestamp,format:unix"`
}

// FixedTermProductData carries the fixed-term product list and paging metadata.
type FixedTermProductData struct {
	List  []FixedTermProduct `json:"list"`
	Total int                `json:"total"`
}

// FixedTermProduct is a single fixed-term earn product and its terms.
type FixedTermProduct struct {
	ID                         int             `json:"id"`
	Name                       string          `json:"name"`
	Asset                      string          `json:"asset"`
	LockUpPeriod               int             `json:"lock_up_period"`
	MinLendAmount              decimal.Decimal `json:"min_lend_amount"`
	UserMaxLendAmount          decimal.Decimal `json:"user_max_lend_amount"`
	UserMaxLendVolume          decimal.Decimal `json:"user_max_lend_volume"`
	TotalLendAmount            decimal.Decimal `json:"total_lend_amount"`
	TotalLendVolume            decimal.Decimal `json:"total_lend_volume"`
	TotalInterest              decimal.Decimal `json:"total_interest"`
	TotalInterestVolume        decimal.Decimal `json:"total_interest_volume"`
	YearRate                   decimal.Decimal `json:"year_rate"`
	Type                       int             `json:"type"`
	ShowPage                   string          `json:"show_page"`
	PreRedeem                  int             `json:"pre_redeem"`
	Reinvest                   int             `json:"reinvest"`
	RedeemAccount              int             `json:"redeem_account"`
	MinVIP                     int             `json:"min_vip"`
	MaxVIP                     int             `json:"max_vip"`
	Status                     int             `json:"status"`
	ShowStatus                 int             `json:"show_status"`
	UserWhiteListTurnOn        int             `json:"user_white_list_turn_on"`
	UseSystemRate              int             `json:"use_system_rate"`
	IsRevolvingLimit           int             `json:"is_revolving_limit"`
	IsForAllUser               int             `json:"is_for_all_user"`
	ForUserCrowd               string          `json:"for_user_crowd"`
	CheckUserCrowdURL          string          `json:"check_user_crowd_url"`
	TagInfo                    string          `json:"tag_info"`
	CouponsInfo                string          `json:"coupons_info"`
	ExtraInfo                  string          `json:"extra_info"`
	SaleStartAt                time.Time       `json:"sale_start_at,format:unix"`
	CreateTime                 string          `json:"create_time"`
	UpdateTime                 string          `json:"update_time"`
	CreateUser                 string          `json:"create_user"`
	UpdateUser                 string          `json:"update_user"`
	AllowedAccount             string          `json:"allowed_account"`
	Title                      string          `json:"title"`
	Subtitle                   string          `json:"subtitle"`
	DialogTitle                string          `json:"dialog_title"`
	DialogContent              string          `json:"dialog_content"`
	DialogURLTitle             string          `json:"dialog_url_title"`
	CouponInfo                 any             `json:"coupon_info"`
	BonusInfo                  any             `json:"bonus_info"`
	BonusBoostInfo             any             `json:"bonus_boost_info"`
	UserTotalAmount            decimal.Decimal `json:"user_total_amount"`
	ProductTotalVolume         decimal.Decimal `json:"product_total_volume"`
	ProductTotalInterestVolume decimal.Decimal `json:"product_total_interest_volume"`
	SaleStatus                 int             `json:"sale_status"`
	Price                      decimal.Decimal `json:"price"`
	ProductCouponSendCount     int             `json:"product_coupon_send_count"`
}

// ListFixedTermProductsByAssetService -- GET /api/v4/earn/fixed-term/product/{asset}/list (private)
//
// Returns the fixed-term earn product list for a single currency.
type ListFixedTermProductsByAssetService struct {
	c      *EarnClient
	asset  string
	params map[string]string
}

func (c *EarnClient) NewListFixedTermProductsByAssetService(asset string) *ListFixedTermProductsByAssetService {
	return &ListFixedTermProductsByAssetService{c: c, asset: asset, params: map[string]string{}}
}

// SetType filters by product type: "" or 1 for regular, 2 for VIP, 0 for all.
func (s *ListFixedTermProductsByAssetService) SetType(productType string) *ListFixedTermProductsByAssetService {
	s.params["type"] = productType
	return s
}

func (s *ListFixedTermProductsByAssetService) Do(ctx context.Context) (*FixedTermProductSimpleResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/earn/fixed-term/product/"+s.asset+"/list", s.params).WithSign()
	return request.Do[FixedTermProductSimpleResponse](req)
}

// FixedTermProductSimpleResponse is the envelope of the single-currency
// fixed-term product list.
type FixedTermProductSimpleResponse struct {
	Code      int                        `json:"code"`
	Message   string                     `json:"message"`
	Data      FixedTermProductSimpleData `json:"data"`
	Timestamp time.Time                  `json:"timestamp,format:unix"`
}

// FixedTermProductSimpleData carries the compact fixed-term product list.
type FixedTermProductSimpleData struct {
	List []FixedTermProductSimple `json:"list"`
}

// FixedTermProductSimple is a compact fixed-term earn product.
type FixedTermProductSimple struct {
	ID              int             `json:"id"`
	Asset           string          `json:"asset"`
	LockUpPeriod    int             `json:"lock_up_period"`
	YearRate        decimal.Decimal `json:"year_rate"`
	Type            int             `json:"type"`
	PreRedeem       int             `json:"pre_redeem"`
	Reinvest        int             `json:"reinvest"`
	RedeemAccount   int             `json:"redeem_account"`
	MinVIP          int             `json:"min_vip"`
	MaxVIP          int             `json:"max_vip"`
	SaleStatus      int             `json:"sale_status"`
	UserTotalAmount decimal.Decimal `json:"user_total_amount"`
	SaleStartAt     time.Time       `json:"sale_start_at,format:unix"`
	Title           string          `json:"title"`
	Subtitle        string          `json:"subtitle"`
	DialogTitle     string          `json:"dialog_title"`
	DialogContent   string          `json:"dialog_content"`
	DialogURLTitle  string          `json:"dialog_url_title"`
	BonusInfo       any             `json:"bonus_info"`
	BonusBoostInfo  any             `json:"bonus_boost_info"`
}

// ListFixedTermLendsService -- GET /api/v4/earn/fixed-term/user/lend (private)
//
// Returns the authenticated user's fixed-term subscription orders, paginated.
type ListFixedTermLendsService struct {
	c      *EarnClient
	params map[string]string
}

// NewListFixedTermLendsService lists subscription orders. orderType is "1" for
// current orders or "2" for historical orders; page and limit paginate.
func (c *EarnClient) NewListFixedTermLendsService(orderType string, page, limit int) *ListFixedTermLendsService {
	return &ListFixedTermLendsService{c: c, params: map[string]string{
		"order_type": orderType,
		"page":       strconv.Itoa(page),
		"limit":      strconv.Itoa(limit),
	}}
}

// SetProductID filters by product id.
func (s *ListFixedTermLendsService) SetProductID(productID int) *ListFixedTermLendsService {
	s.params["product_id"] = strconv.Itoa(productID)
	return s
}

// SetOrderID filters by order id.
func (s *ListFixedTermLendsService) SetOrderID(orderID int64) *ListFixedTermLendsService {
	s.params["order_id"] = strconv.FormatInt(orderID, 10)
	return s
}

// SetAsset filters by currency.
func (s *ListFixedTermLendsService) SetAsset(asset string) *ListFixedTermLendsService {
	s.params["asset"] = asset
	return s
}

// SetSubBusiness filters by sub-business.
func (s *ListFixedTermLendsService) SetSubBusiness(subBusiness int) *ListFixedTermLendsService {
	s.params["sub_business"] = strconv.Itoa(subBusiness)
	return s
}

// SetBusinessFilter sets the business filter (JSON array string, e.g.
// [{"business":1,"sub_business":0}]).
func (s *ListFixedTermLendsService) SetBusinessFilter(businessFilter string) *ListFixedTermLendsService {
	s.params["business_filter"] = businessFilter
	return s
}

func (s *ListFixedTermLendsService) Do(ctx context.Context) (*FixedTermLendResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/earn/fixed-term/user/lend", s.params).WithSign()
	return request.Do[FixedTermLendResponse](req)
}

// FixedTermLendResponse is the envelope of the fixed-term subscription order
// list.
type FixedTermLendResponse struct {
	Code      int               `json:"code"`
	Message   string            `json:"message"`
	Data      FixedTermLendData `json:"data"`
	Timestamp time.Time         `json:"timestamp,format:unix"`
}

// FixedTermLendData carries the fixed-term subscription order list and paging
// metadata.
type FixedTermLendData struct {
	List  []FixedTermLendOrder `json:"list"`
	Total int                  `json:"total"`
}

// FixedTermLendOrder is a single fixed-term earn subscription order.
type FixedTermLendOrder struct {
	ID                   int                  `json:"id"`
	Business             int                  `json:"business"`
	OrderID              int64                `json:"order_id"`
	UserID               int64                `json:"user_id"`
	Asset                string               `json:"asset"`
	ProductID            int                  `json:"product_id"`
	LockUpPeriod         int                  `json:"lock_up_period"`
	Principal            decimal.Decimal      `json:"principal"`
	YearRate             decimal.Decimal      `json:"year_rate"`
	ProductType          int                  `json:"product_type"`
	Interest             decimal.Decimal      `json:"interest"`
	Status               int                  `json:"status"`
	ReinvestStatus       int                  `json:"reinvest_status"`
	RedeemAccountType    int                  `json:"redeem_account_type"`
	OriginOrder          string               `json:"origin_order"`
	RedeemType           int                  `json:"redeem_type"`
	RedeemTime           string               `json:"redeem_time"`
	FinishTime           string               `json:"finish_time"`
	CreateTime           string               `json:"create_time"`
	YearRatePercent      decimal.Decimal      `json:"year_rate_perent"`
	TotalYearRatePercent decimal.Decimal      `json:"total_year_rate_percent"`
	TotalInterest        decimal.Decimal      `json:"total_interest"`
	ProductInfo          FixedTermProductInfo `json:"product_info"`
	BonusInfo            FixedTermBonusInfo   `json:"bonus_info"`
	CouponInfo           FixedTermCouponInfo  `json:"coupon_info"`
	RedeemAt             time.Time            `json:"redeem_at,format:unix"`
	FinishAt             time.Time            `json:"finish_at,format:unix"`
	CreateAt             time.Time            `json:"create_at,format:unix"`
	Icon                 string               `json:"icon"`
}

// FixedTermProductInfo is the product configuration carried on a subscription
// order.
type FixedTermProductInfo struct {
	PreRedeem     int `json:"pre_redeem"`
	Reinvest      int `json:"reinvest"`
	RedeemAccount int `json:"redeem_account"`
	MinVIP        int `json:"min_vip"`
	MaxVIP        int `json:"max_vip"`
}

// FixedTermBonusInfo is the bonus reward campaign attached to a subscription
// order.
type FixedTermBonusInfo struct {
	ID                    int                  `json:"id"`
	ProductID             int                  `json:"product_id"`
	Asset                 string               `json:"asset"`
	BonusAsset            string               `json:"bonus_asset"`
	KYCLimit              string               `json:"kyc_limit"`
	LadderAPR             []FixedTermLadderAPR `json:"ladder_apr"`
	TotalBonusAmount      decimal.Decimal      `json:"total_bonus_amount"`
	UserTotalBonusAmount  decimal.Decimal      `json:"user_total_bonus_amount"`
	Status                int                  `json:"status"`
	StartTime             string               `json:"start_time"`
	EndTime               string               `json:"end_time"`
	CreateTime            string               `json:"create_time"`
	StartAt               time.Time            `json:"start_at,format:unix"`
	EndAt                 time.Time            `json:"end_at,format:unix"`
	TotalIssuedAmount     decimal.Decimal      `json:"total_issued_amount"`
	UserTotalIssuedAmount decimal.Decimal      `json:"user_total_issued_amount"`
	BonusAssetPrice       decimal.Decimal      `json:"bonus_asset_price"`
	ProductAssetPrice     decimal.Decimal      `json:"product_asset_price"`
	ProductYearRate       decimal.Decimal      `json:"product_year_rate"`
}

// FixedTermLadderAPR is one tier of a bonus campaign's tiered annual rate.
type FixedTermLadderAPR struct {
	APR   decimal.Decimal `json:"apr"`
	Left  decimal.Decimal `json:"left"`
	Right decimal.Decimal `json:"right"`
}

// FixedTermCouponInfo is the interest-rate-boost coupon attached to a
// subscription order.
type FixedTermCouponInfo struct {
	ID              int             `json:"id"`
	Business        int             `json:"business"`
	UserID          int64           `json:"user_id"`
	Asset           string          `json:"asset"`
	OrderID         int64           `json:"order_id"`
	FinancialRateID int             `json:"financial_rate_id"`
	BuyLimitLow     decimal.Decimal `json:"buy_limit_low"`
	BuyLimitHigh    decimal.Decimal `json:"buy_limit_high"`
	RateDay         int             `json:"rate_day"`
	RateRatio       decimal.Decimal `json:"rate_ratio"`
	CouponDays      int             `json:"coupon_days"`
	CouponPrincipal decimal.Decimal `json:"coupon_principal"`
	CouponYearRate  decimal.Decimal `json:"coupon_year_rate"`
	CouponInterest  decimal.Decimal `json:"coupon_interest"`
	Status          int             `json:"status"`
	FinishTime      string          `json:"finish_time"`
	CreateTime      string          `json:"create_time"`
}

// SubscribeFixedTermService -- POST /api/v4/earn/fixed-term/user/lend (private)
//
// Subscribes to a fixed-term earn product.
type SubscribeFixedTermService struct {
	c    *EarnClient
	body map[string]any
}

func (c *EarnClient) NewSubscribeFixedTermService(productID int, amount decimal.Decimal) *SubscribeFixedTermService {
	return &SubscribeFixedTermService{c: c, body: map[string]any{
		"product_id": productID,
		"amount":     amount.String(),
	}}
}

// SetYearRate sets the expected annual interest rate.
func (s *SubscribeFixedTermService) SetYearRate(yearRate decimal.Decimal) *SubscribeFixedTermService {
	s.body["year_rate"] = yearRate.String()
	return s
}

// SetReinvestStatus toggles auto-renewal: 0 for disabled, 1 for enabled.
func (s *SubscribeFixedTermService) SetReinvestStatus(reinvestStatus int) *SubscribeFixedTermService {
	s.body["reinvest_status"] = reinvestStatus
	return s
}

// SetRedeemAccountType sets the redemption payout account type: 1 for spot.
func (s *SubscribeFixedTermService) SetRedeemAccountType(redeemAccountType int) *SubscribeFixedTermService {
	s.body["redeem_account_type"] = redeemAccountType
	return s
}

// SetFinancialRateID sets the interest-rate-boost coupon id (0 means not used).
func (s *SubscribeFixedTermService) SetFinancialRateID(financialRateID int) *SubscribeFixedTermService {
	s.body["financial_rate_id"] = financialRateID
	return s
}

// SetSubBusiness sets the sub-business type.
func (s *SubscribeFixedTermService) SetSubBusiness(subBusiness int) *SubscribeFixedTermService {
	s.body["sub_business"] = subBusiness
	return s
}

func (s *SubscribeFixedTermService) Do(ctx context.Context) (*FixedTermLendResult, error) {
	req := request.Post(ctx, s.c, "/api/v4/earn/fixed-term/user/lend", s.body).WithSign()
	return request.Do[FixedTermLendResult](req)
}

// FixedTermLendResult is the envelope of a fixed-term subscription result.
type FixedTermLendResult struct {
	Code      int                     `json:"code"`
	Message   string                  `json:"message"`
	Data      FixedTermLendResultData `json:"data"`
	Timestamp time.Time               `json:"timestamp,format:unix"`
}

// FixedTermLendResultData carries the created subscription order id.
type FixedTermLendResultData struct {
	OrderID int64 `json:"order_id"`
}

// PreRedeemFixedTermService -- POST /api/v4/earn/fixed-term/user/pre-redeem (private)
//
// Redeems a fixed-term earn subscription order.
type PreRedeemFixedTermService struct {
	c    *EarnClient
	body map[string]any
}

func (c *EarnClient) NewPreRedeemFixedTermService(orderID string) *PreRedeemFixedTermService {
	return &PreRedeemFixedTermService{c: c, body: map[string]any{
		"order_id": orderID,
	}}
}

func (s *PreRedeemFixedTermService) Do(ctx context.Context) (*FixedTermPreRedeemResult, error) {
	req := request.Post(ctx, s.c, "/api/v4/earn/fixed-term/user/pre-redeem", s.body).WithSign()
	return request.Do[FixedTermPreRedeemResult](req)
}

// FixedTermPreRedeemResult is the envelope of a fixed-term redemption result;
// data is an empty object on success.
type FixedTermPreRedeemResult struct {
	Code      int                    `json:"code"`
	Message   string                 `json:"message"`
	Data      FixedTermPreRedeemData `json:"data"`
	Timestamp time.Time              `json:"timestamp,format:unix"`
}

// FixedTermPreRedeemData is the (empty) data object of a redemption result.
type FixedTermPreRedeemData struct{}

// ListFixedTermHistoryService -- GET /api/v4/earn/fixed-term/user/history (private)
//
// Returns the authenticated user's fixed-term subscription history, paginated.
type ListFixedTermHistoryService struct {
	c      *EarnClient
	params map[string]string
}

// NewListFixedTermHistoryService lists history records. historyType is "1"
// (subscription), "2" (redemption), "3" (interest) or "4" (bonus reward); page
// and limit paginate.
func (c *EarnClient) NewListFixedTermHistoryService(historyType string, page, limit int) *ListFixedTermHistoryService {
	return &ListFixedTermHistoryService{c: c, params: map[string]string{
		"type":  historyType,
		"page":  strconv.Itoa(page),
		"limit": strconv.Itoa(limit),
	}}
}

// SetProductID filters by product id.
func (s *ListFixedTermHistoryService) SetProductID(productID int) *ListFixedTermHistoryService {
	s.params["product_id"] = strconv.Itoa(productID)
	return s
}

// SetOrderID filters by order id.
func (s *ListFixedTermHistoryService) SetOrderID(orderID string) *ListFixedTermHistoryService {
	s.params["order_id"] = orderID
	return s
}

// SetAsset filters by currency.
func (s *ListFixedTermHistoryService) SetAsset(asset string) *ListFixedTermHistoryService {
	s.params["asset"] = asset
	return s
}

// SetStartAt bounds the result to records at or after this time.
func (s *ListFixedTermHistoryService) SetStartAt(startAt time.Time) *ListFixedTermHistoryService {
	s.params["start_at"] = strconv.FormatInt(startAt.Unix(), 10)
	return s
}

// SetEndAt bounds the result to records at or before this time.
func (s *ListFixedTermHistoryService) SetEndAt(endAt time.Time) *ListFixedTermHistoryService {
	s.params["end_at"] = strconv.FormatInt(endAt.Unix(), 10)
	return s
}

// SetSubBusiness filters by sub-business.
func (s *ListFixedTermHistoryService) SetSubBusiness(subBusiness int) *ListFixedTermHistoryService {
	s.params["sub_business"] = strconv.Itoa(subBusiness)
	return s
}

// SetBusinessFilter sets the business filter (JSON array string, e.g.
// [{"business":1,"sub_business":0}]).
func (s *ListFixedTermHistoryService) SetBusinessFilter(businessFilter string) *ListFixedTermHistoryService {
	s.params["business_filter"] = businessFilter
	return s
}

func (s *ListFixedTermHistoryService) Do(ctx context.Context) (*FixedTermHistoryResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/earn/fixed-term/user/history", s.params).WithSign()
	return request.Do[FixedTermHistoryResponse](req)
}

// FixedTermHistoryResponse is the envelope of the fixed-term history list.
// timestamp is a top-level Unix-second timestamp.
type FixedTermHistoryResponse struct {
	Code      int                  `json:"code"`
	Message   string               `json:"message"`
	Data      FixedTermHistoryData `json:"data"`
	Timestamp time.Time            `json:"timestamp,format:unix"`
}

// FixedTermHistoryData carries the fixed-term history list and paging metadata.
type FixedTermHistoryData struct {
	List  []FixedTermHistoryRecord `json:"list"`
	Total int                      `json:"total"`
}

// FixedTermHistoryRecord is a single fixed-term earn history record.
type FixedTermHistoryRecord struct {
	ID             int             `json:"id"`
	OrderID        int64           `json:"order_id"`
	UserID         int64           `json:"user_id"`
	Asset          string          `json:"asset"`
	UniqTime       string          `json:"uniq_time"`
	BonusID        int             `json:"bonus_id"`
	ProductID      int             `json:"product_id"`
	BonusAsset     string          `json:"bonus_asset"`
	TotalPrincipal decimal.Decimal `json:"total_principal"`
	Amount         decimal.Decimal `json:"amount"`
	AssetPrice     decimal.Decimal `json:"asset_price"`
	Status         int             `json:"status"`
	Detail         string          `json:"detail"`
	CreateTime     string          `json:"create_time"`
	CreateAt       time.Time       `json:"create_at,format:unix"`
	LockUpPeriod   int             `json:"lock_up_period"`
}
