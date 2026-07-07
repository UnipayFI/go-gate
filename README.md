# go-gate

[![Go Reference](https://pkg.go.dev/badge/github.com/UnipayFI/go-gate/v4.svg)](https://pkg.go.dev/github.com/UnipayFI/go-gate/v4)
[![Go 1.26+](https://img.shields.io/badge/Go-1.26%2B-00ADD8?logo=go)](go.mod)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

A Go SDK for the [Gate.com](https://www.gate.com/docs/developers/apiv4/en/) (Gate.io) exchange, covering the entire APIv4 surface — every product line, REST and WebSocket.

| Area | API | Aligned to | Version |
|---|---|---|---|
| REST + WebSocket | `/api/v4` | 2026-07-06 | [v4.106.105](https://www.gate.com/docs/developers/apiv4/en/#changelog) |

Response structs are reconciled against the **live API** (not just the docs), so fields stay in sync — the SDK adds keys the official spec still omits (e.g. `rpi_maker_fee`, futures position vouchers, `market_cap`).

## Install

```bash
go get github.com/UnipayFI/go-gate/v4@latest
```

## Highlights

- One signing/transport core shared by every product; each product line is its own package with a dedicated client.
- Fluent per-endpoint API: `NewXxxService(...).SetFoo(...).Do(ctx)`.
- Amounts as `decimal.Decimal`, timestamps as `time.Time` — Gate's string- and number-encoded numbers, its heterogeneous second/millisecond timestamps, and `""`/`"0"` "not set" sentinels are all decoded for you.
- Every endpoint is tested against the live API, diffing real JSON keys against the struct.

## Quick start

```go
package main

import (
	"context"
	"fmt"

	gate "github.com/UnipayFI/go-gate/v4"
	"github.com/UnipayFI/go-gate/v4/client"
	"github.com/UnipayFI/go-gate/v4/futures"
	"github.com/UnipayFI/go-gate/v4/spot"
	"github.com/shopspring/decimal"
)

func main() {
	ctx := context.Background()

	c := gate.NewSpotClient(
		client.WithAuth("apiKey", "apiSecret"),
		// client.WithProxy("socks5://127.0.0.1:7890"),
	)
	_ = c.SyncServerTime(ctx) // align clock to avoid signature drift

	// Public market data (no auth).
	pair, _ := c.NewGetCurrencyPairService("BTC_USDT").Do(ctx)
	fmt.Println(pair.ID, pair.Precision, pair.MinQuoteAmount)

	// Private account data.
	accounts, _ := c.NewListSpotAccountsService().Do(ctx)
	for _, a := range accounts {
		fmt.Println(a.Currency, a.Available)
	}

	// Place a limit order.
	order, err := c.NewCreateOrderService("BTC_USDT", spot.SideBuy, decimal.RequireFromString("0.0001")).
		SetType(spot.OrderTypeLimit).
		SetPrice(decimal.RequireFromString("30000")).
		SetTimeInForce(spot.TimeInForceGTC).
		Do(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("orderId:", order.ID)

	// Futures — settle currency chosen per request.
	f := gate.NewFuturesClient(client.WithAuth("apiKey", "apiSecret"))
	acct, _ := f.NewListFuturesAccountsService(futures.SettleUSDT).Do(ctx)
	fmt.Println("futures total:", acct.Total)
}
```

## Authentication

Pass credentials from the Gate API-keys page (Gate uses no passphrase):

```go
c := gate.NewSpotClient(client.WithAuth(apiKey, apiSecret))
```

Requests are signed per Gate's APIv4 scheme: `SIGN = hex(HMAC-SHA512(secret, prehash))`, where

```
prehash = METHOD + "\n" + path + "\n" + query + "\n" + hex(SHA512(body)) + "\n" + timestamp
```

is sent in the `KEY` / `SIGN` / `Timestamp` headers (`Timestamp` in whole seconds). For an external signer, pass `client.WithSignFn(fn)`.

Other options: `WithProxy` (http/https/socks5), `WithBaseURL`, `WithNetwork`, `WithTimeOffset`, `WithLogger`, `WithHTTPClient`.

## WebSocket

```go
ws := gate.NewSpotWebSocketClient(
	client.WithWebSocketAuth(apiKey, apiSecret), // private channels only
)

// Public best bid/ask.
done, _, _ := ws.NewSubscribeBookTickerService("BTC_USDT").
	Do(ctx, func(p *request.WsPush[spot.WsBookTicker], err error) {
		if err != nil {
			return
		}
		fmt.Println(p.Result.CurrencyPair, p.Result.BestBid, p.Result.BestAsk)
	})
close(done) // unsubscribe + close

// Private orders (auth attached per subscription).
ws.NewSubscribeOrdersService("!all").Do(ctx, func(p *request.WsPush[[]spot.WsOrder], err error) {
	// p.Result[0].ID, p.Result[0].Event, ...
})
```

Each `Do` returns `(done chan<- struct{}, stop <-chan struct{}, err error)`: close `done` to unsubscribe; `stop` closes when the reader exits. Ping keepalive is automatic. Futures streams take a settle currency: `gate.NewFuturesWebSocketClient(futures.SettleUSDT, ...)`.

## Packages

**Core**

| Package | Scope |
|---------|-------|
| `gate.go` | entry point: `NewSpotClient`/`NewFuturesClient`/… + WebSocket clients |
| `client/` `request/` | REST + WebSocket client, options, HMAC-SHA512 signer, response decode, subscribe framework |
| `common/` | constants, global tolerant `decimal.Decimal` JSON codec |

**Products**

| Package | Scope |
|---------|-------|
| `spot` | currencies, currency-pairs, tickers, order book, trades, candlesticks, fee, accounts, account-book, orders (single/batch/amend/cancel/countdown), my-trades, price-triggered orders + spot WebSocket |
| `futures` | contracts (+delisted), order book, trades, candlesticks, premium index, tickers, funding rate (single + batch), insurance, contract stats, index constituents, liq orders, risk-limit tiers, accounts, positions (single + dual-mode + history + split-mode leverage), position mode, orders (single/batch/amend/price-triggered/BBO), trailing & chase auto-orders, my-trades + futures WebSocket |
| `delivery` | dated-futures market, accounts, positions, orders, settlements, risk-limit tiers, price-triggered orders + delivery WebSocket |
| `options` | underlyings, expirations, contracts, settlements, order book, tickers, candlesticks, trades, accounts, positions, orders, MMP |
| `margin` | isolated + cross margin, funding accounts, auto-repay, margin tiers, unified-margin (uni) lending |
| `unified` | unified account, borrow/repay, quick-repayment, transferables, risk units, mode, delta-neutral mode, leverage config, discount tiers, portfolio calculator |
| `wallet` | deposit address, transfers, sub-account transfers, deposits/withdrawals, balances, trade fee, total balance, dust conversion, **withdrawals** |
| `account` | account detail, rate limit, STP groups, debit fee, main-account API keys |
| `subaccount` | sub-account create/query, API keys, lock/unlock |
| `earn` | dual investment (+ balance / refund / reinvest), structured products, staking (ETH2 + on-chain assets / awards / orders), auto-invest plans, fixed-term products & subscriptions, uni lending |
| `loan` | collateral loan + multi-collateral loan |
| `flashswap` | flash-swap currency pairs, preview, orders |
| `rebate` | agency / partner / broker commission & transaction history, partner applications / eligibility |
| `crossex` | cross-exchange margin & contract trading: accounts, orders, positions & leverage, transfers, convert (flash-swap), rules & fees |
| `tradfi` | TradFi CFDs via MT5: symbols / categories / klines / tickers, orders, positions, user & MT5 account, fund transactions |
| `p2p` | P2P merchant API: account & payment methods, ads, transactions, chat |
| `otc` | OTC fiat & stablecoin conversion + bank-card management |
| `bot` | strategy bots: spot / futures / margin / infinite grid, spot / contract martingale, portfolio management, AIHub recommendations |

## Testing

Tests hit the live API and read credentials from the environment, skipping when unset:

```bash
export GATE_API_KEY=...  GATE_API_SECRET=...
export GATE_PROXY=socks5://127.0.0.1:7890   # optional

go test ./spot/ -run TestSpotAccounts -v             # one module at a time
GATE_TEST_WRITE=1 go test ./spot/ -run TestSpotOrder  # live order tests (tiny, reversible)
```

- Run **per module** (`-run TestXxx`) — Gate rate-limits per endpoint, so the full suite can trip HTTP 429.
- Capability-gated reads (unified account, agency/broker rebate, options, delivery balance) are skipped when the account lacks the capability — signing is still exercised.
- State-changing tests are gated behind `GATE_TEST_WRITE=1` (minimal amounts, large-cap symbols, place → query → cancel). Withdrawals are implemented but **never executed**.

## CHANGE_LOG

- **2026-07-07** — Aligned to v4.106.105. Added five new product packages — `crossex` (cross-exchange margin & contracts), `tradfi` (MT5 stock/forex CFDs), `p2p` (P2P merchant API), `otc` (OTC fiat/stablecoin + bank cards) and `bot` (grid/martingale strategy bots) — plus extensions across existing products: futures trailing & chase auto-orders, BBO orders, split-mode leverage / position mode, `contracts_all`, batch funding rates and positions-timerange; earn auto-invest, fixed-term and dual/staking additions; unified delta-neutral & quick-repayment; rebate partner endpoints; `account/main_keys`; `wallet/getLowCapExchangeList`; options order amend. 150 endpoints added (415 official endpoints now fully covered). Public and account-reachable endpoints reconciled against the live API; capability-gated products (crossex/tradfi/p2p/otc) verified for endpoint + signing correctness.
- **2026-07-01** — Initial release. Full Gate APIv4 coverage: all REST products (spot, futures, delivery, options, margin, unified, wallet, account, sub-account, earn, loan, flash-swap, rebate) and spot/futures/delivery WebSocket public + private channels. Every public and private endpoint reconciled against the live API; order lifecycle (spot + futures, REST + WebSocket) verified with live trades.

## License

[MIT](LICENSE)
