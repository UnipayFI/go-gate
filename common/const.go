package common

import "time"

const (
	GO_GATE_USER_AGENT = "go-gate/1.0"

	// REST base host. Gate serves every business line (spot, margin, unified,
	// perpetual futures, delivery, options, earn, ...) under the /api/v4/* path
	// prefix on a single host; the product is encoded in the path, not the host.
	// The full path (including /api/v4) is what gets signed.
	DEFAULT_REST_BASE_URL = "https://api.gateio.ws"

	// Spot / margin / unified-account WebSocket gateway (v4). One connection
	// multiplexes every spot-family channel; the product is the channel prefix
	// (spot.tickers, spot.orders, ...).
	DEFAULT_WS_SPOT_URL = "wss://api.gateio.ws/ws/v4/"

	// Perpetual-futures WebSocket gateways — one per settle currency. Channels
	// are prefixed "futures." (futures.tickers, futures.orders, ...).
	DEFAULT_WS_FUTURES_USDT_URL = "wss://fx-ws.gateio.ws/v4/ws/usdt"
	DEFAULT_WS_FUTURES_BTC_URL  = "wss://fx-ws.gateio.ws/v4/ws/btc"

	// Delivery (dated-futures) WebSocket gateway. Channels are prefixed
	// "futures." as well but flow over the delivery endpoint.
	DEFAULT_WS_DELIVERY_USDT_URL = "wss://fx-ws.gateio.ws/v4/ws/delivery/usdt"

	// Options WebSocket gateway. Channels are prefixed "options.".
	DEFAULT_WS_OPTIONS_URL = "wss://op-ws.gateio.ws/v4/ws/usdt"

	DEFAULT_KEEP_ALIVE_INTERVAL = 15 * time.Second
	DEFAULT_KEEP_ALIVE_TIMEOUT  = 60 * time.Second
)

// Network identifies which Gate environment a client targets. Gate exposes a
// dedicated futures/delivery testnet (fx-api-testnet.gateio.ws) but has no spot
// testnet; spot clients ignore Testnet and stay on the production host. Use
// client.WithBaseURL to point at any other host.
type Network int

const (
	Mainnet Network = iota
	Testnet
)

// RestBaseURL returns the REST base host for this network.
func (n Network) RestBaseURL() string {
	switch n {
	case Testnet:
		return "https://fx-api-testnet.gateio.ws"
	default:
		return DEFAULT_REST_BASE_URL
	}
}
