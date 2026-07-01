package delivery

// Settle is the settlement currency of a delivery contract, carried as the
// {settle} path segment. Delivery settles in USDT.
type Settle string

const (
	SettleUSDT Settle = "usdt"
)

// OrderStatus is the delivery order lifecycle state.
type OrderStatus string

const (
	OrderStatusOpen     OrderStatus = "open"
	OrderStatusFinished OrderStatus = "finished"
)

// FinishAs explains how a finished order concluded.
type FinishAs string

const (
	FinishAsFilled         FinishAs = "filled"
	FinishAsCancelled      FinishAs = "cancelled"
	FinishAsLiquidated     FinishAs = "liquidated"
	FinishAsIOC            FinishAs = "ioc"
	FinishAsReduceOnly     FinishAs = "reduce_only"
	FinishAsPositionClosed FinishAs = "position_closed"
	FinishAsStp            FinishAs = "stp"
)

// TimeInForce determines how long a delivery order stays active.
type TimeInForce string

const (
	TimeInForceGTC TimeInForce = "gtc"
	TimeInForceIOC TimeInForce = "ioc"
	TimeInForcePOC TimeInForce = "poc"
	TimeInForceFOK TimeInForce = "fok"
)

// StpAct is the self-trade-prevention action.
type StpAct string

const (
	StpActCancelNewest StpAct = "cn"
	StpActCancelOldest StpAct = "co"
	StpActCancelBoth   StpAct = "cb"
	StpActNone         StpAct = "-"
)
