package spot

// TradeStatus is the tradability state of a spot pair.
type TradeStatus string

const (
	TradeStatusTradable   TradeStatus = "tradable"
	TradeStatusUntradable TradeStatus = "untradable"
	TradeStatusBuyable    TradeStatus = "buyable"
	TradeStatusSellable   TradeStatus = "sellable"
)

// Side is the order direction.
type Side string

const (
	SideBuy  Side = "buy"
	SideSell Side = "sell"
)

// OrderType is the order execution method.
type OrderType string

const (
	OrderTypeLimit  OrderType = "limit"
	OrderTypeMarket OrderType = "market"
)

// Account selects which balance an order draws on.
type Account string

const (
	AccountSpot        Account = "spot"
	AccountMargin      Account = "margin"
	AccountCrossMargin Account = "cross_margin"
	AccountUnified     Account = "unified"
	// AccountNormal is spot price-triggered orders' name for the spot account
	// (put.account uses "normal" rather than "spot").
	AccountNormal Account = "normal"
)

// TimeInForce determines how long an order stays active.
type TimeInForce string

const (
	TimeInForceGTC TimeInForce = "gtc" // GoodTillCancelled
	TimeInForceIOC TimeInForce = "ioc" // ImmediateOrCancelled
	TimeInForcePOC TimeInForce = "poc" // PendingOrCancelled (post-only)
	TimeInForceFOK TimeInForce = "fok" // FillOrKill
)

// OrderStatus is the order lifecycle state.
type OrderStatus string

const (
	OrderStatusOpen      OrderStatus = "open"
	OrderStatusClosed    OrderStatus = "closed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// StpAct is the self-trade-prevention action.
type StpAct string

const (
	StpActCancelNewest StpAct = "cn"
	StpActCancelOldest StpAct = "co"
	StpActCancelBoth   StpAct = "cb"
	StpActNone         StpAct = "-"
)
