package tradfi

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// ListPositionsService -- GET /api/v4/tradfi/positions (private)
//
// Lists the caller's active positions.
type ListPositionsService struct {
	c *TradfiClient
}

func (c *TradfiClient) NewListPositionsService() *ListPositionsService {
	return &ListPositionsService{c: c}
}

func (s *ListPositionsService) Do(ctx context.Context) (*TradfiPositionsResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/tradfi/positions").WithSign()
	return request.Do[TradfiPositionsResponse](req)
}

// TradfiPositionsResponse is the envelope of the active-position query.
type TradfiPositionsResponse struct {
	Label     string    `json:"label"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp,format:unixmilli"`
	Data      struct {
		List []TradfiPosition `json:"list"`
	} `json:"data"`
}

// TradfiPosition is a single active position. position_dir is "Long" or "Short".
type TradfiPosition struct {
	PositionID        int64           `json:"position_id"`
	Symbol            string          `json:"symbol"`
	SymbolDesc        string          `json:"symbol_desc"`
	Margin            decimal.Decimal `json:"margin"`
	UnrealizedPnL     decimal.Decimal `json:"unrealized_pnl"`
	UnrealizedPnLRate decimal.Decimal `json:"unrealized_pnl_rate"`
	Volume            decimal.Decimal `json:"volume"`
	PriceOpen         decimal.Decimal `json:"price_open"`
	PositionDir       string          `json:"position_dir"`
}

// ModifyPositionService -- PUT /api/v4/tradfi/positions/{position_id} (private)
//
// Modifies a position's take-profit and/or stop-loss prices.
type ModifyPositionService struct {
	c          *TradfiClient
	positionID int64
	body       map[string]any
}

func (c *TradfiClient) NewModifyPositionService(positionID int64) *ModifyPositionService {
	return &ModifyPositionService{c: c, positionID: positionID, body: map[string]any{}}
}

// SetPriceTP sets the take-profit price. Passing "0" clears the existing one.
func (s *ModifyPositionService) SetPriceTP(priceTP decimal.Decimal) *ModifyPositionService {
	s.body["price_tp"] = priceTP.String()
	return s
}

// SetPriceSL sets the stop-loss price. Passing "0" clears the existing one.
func (s *ModifyPositionService) SetPriceSL(priceSL decimal.Decimal) *ModifyPositionService {
	s.body["price_sl"] = priceSL.String()
	return s
}

func (s *ModifyPositionService) Do(ctx context.Context) (*TradfiModifyPositionResponse, error) {
	req := request.Put(ctx, s.c, "/api/v4/tradfi/positions/"+strconv.FormatInt(s.positionID, 10), s.body).WithSign()
	return request.Do[TradfiModifyPositionResponse](req)
}

// TradfiModifyPositionResponse is the envelope of the modify-position request.
// data is an empty object on success.
type TradfiModifyPositionResponse struct {
	Timestamp time.Time `json:"timestamp,format:unixmilli"`
	Data      struct{}  `json:"data"`
}

// ClosePositionService -- POST /api/v4/tradfi/positions/{position_id}/close (private)
//
// Closes a position. closeType is 1 for a partial close (closeVolume required)
// or 2 for a full close (closeVolume ignored).
type ClosePositionService struct {
	c          *TradfiClient
	positionID int64
	body       map[string]any
}

func (c *TradfiClient) NewClosePositionService(positionID int64, closeType int) *ClosePositionService {
	return &ClosePositionService{c: c, positionID: positionID, body: map[string]any{
		"close_type": closeType,
	}}
}

// SetCloseVolume sets the volume to close (required when closeType is 1).
func (s *ClosePositionService) SetCloseVolume(closeVolume decimal.Decimal) *ClosePositionService {
	s.body["close_volume"] = closeVolume.String()
	return s
}

func (s *ClosePositionService) Do(ctx context.Context) (*TradfiClosePositionResponse, error) {
	req := request.Post(ctx, s.c, "/api/v4/tradfi/positions/"+strconv.FormatInt(s.positionID, 10)+"/close", s.body).WithSign()
	return request.Do[TradfiClosePositionResponse](req)
}

// TradfiClosePositionResponse is the envelope of the close-position request.
// data is an empty object on success.
type TradfiClosePositionResponse struct {
	Timestamp time.Time `json:"timestamp,format:unixmilli"`
	Data      struct{}  `json:"data"`
}

// ListPositionHistoryService -- GET /api/v4/tradfi/positions/history (private)
//
// Lists the caller's closed-position history (earliest queryable one month ago).
type ListPositionHistoryService struct {
	c      *TradfiClient
	params map[string]string
}

func (c *TradfiClient) NewListPositionHistoryService() *ListPositionHistoryService {
	return &ListPositionHistoryService{c: c, params: map[string]string{}}
}

// SetPage selects the result page (defaults to 1).
func (s *ListPositionHistoryService) SetPage(page int64) *ListPositionHistoryService {
	s.params["page"] = strconv.FormatInt(page, 10)
	return s
}

// SetPageSize caps the number of records per page (defaults to 10, max 100).
func (s *ListPositionHistoryService) SetPageSize(pageSize int64) *ListPositionHistoryService {
	s.params["page_size"] = strconv.FormatInt(pageSize, 10)
	return s
}

// SetBeginTime bounds the result to positions at or after this time.
func (s *ListPositionHistoryService) SetBeginTime(beginTime time.Time) *ListPositionHistoryService {
	s.params["begin_time"] = strconv.FormatInt(beginTime.Unix(), 10)
	return s
}

// SetEndTime bounds the result to positions at or before this time.
func (s *ListPositionHistoryService) SetEndTime(endTime time.Time) *ListPositionHistoryService {
	s.params["end_time"] = strconv.FormatInt(endTime.Unix(), 10)
	return s
}

// SetSymbol narrows the result to a single symbol (e.g. "EURUSD").
func (s *ListPositionHistoryService) SetSymbol(symbol string) *ListPositionHistoryService {
	s.params["symbol"] = symbol
	return s
}

// SetPositionDir narrows the result to a direction ("Long" or "Short").
func (s *ListPositionHistoryService) SetPositionDir(positionDir string) *ListPositionHistoryService {
	s.params["position_dir"] = positionDir
	return s
}

func (s *ListPositionHistoryService) Do(ctx context.Context) (*TradfiPositionHistoryResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/tradfi/positions/history", s.params).WithSign()
	return request.Do[TradfiPositionHistoryResponse](req)
}

// TradfiPositionHistoryResponse is the envelope of the closed-position query.
type TradfiPositionHistoryResponse struct {
	Timestamp time.Time `json:"timestamp,format:unixmilli"`
	Data      struct {
		Total     int                     `json:"total"`
		TotalPage int                     `json:"total_page"`
		List      []TradfiPositionHistory `json:"list"`
	} `json:"data"`
}

// TradfiPositionHistory is a single closed-position record. time_create /
// time_close are string-encoded integer-second Unix timestamps. close_detail is
// null for a normal close and populated on forced liquidation.
type TradfiPositionHistory struct {
	PositionID        int64                      `json:"position_id"`
	Symbol            string                     `json:"symbol"`
	RealizedPnL       decimal.Decimal            `json:"realized_pnl"`
	RealizedPnLRate   decimal.Decimal            `json:"realized_pnl_rate"`
	Volume            decimal.Decimal            `json:"volume"`
	VolumeClosed      decimal.Decimal            `json:"volume_closed"`
	PriceOpen         decimal.Decimal            `json:"price_open"`
	PositionDir       string                     `json:"position_dir"`
	PriceTP           decimal.Decimal            `json:"price_tp"`
	PriceSL           decimal.Decimal            `json:"price_sl"`
	CounterpartyPrice decimal.Decimal            `json:"counterparty_price"`
	ClosePrice        decimal.Decimal            `json:"close_price"`
	TimeCreate        time.Time                  `json:"time_create,string,format:unix"`
	TimeClose         time.Time                  `json:"time_close,string,format:unix"`
	PositionStatus    string                     `json:"position_status"`
	CloseDetail       *TradfiPositionCloseDetail `json:"close_detail"`
	RealizedPnLDetail TradfiRealizedPnLDetail    `json:"realized_pnl_detail"`
}

// TradfiPositionCloseDetail is the liquidation breakdown of a forced close.
// margin_level and stop_out_level are percentages multiplied by 100.
type TradfiPositionCloseDetail struct {
	MarginLevel  decimal.Decimal `json:"margin_level"`
	Margin       decimal.Decimal `json:"margin"`
	Equity       decimal.Decimal `json:"equity"`
	StopOutLevel decimal.Decimal `json:"stop_out_level"`
}

// TradfiRealizedPnLDetail is the realized-PnL breakdown of a closed position.
type TradfiRealizedPnLDetail struct {
	ClosedPnL decimal.Decimal `json:"closed_pnl"`
	Swap      decimal.Decimal `json:"swap"`
	Fee       decimal.Decimal `json:"fee"`
}
