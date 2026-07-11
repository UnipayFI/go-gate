// Package gate is the entry point of the Gate.com (Gate.io) exchange Go SDK.
//
// Install: go get github.com/UnipayFI/go-gate/v4
// Import:  import gate "github.com/UnipayFI/go-gate/v4"
//
// The SDK covers Gate's v4 REST API and v4 WebSocket streams across every
// product line, each exposed as its own package with a dedicated client:
//
//   - spot        — /api/v4/spot/*      (spot, margin & unified order flow)
//   - futures     — /api/v4/futures/*   (USDT- & BTC-settled perpetuals)
//   - delivery    — /api/v4/delivery/*  (dated futures)
//   - options     — /api/v4/options/*
//   - margin      — /api/v4/margin/*    (isolated + unified-margin lending)
//   - unified     — /api/v4/unified/*   (unified account)
//   - wallet      — /api/v4/wallet/* + /api/v4/withdrawals/*
//   - account     — /api/v4/account/*
//   - subaccount  — /api/v4/sub_accounts/*
//   - earn        — /api/v4/earn/*
//   - loan        — /api/v4/loan/*      (collateral + multi-collateral)
//   - flashswap   — /api/v4/flash_swap/*
//   - rebate      — /api/v4/rebate/*
//   - crossex     — /api/v4/crossex/*   (cross-exchange margin & contracts)
//   - tradfi      — /api/v4/tradfi/*    (MT5 stock / forex CFDs)
//   - stock       — /api/v4/stock/*     (traditional-finance stock spot)
//   - p2p         — /api/v4/p2p/*       (P2P merchant API)
//   - otc         — /api/v4/otc/*       (OTC fiat / stablecoin conversion)
//   - bot         — /api/v4/bot/*       (grid / martingale strategy bots)
//
// Authentication uses Gate's APIv4 scheme (KEY / SIGN / Timestamp headers,
// HMAC-SHA512 over method\npath\nquery\nSHA512(body)\ntimestamp); the shared
// client/request/common layers are reused by every product.
//
// Quick start (spot):
//
//	c := gate.NewSpotClient(client.WithAuth(apiKey, apiSecret))
//	if err := c.SyncServerTime(ctx); err != nil { /* ... */ }
//	accounts, err := c.NewListSpotAccountsService().Do(ctx)
//
// Quick start (futures):
//
//	f := gate.NewFuturesClient(client.WithAuth(apiKey, apiSecret))
//	acct, err := f.NewListFuturesAccountsService(futures.SettleUSDT).Do(ctx)
package gate

import (
	"github.com/UnipayFI/go-gate/v4/account"
	"github.com/UnipayFI/go-gate/v4/bot"
	"github.com/UnipayFI/go-gate/v4/client"
	"github.com/UnipayFI/go-gate/v4/crossex"
	"github.com/UnipayFI/go-gate/v4/delivery"
	"github.com/UnipayFI/go-gate/v4/earn"
	"github.com/UnipayFI/go-gate/v4/flashswap"
	"github.com/UnipayFI/go-gate/v4/futures"
	"github.com/UnipayFI/go-gate/v4/loan"
	"github.com/UnipayFI/go-gate/v4/margin"
	"github.com/UnipayFI/go-gate/v4/options"
	"github.com/UnipayFI/go-gate/v4/otc"
	"github.com/UnipayFI/go-gate/v4/p2p"
	"github.com/UnipayFI/go-gate/v4/rebate"
	"github.com/UnipayFI/go-gate/v4/spot"
	"github.com/UnipayFI/go-gate/v4/stock"
	"github.com/UnipayFI/go-gate/v4/subaccount"
	"github.com/UnipayFI/go-gate/v4/tradfi"
	"github.com/UnipayFI/go-gate/v4/unified"
	"github.com/UnipayFI/go-gate/v4/wallet"
)

// --- REST clients ---

// NewSpotClient constructs a spot / margin / unified REST client.
func NewSpotClient(opts ...client.Options) *spot.SpotClient {
	return spot.NewSpotClient(opts...)
}

// NewFuturesClient constructs a perpetual-futures REST client (settle chosen
// per request).
func NewFuturesClient(opts ...client.Options) *futures.FuturesClient {
	return futures.NewFuturesClient(opts...)
}

// NewDeliveryClient constructs a delivery (dated-futures) REST client.
func NewDeliveryClient(opts ...client.Options) *delivery.DeliveryClient {
	return delivery.NewDeliveryClient(opts...)
}

// NewOptionsClient constructs an options REST client.
func NewOptionsClient(opts ...client.Options) *options.OptionsClient {
	return options.NewOptionsClient(opts...)
}

// NewMarginClient constructs a margin / unified-margin REST client.
func NewMarginClient(opts ...client.Options) *margin.MarginClient {
	return margin.NewMarginClient(opts...)
}

// NewUnifiedClient constructs a unified-account REST client.
func NewUnifiedClient(opts ...client.Options) *unified.UnifiedClient {
	return unified.NewUnifiedClient(opts...)
}

// NewWalletClient constructs a wallet / withdrawal REST client.
func NewWalletClient(opts ...client.Options) *wallet.WalletClient {
	return wallet.NewWalletClient(opts...)
}

// NewAccountClient constructs an account REST client.
func NewAccountClient(opts ...client.Options) *account.AccountClient {
	return account.NewAccountClient(opts...)
}

// NewSubAccountClient constructs a sub-account REST client.
func NewSubAccountClient(opts ...client.Options) *subaccount.SubAccountClient {
	return subaccount.NewSubAccountClient(opts...)
}

// NewEarnClient constructs an earn REST client.
func NewEarnClient(opts ...client.Options) *earn.EarnClient {
	return earn.NewEarnClient(opts...)
}

// NewLoanClient constructs a collateral-loan REST client.
func NewLoanClient(opts ...client.Options) *loan.LoanClient {
	return loan.NewLoanClient(opts...)
}

// NewFlashSwapClient constructs a flash-swap REST client.
func NewFlashSwapClient(opts ...client.Options) *flashswap.FlashSwapClient {
	return flashswap.NewFlashSwapClient(opts...)
}

// NewRebateClient constructs a rebate REST client.
func NewRebateClient(opts ...client.Options) *rebate.RebateClient {
	return rebate.NewRebateClient(opts...)
}

// NewCrossexClient constructs a Cross-Exchange (cross-venue margin & contract)
// REST client.
func NewCrossexClient(opts ...client.Options) *crossex.CrossexClient {
	return crossex.NewCrossexClient(opts...)
}

// NewTradfiClient constructs a TradFi (MT5-backed stock / forex CFD) REST client.
func NewTradfiClient(opts ...client.Options) *tradfi.TradfiClient {
	return tradfi.NewTradfiClient(opts...)
}

// NewStockClient constructs a Stock (traditional-finance stock spot) REST client.
func NewStockClient(opts ...client.Options) *stock.StockClient {
	return stock.NewStockClient(opts...)
}

// NewP2PClient constructs a P2P merchant REST client.
func NewP2PClient(opts ...client.Options) *p2p.P2PClient {
	return p2p.NewP2PClient(opts...)
}

// NewOTCClient constructs an OTC fiat / stablecoin conversion REST client.
func NewOTCClient(opts ...client.Options) *otc.OTCClient {
	return otc.NewOTCClient(opts...)
}

// NewBotClient constructs a strategy-bot (grid / martingale) REST client.
func NewBotClient(opts ...client.Options) *bot.BotClient {
	return bot.NewBotClient(opts...)
}

// --- WebSocket clients ---

// NewSpotWebSocketClient constructs a spot / margin / unified stream client.
func NewSpotWebSocketClient(opts ...client.WebSocketOptions) *spot.SpotWebSocketClient {
	return spot.NewSpotWebSocketClient(opts...)
}

// NewFuturesWebSocketClient constructs a perpetual-futures stream client for the
// given settle currency (futures.SettleUSDT or futures.SettleBTC).
func NewFuturesWebSocketClient(settle futures.Settle, opts ...client.WebSocketOptions) *futures.FuturesWebSocketClient {
	return futures.NewFuturesWebSocketClient(settle, opts...)
}

// NewDeliveryWebSocketClient constructs a delivery stream client.
func NewDeliveryWebSocketClient(opts ...client.WebSocketOptions) *delivery.DeliveryWebSocketClient {
	return delivery.NewDeliveryWebSocketClient(opts...)
}
