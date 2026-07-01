package futures

// Settle is the settlement currency of a futures contract, carried as the
// {settle} path segment.
type Settle string

const (
	SettleUSDT Settle = "usdt"
	SettleBTC  Settle = "btc"
	SettleUSD  Settle = "usd"
)

// OrderStatus is the futures order lifecycle state.
type OrderStatus string

const (
	OrderStatusOpen     OrderStatus = "open"
	OrderStatusFinished OrderStatus = "finished"
)

// FinishAs explains how a finished order concluded.
type FinishAs string

const (
	FinishAsFilled          FinishAs = "filled"
	FinishAsCancelled       FinishAs = "cancelled"
	FinishAsLiquidated      FinishAs = "liquidated"
	FinishAsIOC             FinishAs = "ioc"
	FinishAsAutoDeleveraged FinishAs = "auto_deleveraged"
	FinishAsReduceOnly      FinishAs = "reduce_only"
	FinishAsPositionClosed  FinishAs = "position_closed"
	FinishAsReduceOut       FinishAs = "reduce_out"
	FinishAsSTP             FinishAs = "stp"
)

// TimeInForce determines how long a futures order stays active.
type TimeInForce string

const (
	TimeInForceGTC TimeInForce = "gtc" // GoodTillCancelled
	TimeInForceIOC TimeInForce = "ioc" // ImmediateOrCancelled (taker only)
	TimeInForcePOC TimeInForce = "poc" // PendingOrCancelled (post-only)
	TimeInForceFOK TimeInForce = "fok" // FillOrKill
)

// AutoSize closes one leg of a dual-mode position (size must be 0).
type AutoSize string

const (
	AutoSizeCloseLong  AutoSize = "close_long"
	AutoSizeCloseShort AutoSize = "close_short"
)

// StpAct is the self-trade-prevention action.
type StpAct string

const (
	StpActCancelNewest StpAct = "cn"
	StpActCancelOldest StpAct = "co"
	StpActCancelBoth   StpAct = "cb"
	StpActNone         StpAct = "-"
)
