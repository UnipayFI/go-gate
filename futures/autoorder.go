package futures

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// ============================================================================
// Trail (trailing) auto orders -- /api/v4/futures/{settle}/autoorder/v1/trail/*
// ============================================================================

// CreateTrailOrderService -- POST /api/v4/futures/{settle}/autoorder/v1/trail/create (private)
//
// Creates a trailing (trail) auto order that follows the market by a callback
// ratio or price distance. amount is the signed contract quantity (positive to
// buy, negative to sell).
type CreateTrailOrderService struct {
	c      *FuturesClient
	settle Settle
	body   map[string]any
}

func (c *FuturesClient) NewCreateTrailOrderService(settle Settle, contract string, amount decimal.Decimal) *CreateTrailOrderService {
	return &CreateTrailOrderService{c: c, settle: settle, body: map[string]any{
		"contract": contract,
		"amount":   amount.String(),
	}}
}

// SetActivationPrice sets the activation price (0 activates immediately).
func (s *CreateTrailOrderService) SetActivationPrice(activationPrice decimal.Decimal) *CreateTrailOrderService {
	s.body["activation_price"] = activationPrice.String()
	return s
}

// SetIsGte activates when the market price is >= the activation price (true) or
// <= the activation price (false).
func (s *CreateTrailOrderService) SetIsGte(isGte bool) *CreateTrailOrderService {
	s.body["is_gte"] = isGte
	return s
}

// SetPriceType selects the activation reference price (1 latest, 2 index, 3 mark).
func (s *CreateTrailOrderService) SetPriceType(priceType int) *CreateTrailOrderService {
	s.body["price_type"] = priceType
	return s
}

// SetPriceOffset sets the callback ratio or price distance, e.g. "0.1" or "0.1%".
func (s *CreateTrailOrderService) SetPriceOffset(priceOffset string) *CreateTrailOrderService {
	s.body["price_offset"] = priceOffset
	return s
}

// SetReduceOnly marks the order as reduce-only.
func (s *CreateTrailOrderService) SetReduceOnly(reduceOnly bool) *CreateTrailOrderService {
	s.body["reduce_only"] = reduceOnly
	return s
}

// SetPositionRelated binds the order to a position.
func (s *CreateTrailOrderService) SetPositionRelated(positionRelated bool) *CreateTrailOrderService {
	s.body["position_related"] = positionRelated
	return s
}

// SetText attaches custom order information identifying the order source.
func (s *CreateTrailOrderService) SetText(text string) *CreateTrailOrderService {
	s.body["text"] = text
	return s
}

// SetPosMarginMode sets the position margin mode ("isolated" or "cross").
func (s *CreateTrailOrderService) SetPosMarginMode(posMarginMode string) *CreateTrailOrderService {
	s.body["pos_margin_mode"] = posMarginMode
	return s
}

// SetPositionMode sets the position mode ("single", "dual" or "dual_plus").
func (s *CreateTrailOrderService) SetPositionMode(positionMode string) *CreateTrailOrderService {
	s.body["position_mode"] = positionMode
	return s
}

func (s *CreateTrailOrderService) Do(ctx context.Context) (*FuturesTrailCreateResponse, error) {
	req := request.Post(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/autoorder/v1/trail/create", s.body).WithSign()
	return request.Do[FuturesTrailCreateResponse](req)
}

// StopTrailOrderService -- POST /api/v4/futures/{settle}/autoorder/v1/trail/stop (private)
//
// Terminates a single trail order, identified by id or (when id is omitted) by
// the account's custom text.
type StopTrailOrderService struct {
	c      *FuturesClient
	settle Settle
	body   map[string]any
}

func (c *FuturesClient) NewStopTrailOrderService(settle Settle) *StopTrailOrderService {
	return &StopTrailOrderService{c: c, settle: settle, body: map[string]any{}}
}

// SetID identifies the trail order to terminate (when set, text is not needed).
func (s *StopTrailOrderService) SetID(id int64) *StopTrailOrderService {
	s.body["id"] = id
	return s
}

// SetText terminates the order by custom text when no id is provided.
func (s *StopTrailOrderService) SetText(text string) *StopTrailOrderService {
	s.body["text"] = text
	return s
}

func (s *StopTrailOrderService) Do(ctx context.Context) (*FuturesTrailOrderResult, error) {
	req := request.Post(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/autoorder/v1/trail/stop", s.body).WithSign()
	return request.Do[FuturesTrailOrderResult](req)
}

// StopAllTrailOrdersService -- POST /api/v4/futures/{settle}/autoorder/v1/trail/stop_all (private)
//
// Batch-terminates trail orders, optionally limited to one contract and/or the
// orders bound to a single position.
type StopAllTrailOrdersService struct {
	c      *FuturesClient
	settle Settle
	body   map[string]any
}

func (c *FuturesClient) NewStopAllTrailOrdersService(settle Settle) *StopAllTrailOrdersService {
	return &StopAllTrailOrdersService{c: c, settle: settle, body: map[string]any{}}
}

// SetContract limits the termination to a single futures contract.
func (s *StopAllTrailOrdersService) SetContract(contract string) *StopAllTrailOrdersService {
	s.body["contract"] = contract
	return s
}

// SetRelatedPosition cancels only orders bound to the given position (1 long, 2 short).
func (s *StopAllTrailOrdersService) SetRelatedPosition(relatedPosition int) *StopAllTrailOrdersService {
	s.body["related_position"] = relatedPosition
	return s
}

func (s *StopAllTrailOrdersService) Do(ctx context.Context) (*FuturesTrailOrdersResult, error) {
	req := request.Post(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/autoorder/v1/trail/stop_all", s.body).WithSign()
	return request.Do[FuturesTrailOrdersResult](req)
}

// ListTrailOrdersService -- GET /api/v4/futures/{settle}/autoorder/v1/trail/list (private)
//
// Lists the account's trail orders, filtered by contract, status, time range,
// side and other optional criteria.
type ListTrailOrdersService struct {
	c      *FuturesClient
	settle Settle
	params map[string]string
}

func (c *FuturesClient) NewListTrailOrdersService(settle Settle) *ListTrailOrdersService {
	return &ListTrailOrdersService{c: c, settle: settle, params: map[string]string{}}
}

// SetContract narrows the result to a single futures contract.
func (s *ListTrailOrdersService) SetContract(contract string) *ListTrailOrdersService {
	s.params["contract"] = contract
	return s
}

// SetIsFinished selects historical (true) or in-progress (false) orders.
func (s *ListTrailOrdersService) SetIsFinished(isFinished bool) *ListTrailOrdersService {
	s.params["is_finished"] = strconv.FormatBool(isFinished)
	return s
}

// SetStartAt sets the start of the query time range.
func (s *ListTrailOrdersService) SetStartAt(startAt int64) *ListTrailOrdersService {
	s.params["start_at"] = strconv.FormatInt(startAt, 10)
	return s
}

// SetEndAt sets the end of the query time range.
func (s *ListTrailOrdersService) SetEndAt(endAt int64) *ListTrailOrdersService {
	s.params["end_at"] = strconv.FormatInt(endAt, 10)
	return s
}

// SetPageNum sets the page number, starting from 1.
func (s *ListTrailOrdersService) SetPageNum(pageNum int) *ListTrailOrdersService {
	s.params["page_num"] = strconv.Itoa(pageNum)
	return s
}

// SetPageSize sets the number of items per page.
func (s *ListTrailOrdersService) SetPageSize(pageSize int) *ListTrailOrdersService {
	s.params["page_size"] = strconv.Itoa(pageSize)
	return s
}

// SetSortBy selects the sort field (1 creation time, 2 end time).
func (s *ListTrailOrdersService) SetSortBy(sortBy int) *ListTrailOrdersService {
	s.params["sort_by"] = strconv.Itoa(sortBy)
	return s
}

// SetHideCancel hides cancelled orders when true.
func (s *ListTrailOrdersService) SetHideCancel(hideCancel bool) *ListTrailOrdersService {
	s.params["hide_cancel"] = strconv.FormatBool(hideCancel)
	return s
}

// SetRelatedPosition returns only orders bound to the given position (1 long, 2 short).
func (s *ListTrailOrdersService) SetRelatedPosition(relatedPosition int) *ListTrailOrdersService {
	s.params["related_position"] = strconv.Itoa(relatedPosition)
	return s
}

// SetSortByTrigger sorts by trigger/activation price so orders closest to firing
// come first (current orders only).
func (s *ListTrailOrdersService) SetSortByTrigger(sortByTrigger bool) *ListTrailOrdersService {
	s.params["sort_by_trigger"] = strconv.FormatBool(sortByTrigger)
	return s
}

// SetReduceOnly filters by reduce-only flag (1 yes, 2 no).
func (s *ListTrailOrdersService) SetReduceOnly(reduceOnly int) *ListTrailOrdersService {
	s.params["reduce_only"] = strconv.Itoa(reduceOnly)
	return s
}

// SetSide filters by direction (1 long position, 2 short position).
func (s *ListTrailOrdersService) SetSide(side int) *ListTrailOrdersService {
	s.params["side"] = strconv.Itoa(side)
	return s
}

func (s *ListTrailOrdersService) Do(ctx context.Context) (*FuturesTrailOrdersResult, error) {
	req := request.Get(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/autoorder/v1/trail/list", s.params).WithSign()
	return request.Do[FuturesTrailOrdersResult](req)
}

// GetTrailOrderService -- GET /api/v4/futures/{settle}/autoorder/v1/trail/detail (private)
//
// Returns the details of a single trail order by its id.
type GetTrailOrderService struct {
	c      *FuturesClient
	settle Settle
	params map[string]string
}

func (c *FuturesClient) NewGetTrailOrderService(settle Settle, id int64) *GetTrailOrderService {
	return &GetTrailOrderService{c: c, settle: settle, params: map[string]string{
		"id": strconv.FormatInt(id, 10),
	}}
}

func (s *GetTrailOrderService) Do(ctx context.Context) (*FuturesTrailDetailResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/autoorder/v1/trail/detail", s.params).WithSign()
	return request.Do[FuturesTrailDetailResponse](req)
}

// UpdateTrailOrderService -- POST /api/v4/futures/{settle}/autoorder/v1/trail/update (private)
//
// Modifies an existing trail order in place. Only the fields you set are changed;
// unset (empty/zero) fields are left as they are.
type UpdateTrailOrderService struct {
	c      *FuturesClient
	settle Settle
	body   map[string]any
}

func (c *FuturesClient) NewUpdateTrailOrderService(settle Settle, id int64) *UpdateTrailOrderService {
	return &UpdateTrailOrderService{c: c, settle: settle, body: map[string]any{
		"id": id,
	}}
}

// SetAmount sets the new total signed contract quantity (0 means no modification).
func (s *UpdateTrailOrderService) SetAmount(amount decimal.Decimal) *UpdateTrailOrderService {
	s.body["amount"] = amount.String()
	return s
}

// SetActivationPrice sets the new activation price (0 activates immediately,
// empty means no modification).
func (s *UpdateTrailOrderService) SetActivationPrice(activationPrice decimal.Decimal) *UpdateTrailOrderService {
	s.body["activation_price"] = activationPrice.String()
	return s
}

// SetIsGteStr sets the activation direction as a string ("true", "false", or
// empty for no modification).
func (s *UpdateTrailOrderService) SetIsGteStr(isGteStr string) *UpdateTrailOrderService {
	s.body["is_gte_str"] = isGteStr
	return s
}

// SetPriceType sets the new activation reference price (0 no modification, 1
// latest, 2 index, 3 mark).
func (s *UpdateTrailOrderService) SetPriceType(priceType int) *UpdateTrailOrderService {
	s.body["price_type"] = priceType
	return s
}

// SetPriceOffset sets the new callback ratio or price distance (empty means no
// modification), e.g. "0.1" or "0.1%".
func (s *UpdateTrailOrderService) SetPriceOffset(priceOffset string) *UpdateTrailOrderService {
	s.body["price_offset"] = priceOffset
	return s
}

func (s *UpdateTrailOrderService) Do(ctx context.Context) (*FuturesTrailOrderResult, error) {
	req := request.Post(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/autoorder/v1/trail/update", s.body).WithSign()
	return request.Do[FuturesTrailOrderResult](req)
}

// GetTrailChangeLogService -- GET /api/v4/futures/{settle}/autoorder/v1/trail/change_log (private)
//
// Returns the user modification records (create/update history) of a trail order.
type GetTrailChangeLogService struct {
	c      *FuturesClient
	settle Settle
	params map[string]string
}

func (c *FuturesClient) NewGetTrailChangeLogService(settle Settle, id int64) *GetTrailChangeLogService {
	return &GetTrailChangeLogService{c: c, settle: settle, params: map[string]string{
		"id": strconv.FormatInt(id, 10),
	}}
}

// SetPageNum sets the page number, starting from 1.
func (s *GetTrailChangeLogService) SetPageNum(pageNum int) *GetTrailChangeLogService {
	s.params["page_num"] = strconv.Itoa(pageNum)
	return s
}

// SetPageSize sets the number of items per page.
func (s *GetTrailChangeLogService) SetPageSize(pageSize int) *GetTrailChangeLogService {
	s.params["page_size"] = strconv.Itoa(pageSize)
	return s
}

func (s *GetTrailChangeLogService) Do(ctx context.Context) (*FuturesTrailChangeLogResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/autoorder/v1/trail/change_log", s.params).WithSign()
	return request.Do[FuturesTrailChangeLogResponse](req)
}

// ============================================================================
// Chase auto orders -- /api/v4/futures/{settle}/autoorder/v1/chase/*
// ============================================================================

// CreateChaseOrderService -- POST /api/v4/futures/{settle}/autoorder/v1/chase/create (private)
//
// Creates a chase order that repeatedly re-prices to follow the best bid/ask up
// to a price limit. amount is the signed contract quantity (positive to buy,
// negative to sell); priceLimit is the maximum chase price ("0" for no limit).
type CreateChaseOrderService struct {
	c      *FuturesClient
	settle Settle
	body   map[string]any
}

func (c *FuturesClient) NewCreateChaseOrderService(settle Settle, contract string, amount, priceLimit decimal.Decimal) *CreateChaseOrderService {
	return &CreateChaseOrderService{c: c, settle: settle, body: map[string]any{
		"contract":    contract,
		"amount":      amount.String(),
		"price_limit": priceLimit.String(),
	}}
}

// SetSettle overrides the body settle currency (the path parameter takes precedence).
func (s *CreateChaseOrderService) SetSettle(settle string) *CreateChaseOrderService {
	s.body["settle"] = settle
	return s
}

// SetOffsetLimit sets the maximum chasing distance from the best price (mutually
// exclusive with price_limit).
func (s *CreateChaseOrderService) SetOffsetLimit(offsetLimit decimal.Decimal) *CreateChaseOrderService {
	s.body["offset_limit"] = offsetLimit.String()
	return s
}

// SetReduceOnly marks the order as reduce-only.
func (s *CreateChaseOrderService) SetReduceOnly(reduceOnly bool) *CreateChaseOrderService {
	s.body["reduce_only"] = reduceOnly
	return s
}

// SetText attaches an optional custom tag.
func (s *CreateChaseOrderService) SetText(text string) *CreateChaseOrderService {
	s.body["text"] = text
	return s
}

// SetIsDualMode enables dual-position mode.
func (s *CreateChaseOrderService) SetIsDualMode(isDualMode bool) *CreateChaseOrderService {
	s.body["is_dual_mode"] = isDualMode
	return s
}

// SetPriceType selects the price type (1 best bid/ask, 2 distance from best bid/ask).
func (s *CreateChaseOrderService) SetPriceType(priceType int64) *CreateChaseOrderService {
	s.body["price_type"] = priceType
	return s
}

// SetPriceGapType selects the gap unit when price_type == 2 (1 absolute price gap,
// 2 percentage).
func (s *CreateChaseOrderService) SetPriceGapType(priceGapType int64) *CreateChaseOrderService {
	s.body["price_gap_type"] = priceGapType
	return s
}

// SetPriceGapValue sets the price gap value paired with price_gap_type.
func (s *CreateChaseOrderService) SetPriceGapValue(priceGapValue string) *CreateChaseOrderService {
	s.body["price_gap_value"] = priceGapValue
	return s
}

// SetPosMarginMode sets the position margin mode ("isolated" or "cross").
func (s *CreateChaseOrderService) SetPosMarginMode(posMarginMode string) *CreateChaseOrderService {
	s.body["pos_margin_mode"] = posMarginMode
	return s
}

// SetPositionMode sets the position mode ("single", "dual" or "dual_plus").
func (s *CreateChaseOrderService) SetPositionMode(positionMode string) *CreateChaseOrderService {
	s.body["position_mode"] = positionMode
	return s
}

func (s *CreateChaseOrderService) Do(ctx context.Context) (*FuturesChaseCreateResponse, error) {
	req := request.Post(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/autoorder/v1/chase/create", s.body).WithSign()
	return request.Do[FuturesChaseCreateResponse](req)
}

// StopChaseOrderService -- POST /api/v4/futures/{settle}/autoorder/v1/chase/stop (private)
//
// Stops a single chase order, identified by id or (when id is omitted) by the
// account's custom text.
type StopChaseOrderService struct {
	c      *FuturesClient
	settle Settle
	body   map[string]any
}

func (c *FuturesClient) NewStopChaseOrderService(settle Settle) *StopChaseOrderService {
	return &StopChaseOrderService{c: c, settle: settle, body: map[string]any{}}
}

// SetID identifies the chase order to stop (either id or text must be provided).
func (s *StopChaseOrderService) SetID(id string) *StopChaseOrderService {
	s.body["id"] = id
	return s
}

// SetText stops the order by custom text when no id is provided.
func (s *StopChaseOrderService) SetText(text string) *StopChaseOrderService {
	s.body["text"] = text
	return s
}

// SetSettle overrides the body settle currency (the path parameter takes precedence).
func (s *StopChaseOrderService) SetSettle(settle string) *StopChaseOrderService {
	s.body["settle"] = settle
	return s
}

func (s *StopChaseOrderService) Do(ctx context.Context) (*FuturesChaseOrderResult, error) {
	req := request.Post(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/autoorder/v1/chase/stop", s.body).WithSign()
	return request.Do[FuturesChaseOrderResult](req)
}

// StopAllChaseOrdersService -- POST /api/v4/futures/{settle}/autoorder/v1/chase/stop_all (private)
//
// Stops chase orders in batch, optionally limited to one contract and/or margin
// mode.
type StopAllChaseOrdersService struct {
	c      *FuturesClient
	settle Settle
	body   map[string]any
}

func (c *FuturesClient) NewStopAllChaseOrdersService(settle Settle) *StopAllChaseOrdersService {
	return &StopAllChaseOrdersService{c: c, settle: settle, body: map[string]any{}}
}

// SetContract limits the batch stop to a single futures contract.
func (s *StopAllChaseOrdersService) SetContract(contract string) *StopAllChaseOrdersService {
	s.body["contract"] = contract
	return s
}

// SetSettle overrides the body settle currency (the path parameter takes precedence).
func (s *StopAllChaseOrdersService) SetSettle(settle string) *StopAllChaseOrdersService {
	s.body["settle"] = settle
	return s
}

// SetPosMarginMode limits the batch stop to a margin mode ("isolated" or "cross").
func (s *StopAllChaseOrdersService) SetPosMarginMode(posMarginMode string) *StopAllChaseOrdersService {
	s.body["pos_margin_mode"] = posMarginMode
	return s
}

func (s *StopAllChaseOrdersService) Do(ctx context.Context) (*FuturesChaseOrdersResult, error) {
	req := request.Post(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/autoorder/v1/chase/stop_all", s.body).WithSign()
	return request.Do[FuturesChaseOrdersResult](req)
}

// ListChaseOrdersService -- GET /api/v4/futures/{settle}/autoorder/v1/chase/list (private)
//
// Lists the account's chase orders. sortBy is required (1 by created time, 2 by
// finished time).
type ListChaseOrdersService struct {
	c      *FuturesClient
	settle Settle
	params map[string]string
}

func (c *FuturesClient) NewListChaseOrdersService(settle Settle, sortBy int) *ListChaseOrdersService {
	return &ListChaseOrdersService{c: c, settle: settle, params: map[string]string{
		"sort_by": strconv.Itoa(sortBy),
	}}
}

// SetContract narrows the result to a single futures contract.
func (s *ListChaseOrdersService) SetContract(contract string) *ListChaseOrdersService {
	s.params["contract"] = contract
	return s
}

// SetIsFinished selects finished (true) or in-progress (false) orders.
func (s *ListChaseOrdersService) SetIsFinished(isFinished bool) *ListChaseOrdersService {
	s.params["is_finished"] = strconv.FormatBool(isFinished)
	return s
}

// SetStartAt sets the lower time bound of the history list (paired with end_at).
func (s *ListChaseOrdersService) SetStartAt(startAt int64) *ListChaseOrdersService {
	s.params["start_at"] = strconv.FormatInt(startAt, 10)
	return s
}

// SetEndAt sets the upper time bound of the history list (paired with start_at).
func (s *ListChaseOrdersService) SetEndAt(endAt int64) *ListChaseOrdersService {
	s.params["end_at"] = strconv.FormatInt(endAt, 10)
	return s
}

// SetPageNum sets the page number, starting from 1.
func (s *ListChaseOrdersService) SetPageNum(pageNum int) *ListChaseOrdersService {
	s.params["page_num"] = strconv.Itoa(pageNum)
	return s
}

// SetPageSize sets the page size (must be between 1 and 100).
func (s *ListChaseOrdersService) SetPageSize(pageSize int) *ListChaseOrdersService {
	s.params["page_size"] = strconv.Itoa(pageSize)
	return s
}

// SetHideCancel hides cancelled orders when true.
func (s *ListChaseOrdersService) SetHideCancel(hideCancel bool) *ListChaseOrdersService {
	s.params["hide_cancel"] = strconv.FormatBool(hideCancel)
	return s
}

// SetReduceOnly filters by reduce-only flag (0 unknown, 1 true, 2 false).
func (s *ListChaseOrdersService) SetReduceOnly(reduceOnly int) *ListChaseOrdersService {
	s.params["reduce_only"] = strconv.Itoa(reduceOnly)
	return s
}

// SetSide filters by long/short side (1 long, 2 short).
func (s *ListChaseOrdersService) SetSide(side int) *ListChaseOrdersService {
	s.params["side"] = strconv.Itoa(side)
	return s
}

func (s *ListChaseOrdersService) Do(ctx context.Context) (*FuturesChaseOrdersResult, error) {
	req := request.Get(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/autoorder/v1/chase/list", s.params).WithSign()
	return request.Do[FuturesChaseOrdersResult](req)
}

// GetChaseOrderService -- GET /api/v4/futures/{settle}/autoorder/v1/chase/detail (private)
//
// Returns the details of a single chase order by its id.
type GetChaseOrderService struct {
	c      *FuturesClient
	settle Settle
	params map[string]string
}

func (c *FuturesClient) NewGetChaseOrderService(settle Settle, id string) *GetChaseOrderService {
	return &GetChaseOrderService{c: c, settle: settle, params: map[string]string{
		"id": id,
	}}
}

func (s *GetChaseOrderService) Do(ctx context.Context) (*FuturesChaseOrderResult, error) {
	req := request.Get(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/autoorder/v1/chase/detail", s.params).WithSign()
	return request.Do[FuturesChaseOrderResult](req)
}

// ============================================================================
// Response / model types
// ============================================================================

// FuturesAutoOrderEnvelope is the {code,message,data,timestamp} business envelope
// every futures autoorder (trail/chase) endpoint wraps its payload in — unlike
// Gate's core v4 API, which returns the payload directly. timestamp is a
// millisecond Unix timestamp; Code 0 means success.
type FuturesAutoOrderEnvelope[T any] struct {
	Code      int       `json:"code"`
	Message   string    `json:"message"`
	Data      T         `json:"data"`
	Timestamp time.Time `json:"timestamp,format:unixmilli"`
}

// FuturesTrailCreateData holds the id of the newly created trail order.
type FuturesTrailCreateData struct {
	ID string `json:"id"`
}

// FuturesTrailDetailData holds a single trail order (create/detail/stop/update data).
type FuturesTrailDetailData struct {
	Order FuturesTrailOrder `json:"order"`
}

// FuturesTrailOrdersData holds a list of trail orders (list/stop_all data).
type FuturesTrailOrdersData struct {
	Orders []FuturesTrailOrder `json:"orders"`
}

// FuturesTrailChangeLogData holds a trail order's modification records.
type FuturesTrailChangeLogData struct {
	ChangeLog []FuturesTrailChangeLog `json:"change_log"`
}

// FuturesTrailCreateResponse is the create-trail-order envelope (data.id).
type FuturesTrailCreateResponse = FuturesAutoOrderEnvelope[FuturesTrailCreateData]

// FuturesTrailOrderResult is the single-trail-order envelope (stop / update, data.order).
type FuturesTrailOrderResult = FuturesAutoOrderEnvelope[FuturesTrailDetailData]

// FuturesTrailOrdersResult is the trail-order-list envelope (list / stop_all, data.orders).
type FuturesTrailOrdersResult = FuturesAutoOrderEnvelope[FuturesTrailOrdersData]

// FuturesTrailDetailResponse is the trail-order-detail envelope (data.order).
type FuturesTrailDetailResponse = FuturesAutoOrderEnvelope[FuturesTrailDetailData]

// FuturesTrailChangeLogResponse is the trail-order change-log envelope (data.change_log).
type FuturesTrailChangeLogResponse = FuturesAutoOrderEnvelope[FuturesTrailChangeLogData]

// FuturesTrailOrder is a trailing (trail) auto order and its live state. The
// integer timestamps are Unix seconds; the *_precise fields carry the same times
// as high-precision "seconds.microseconds" strings. price_offset is a callback
// ratio or price distance that may be a percentage (e.g. "0.1%").
type FuturesTrailOrder struct {
	ID                 int64           `json:"id"`
	UserID             int64           `json:"user_id"`
	User               int64           `json:"user"`
	Contract           string          `json:"contract"`
	Settle             string          `json:"settle"`
	Amount             decimal.Decimal `json:"amount"`
	IsGte              bool            `json:"is_gte"`
	ActivationPrice    decimal.Decimal `json:"activation_price"`
	PriceType          int             `json:"price_type"`
	PriceOffset        string          `json:"price_offset"`
	Text               string          `json:"text"`
	ReduceOnly         bool            `json:"reduce_only"`
	PositionRelated    bool            `json:"position_related"`
	CreatedAt          time.Time       `json:"created_at,format:unix"`
	ActivatedAt        time.Time       `json:"activated_at,format:unix"`
	FinishedAt         time.Time       `json:"finished_at,format:unix"`
	CreateTime         time.Time       `json:"create_time,format:unix"`
	ActiveTime         time.Time       `json:"active_time,format:unix"`
	FinishTime         time.Time       `json:"finish_time,format:unix"`
	Reason             string          `json:"reason"`
	SuborderText       string          `json:"suborder_text"`
	IsDualMode         bool            `json:"is_dual_mode"`
	TriggerPrice       decimal.Decimal `json:"trigger_price"`
	SuborderID         int64           `json:"suborder_id"`
	SideLabel          string          `json:"side_label"`
	OriginalStatus     int             `json:"original_status"`
	Status             OrderStatus     `json:"status"`
	PositionSideOutput string          `json:"position_side_output"`
	UpdatedAt          time.Time       `json:"updated_at,format:unix"`
	ExtremumPrice      decimal.Decimal `json:"extremum_price"`
	StatusCode         string          `json:"status_code"`
	CreatedAtPrecise   string          `json:"created_at_precise"`
	FinishedAtPrecise  string          `json:"finished_at_precise"`
	ActivatedAtPrecise string          `json:"activated_at_precise"`
	StatusLabel        string          `json:"status_label"`
	PosMarginMode      string          `json:"pos_margin_mode"`
	PositionMode       string          `json:"position_mode"`
	ErrorLabel         string          `json:"error_label"`
	Leverage           decimal.Decimal `json:"leverage"`
}

// FuturesTrailChangeLog is one create/modify record of a trail order. updated_at
// is a Unix-second timestamp; is_create marks the record as the creation (true)
// or a later modification (false).
type FuturesTrailChangeLog struct {
	UpdatedAt       time.Time       `json:"updated_at,format:unix"`
	Amount          decimal.Decimal `json:"amount"`
	IsGte           bool            `json:"is_gte"`
	ActivationPrice decimal.Decimal `json:"activation_price"`
	PriceType       int             `json:"price_type"`
	PriceOffset     string          `json:"price_offset"`
	IsCreate        bool            `json:"is_create"`
}

// FuturesChaseCreateData holds the id of the newly created chase order.
type FuturesChaseCreateData struct {
	ID string `json:"id"`
}

// FuturesChaseOrderData holds a single chase order (stop/detail data).
type FuturesChaseOrderData struct {
	Order FuturesChaseOrder `json:"order"`
}

// FuturesChaseOrdersData holds a list of chase orders (list/stop_all data).
type FuturesChaseOrdersData struct {
	Orders []FuturesChaseOrder `json:"orders"`
}

// FuturesChaseCreateResponse is the create-chase-order envelope (data.id).
type FuturesChaseCreateResponse = FuturesAutoOrderEnvelope[FuturesChaseCreateData]

// FuturesChaseOrderResult is the single-chase-order envelope (stop / detail, data.order).
type FuturesChaseOrderResult = FuturesAutoOrderEnvelope[FuturesChaseOrderData]

// FuturesChaseOrdersResult is the chase-order-list envelope (list / stop_all, data.orders).
type FuturesChaseOrdersResult = FuturesAutoOrderEnvelope[FuturesChaseOrdersData]

// FuturesChaseOrder is a chase auto order and its live state. amount is the
// signed contract quantity (positive buy, negative sell); create_time /
// finish_time / updated_at are Unix-second timestamps, and the *_precise fields
// carry the same times as high-precision "seconds.microseconds" strings.
type FuturesChaseOrder struct {
	ID                 string          `json:"id"`
	User               string          `json:"user"`
	Contract           string          `json:"contract"`
	Settle             string          `json:"settle"`
	Amount             decimal.Decimal `json:"amount"`
	PriceLimit         decimal.Decimal `json:"price_limit"`
	ReduceOnly         bool            `json:"reduce_only"`
	Text               string          `json:"text"`
	CreateTime         time.Time       `json:"create_time,format:unix"`
	FinishTime         time.Time       `json:"finish_time,format:unix"`
	OriginalStatus     int             `json:"original_status"`
	Status             OrderStatus     `json:"status"`
	Reason             string          `json:"reason"`
	FillAmount         decimal.Decimal `json:"fill_amount"`
	AverageFillPrice   decimal.Decimal `json:"average_fill_price"`
	SuborderID         string          `json:"suborder_id"`
	IsDualMode         bool            `json:"is_dual_mode"`
	SideLabel          string          `json:"side_label"`
	PositionSideOutput string          `json:"position_side_output"`
	ChasePrice         decimal.Decimal `json:"chase_price"`
	IntervalSec        int             `json:"interval_sec"`
	UpdatedAt          time.Time       `json:"updated_at,format:unix"`
	SuborderPrice      decimal.Decimal `json:"suborder_price"`
	SuborderOngoing    bool            `json:"suborder_ongoing"`
	SuborderFinishAs   string          `json:"suborder_finish_as"`
	PriceType          int             `json:"price_type"`
	PriceGapType       string          `json:"price_gap_type"`
	PriceGapValue      string          `json:"price_gap_value"`
	StatusCode         string          `json:"status_code"`
	CreateTimePrecise  string          `json:"create_time_precise"`
	FinishTimePrecise  string          `json:"finish_time_precise"`
	PosMarginMode      string          `json:"pos_margin_mode"`
	PositionMode       string          `json:"position_mode"`
	Leverage           decimal.Decimal `json:"leverage"`
	ErrorLabel         string          `json:"error_label"`
}
