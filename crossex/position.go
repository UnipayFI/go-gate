package crossex

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// GetPositionLeverageService -- GET /api/v4/crossex/positions/leverage (private)
//
// Returns the contract leverage multiplier per trading pair, keyed by symbol.
type GetPositionLeverageService struct {
	c      *CrossexClient
	params map[string]string
}

func (c *CrossexClient) NewGetPositionLeverageService() *GetPositionLeverageService {
	return &GetPositionLeverageService{c: c, params: map[string]string{}}
}

// SetSymbols narrows the result to a comma-separated list of trading pairs.
func (s *GetPositionLeverageService) SetSymbols(symbols string) *GetPositionLeverageService {
	s.params["symbols"] = symbols
	return s
}

func (s *GetPositionLeverageService) Do(ctx context.Context) (map[string]decimal.Decimal, error) {
	req := request.Get(ctx, s.c, "/api/v4/crossex/positions/leverage", s.params).WithSign()
	resp, err := request.Do[map[string]decimal.Decimal](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// UpdatePositionLeverageService -- POST /api/v4/crossex/positions/leverage (private)
//
// Modifies the contract leverage multiplier for a trading pair.
type UpdatePositionLeverageService struct {
	c    *CrossexClient
	body map[string]any
}

func (c *CrossexClient) NewUpdatePositionLeverageService(symbol string, leverage decimal.Decimal) *UpdatePositionLeverageService {
	return &UpdatePositionLeverageService{c: c, body: map[string]any{
		"symbol":   symbol,
		"leverage": leverage.String(),
	}}
}

func (s *UpdatePositionLeverageService) Do(ctx context.Context) (*CrossexLeverage, error) {
	req := request.Post(ctx, s.c, "/api/v4/crossex/positions/leverage", s.body).WithSign()
	return request.Do[CrossexLeverage](req)
}

// CrossexLeverage echoes the leverage change that was requested for a pair.
type CrossexLeverage struct {
	Symbol   string          `json:"symbol"`
	Leverage decimal.Decimal `json:"leverage"`
}

// GetMarginPositionLeverageService -- GET /api/v4/crossex/margin_positions/leverage (private)
//
// Returns the margin (leveraged) trading-pair leverage multiplier, keyed by symbol.
type GetMarginPositionLeverageService struct {
	c      *CrossexClient
	params map[string]string
}

func (c *CrossexClient) NewGetMarginPositionLeverageService() *GetMarginPositionLeverageService {
	return &GetMarginPositionLeverageService{c: c, params: map[string]string{}}
}

// SetSymbols narrows the result to a comma-separated list of trading pairs.
func (s *GetMarginPositionLeverageService) SetSymbols(symbols string) *GetMarginPositionLeverageService {
	s.params["symbols"] = symbols
	return s
}

func (s *GetMarginPositionLeverageService) Do(ctx context.Context) (map[string]decimal.Decimal, error) {
	req := request.Get(ctx, s.c, "/api/v4/crossex/margin_positions/leverage", s.params).WithSign()
	resp, err := request.Do[map[string]decimal.Decimal](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// UpdateMarginPositionLeverageService -- POST /api/v4/crossex/margin_positions/leverage (private)
//
// Modifies the margin (leveraged) trading-pair leverage multiplier.
type UpdateMarginPositionLeverageService struct {
	c    *CrossexClient
	body map[string]any
}

func (c *CrossexClient) NewUpdateMarginPositionLeverageService(symbol string, leverage decimal.Decimal) *UpdateMarginPositionLeverageService {
	return &UpdateMarginPositionLeverageService{c: c, body: map[string]any{
		"symbol":   symbol,
		"leverage": leverage.String(),
	}}
}

func (s *UpdateMarginPositionLeverageService) Do(ctx context.Context) (*CrossexLeverage, error) {
	req := request.Post(ctx, s.c, "/api/v4/crossex/margin_positions/leverage", s.body).WithSign()
	return request.Do[CrossexLeverage](req)
}

// ClosePositionService -- POST /api/v4/crossex/position (private)
//
// Fully closes a contract or leveraged position for a trading pair.
type ClosePositionService struct {
	c    *CrossexClient
	body map[string]any
}

func (c *CrossexClient) NewClosePositionService(symbol string) *ClosePositionService {
	return &ClosePositionService{c: c, body: map[string]any{
		"symbol": symbol,
	}}
}

// SetPositionSide sets the position side (required for leveraged positions).
func (s *ClosePositionService) SetPositionSide(positionSide string) *ClosePositionService {
	s.body["position_side"] = positionSide
	return s
}

func (s *ClosePositionService) Do(ctx context.Context) (*CrossexOrderResult, error) {
	req := request.Post(ctx, s.c, "/api/v4/crossex/position", s.body).WithSign()
	return request.Do[CrossexOrderResult](req)
}

// ListPositionsService -- GET /api/v4/crossex/positions (private)
//
// Returns the account's open contract positions.
type ListPositionsService struct {
	c      *CrossexClient
	params map[string]string
}

func (c *CrossexClient) NewListPositionsService() *ListPositionsService {
	return &ListPositionsService{c: c, params: map[string]string{}}
}

// SetSymbol narrows the result to a single trading pair.
func (s *ListPositionsService) SetSymbol(symbol string) *ListPositionsService {
	s.params["symbol"] = symbol
	return s
}

// SetExchangeType narrows the result to a single venue.
func (s *ListPositionsService) SetExchangeType(exchangeType string) *ListPositionsService {
	s.params["exchange_type"] = exchangeType
	return s
}

func (s *ListPositionsService) Do(ctx context.Context) ([]CrossexPosition, error) {
	req := request.Get(ctx, s.c, "/api/v4/crossex/positions", s.params).WithSign()
	resp, err := request.Do[[]CrossexPosition](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CrossexPosition is one open contract position. funding_time, create_time and
// update_time are millisecond Unix timestamps (funding_time 0 means the funding
// fee has not been collected yet).
type CrossexPosition struct {
	UserID            string          `json:"user_id"`
	PositionID        string          `json:"position_id"`
	Symbol            string          `json:"symbol"`
	PositionSide      string          `json:"position_side"`
	InitialMargin     decimal.Decimal `json:"initial_margin"`
	MaintenanceMargin decimal.Decimal `json:"maintenance_margin"`
	PositionQty       decimal.Decimal `json:"position_qty"`
	PositionValue     decimal.Decimal `json:"position_value"`
	UPnL              decimal.Decimal `json:"upnl"`
	UPnLRate          decimal.Decimal `json:"upnl_rate"`
	EntryPrice        decimal.Decimal `json:"entry_price"`
	MarkPrice         decimal.Decimal `json:"mark_price"`
	Leverage          decimal.Decimal `json:"leverage"`
	MaxLeverage       decimal.Decimal `json:"max_leverage"`
	RiskLimit         decimal.Decimal `json:"risk_limit"`
	Fee               decimal.Decimal `json:"fee"`
	FundingFee        decimal.Decimal `json:"funding_fee"`
	FundingTime       time.Time       `json:"funding_time,string,format:unixmilli"`
	CreateTime        time.Time       `json:"create_time,string,format:unixmilli"`
	UpdateTime        time.Time       `json:"update_time,string,format:unixmilli"`
	ClosedPnL         decimal.Decimal `json:"closed_pnl"`
}

// ListMarginPositionsService -- GET /api/v4/crossex/margin_positions (private)
//
// Returns the account's open leveraged (margin) positions.
type ListMarginPositionsService struct {
	c      *CrossexClient
	params map[string]string
}

func (c *CrossexClient) NewListMarginPositionsService() *ListMarginPositionsService {
	return &ListMarginPositionsService{c: c, params: map[string]string{}}
}

// SetSymbol narrows the result to a single trading pair.
func (s *ListMarginPositionsService) SetSymbol(symbol string) *ListMarginPositionsService {
	s.params["symbol"] = symbol
	return s
}

// SetExchangeType narrows the result to a single venue.
func (s *ListMarginPositionsService) SetExchangeType(exchangeType string) *ListMarginPositionsService {
	s.params["exchange_type"] = exchangeType
	return s
}

func (s *ListMarginPositionsService) Do(ctx context.Context) ([]CrossexMarginPosition, error) {
	req := request.Get(ctx, s.c, "/api/v4/crossex/margin_positions", s.params).WithSign()
	resp, err := request.Do[[]CrossexMarginPosition](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CrossexMarginPosition is one open leveraged (margin) position. create_time and
// update_time are millisecond Unix timestamps.
type CrossexMarginPosition struct {
	UserID            string          `json:"user_id"`
	PositionID        string          `json:"position_id"`
	Symbol            string          `json:"symbol"`
	PositionSide      string          `json:"position_side"`
	InitialMargin     decimal.Decimal `json:"initial_margin"`
	MaintenanceMargin decimal.Decimal `json:"maintenance_margin"`
	AssetQty          decimal.Decimal `json:"asset_qty"`
	AssetCoin         string          `json:"asset_coin"`
	PositionValue     decimal.Decimal `json:"position_value"`
	Liability         decimal.Decimal `json:"liability"`
	LiabilityCoin     string          `json:"liability_coin"`
	Interest          decimal.Decimal `json:"interest"`
	MaxPositionQty    decimal.Decimal `json:"max_position_qty"`
	EntryPrice        decimal.Decimal `json:"entry_price"`
	IndexPrice        decimal.Decimal `json:"index_price"`
	UPnL              decimal.Decimal `json:"upnl"`
	UPnLRate          decimal.Decimal `json:"upnl_rate"`
	Leverage          decimal.Decimal `json:"leverage"`
	MaxLeverage       decimal.Decimal `json:"max_leverage"`
	CreateTime        time.Time       `json:"create_time,string,format:unixmilli"`
	UpdateTime        time.Time       `json:"update_time,string,format:unixmilli"`
}

// GetADLRankService -- GET /api/v4/crossex/adl_rank (private)
//
// Returns the account's ADL (auto-deleveraging) reduction ranking for a pair.
type GetADLRankService struct {
	c      *CrossexClient
	params map[string]string
}

func (c *CrossexClient) NewGetADLRankService(symbol string) *GetADLRankService {
	return &GetADLRankService{c: c, params: map[string]string{
		"symbol": symbol,
	}}
}

func (s *GetADLRankService) Do(ctx context.Context) (*CrossexADLRank, error) {
	req := request.Get(ctx, s.c, "/api/v4/crossex/adl_rank", s.params).WithSign()
	return request.Do[CrossexADLRank](req)
}

// CrossexADLRank is the ADL reduction ranking for a trading pair.
type CrossexADLRank struct {
	UserID          string `json:"user_id"`
	Symbol          string `json:"symbol"`
	CrossexADLRank  string `json:"crossex_adl_rank"`
	ExchangeADLRank string `json:"exchange_adl_rank"`
}

// ListHistoryPositionsService -- GET /api/v4/crossex/history_positions (private)
//
// Returns the account's closed contract-position history.
type ListHistoryPositionsService struct {
	c      *CrossexClient
	params map[string]string
}

func (c *CrossexClient) NewListHistoryPositionsService() *ListHistoryPositionsService {
	return &ListHistoryPositionsService{c: c, params: map[string]string{}}
}

// SetPage selects the result page.
func (s *ListHistoryPositionsService) SetPage(page int) *ListHistoryPositionsService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit caps the number of records returned in a single list (max 1000).
func (s *ListHistoryPositionsService) SetLimit(limit int) *ListHistoryPositionsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetSymbol narrows the result to a single trading pair.
func (s *ListHistoryPositionsService) SetSymbol(symbol string) *ListHistoryPositionsService {
	s.params["symbol"] = symbol
	return s
}

// SetFrom sets the start time (millisecond precision).
func (s *ListHistoryPositionsService) SetFrom(from time.Time) *ListHistoryPositionsService {
	s.params["from"] = strconv.FormatInt(from.UnixMilli(), 10)
	return s
}

// SetTo sets the end time (millisecond precision).
func (s *ListHistoryPositionsService) SetTo(to time.Time) *ListHistoryPositionsService {
	s.params["to"] = strconv.FormatInt(to.UnixMilli(), 10)
	return s
}

func (s *ListHistoryPositionsService) Do(ctx context.Context) ([]CrossexHistoricalPosition, error) {
	req := request.Get(ctx, s.c, "/api/v4/crossex/history_positions", s.params).WithSign()
	resp, err := request.Do[[]CrossexHistoricalPosition](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CrossexHistoricalPosition is one closed contract-position record. create_time
// and update_time are millisecond Unix timestamps.
type CrossexHistoricalPosition struct {
	PositionID     string          `json:"position_id"`
	UserID         string          `json:"user_id"`
	Symbol         string          `json:"symbol"`
	ClosedType     string          `json:"closed_type"`
	ClosedPnL      decimal.Decimal `json:"closed_pnl"`
	ClosedPnLRate  decimal.Decimal `json:"closed_pnl_rate"`
	OpenAvgPrice   decimal.Decimal `json:"open_avg_price"`
	ClosedAvgPrice decimal.Decimal `json:"closed_avg_price"`
	MaxPositionQty decimal.Decimal `json:"max_position_qty"`
	ClosedQty      decimal.Decimal `json:"closed_qty"`
	ClosedValue    decimal.Decimal `json:"closed_value"`
	Fee            decimal.Decimal `json:"fee"`
	LiqFee         decimal.Decimal `json:"liq_fee"`
	FundingFee     decimal.Decimal `json:"funding_fee"`
	PositionSide   string          `json:"position_side"`
	PositionMode   string          `json:"position_mode"`
	Leverage       decimal.Decimal `json:"leverage"`
	BusinessType   string          `json:"business_type"`
	CreateTime     time.Time       `json:"create_time,string,format:unixmilli"`
	UpdateTime     time.Time       `json:"update_time,string,format:unixmilli"`
}

// ListHistoryMarginPositionsService -- GET /api/v4/crossex/history_margin_positions (private)
//
// Returns the account's closed leveraged (margin) position history.
type ListHistoryMarginPositionsService struct {
	c      *CrossexClient
	params map[string]string
}

func (c *CrossexClient) NewListHistoryMarginPositionsService() *ListHistoryMarginPositionsService {
	return &ListHistoryMarginPositionsService{c: c, params: map[string]string{}}
}

// SetPage selects the result page.
func (s *ListHistoryMarginPositionsService) SetPage(page int) *ListHistoryMarginPositionsService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit caps the number of records returned in a single list (max 1000).
func (s *ListHistoryMarginPositionsService) SetLimit(limit int) *ListHistoryMarginPositionsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetSymbol narrows the result to a single trading pair.
func (s *ListHistoryMarginPositionsService) SetSymbol(symbol string) *ListHistoryMarginPositionsService {
	s.params["symbol"] = symbol
	return s
}

// SetFrom sets the start time (millisecond precision).
func (s *ListHistoryMarginPositionsService) SetFrom(from time.Time) *ListHistoryMarginPositionsService {
	s.params["from"] = strconv.FormatInt(from.UnixMilli(), 10)
	return s
}

// SetTo sets the end time (millisecond precision).
func (s *ListHistoryMarginPositionsService) SetTo(to time.Time) *ListHistoryMarginPositionsService {
	s.params["to"] = strconv.FormatInt(to.UnixMilli(), 10)
	return s
}

func (s *ListHistoryMarginPositionsService) Do(ctx context.Context) ([]CrossexHistoricalMarginPosition, error) {
	req := request.Get(ctx, s.c, "/api/v4/crossex/history_margin_positions", s.params).WithSign()
	resp, err := request.Do[[]CrossexHistoricalMarginPosition](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CrossexHistoricalMarginPosition is one closed leveraged-position record.
// create_time and update_time are millisecond Unix timestamps.
type CrossexHistoricalMarginPosition struct {
	PositionID     string          `json:"position_id"`
	UserID         string          `json:"user_id"`
	Symbol         string          `json:"symbol"`
	ClosedType     string          `json:"closed_type"`
	ClosedPnL      decimal.Decimal `json:"closed_pnl"`
	ClosedPnLRate  decimal.Decimal `json:"closed_pnl_rate"`
	OpenAvgPrice   decimal.Decimal `json:"open_avg_price"`
	ClosedAvgPrice decimal.Decimal `json:"closed_avg_price"`
	MaxPositionQty decimal.Decimal `json:"max_position_qty"`
	ClosedQty      decimal.Decimal `json:"closed_qty"`
	ClosedValue    decimal.Decimal `json:"closed_value"`
	LiqFee         decimal.Decimal `json:"liq_fee"`
	PositionSide   string          `json:"position_side"`
	Leverage       decimal.Decimal `json:"leverage"`
	Interest       decimal.Decimal `json:"interest"`
	BusinessType   string          `json:"business_type"`
	CreateTime     time.Time       `json:"create_time,string,format:unixmilli"`
	UpdateTime     time.Time       `json:"update_time,string,format:unixmilli"`
}

// ListHistoryMarginInterestsService -- GET /api/v4/crossex/history_margin_interests (private)
//
// Returns the account's leveraged-interest deduction history.
type ListHistoryMarginInterestsService struct {
	c      *CrossexClient
	params map[string]string
}

func (c *CrossexClient) NewListHistoryMarginInterestsService() *ListHistoryMarginInterestsService {
	return &ListHistoryMarginInterestsService{c: c, params: map[string]string{}}
}

// SetSymbol narrows the result to a single trading pair.
func (s *ListHistoryMarginInterestsService) SetSymbol(symbol string) *ListHistoryMarginInterestsService {
	s.params["symbol"] = symbol
	return s
}

// SetFrom sets the start time (millisecond precision).
func (s *ListHistoryMarginInterestsService) SetFrom(from time.Time) *ListHistoryMarginInterestsService {
	s.params["from"] = strconv.FormatInt(from.UnixMilli(), 10)
	return s
}

// SetTo sets the end time (millisecond precision).
func (s *ListHistoryMarginInterestsService) SetTo(to time.Time) *ListHistoryMarginInterestsService {
	s.params["to"] = strconv.FormatInt(to.UnixMilli(), 10)
	return s
}

// SetPage selects the result page.
func (s *ListHistoryMarginInterestsService) SetPage(page int) *ListHistoryMarginInterestsService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit caps the number of records returned in a single list.
func (s *ListHistoryMarginInterestsService) SetLimit(limit int) *ListHistoryMarginInterestsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetExchangeType narrows the result to a single venue.
func (s *ListHistoryMarginInterestsService) SetExchangeType(exchangeType string) *ListHistoryMarginInterestsService {
	s.params["exchange_type"] = exchangeType
	return s
}

func (s *ListHistoryMarginInterestsService) Do(ctx context.Context) ([]CrossexMarginInterestRecord, error) {
	req := request.Get(ctx, s.c, "/api/v4/crossex/history_margin_interests", s.params).WithSign()
	resp, err := request.Do[[]CrossexMarginInterestRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CrossexMarginInterestRecord is one leveraged-interest deduction record.
// create_time is a millisecond Unix timestamp. Note the userId key is camelCase
// on the wire (unlike the snake_case used elsewhere in this product).
type CrossexMarginInterestRecord struct {
	UserID        string          `json:"userId"`
	Symbol        string          `json:"symbol"`
	InterestID    string          `json:"interest_id"`
	LiabilityID   string          `json:"liability_id"`
	Liability     decimal.Decimal `json:"liability"`
	LiabilityCoin string          `json:"liability_coin"`
	Interest      decimal.Decimal `json:"interest"`
	InterestRate  decimal.Decimal `json:"interest_rate"`
	InterestType  string          `json:"interest_type"`
	CreateTime    time.Time       `json:"create_time,string,format:unixmilli"`
	ExchangeType  string          `json:"exchange_type"`
}
