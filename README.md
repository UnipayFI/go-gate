# go-gate

[![Go Reference](https://pkg.go.dev/badge/github.com/UnipayFI/go-gate/v4.svg)](https://pkg.go.dev/github.com/UnipayFI/go-gate/v4)
[![Go 1.27+](https://img.shields.io/badge/Go-1.27%2B-00ADD8?logo=go)](go.mod)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

A Go SDK for the [Gate.com](https://www.gate.com/docs/developers/apiv4/en/) (Gate.io) exchange, covering the entire APIv4 surface — every product line, REST and WebSocket.

| Area | API | Aligned to | Version |
|---|---|---|---|
| REST + WebSocket | `/api/v4` | 2026-09-17 | [v4.106.144](https://www.gate.com/docs/developers/apiv4/en/#changelog) |

Response structs are reconciled against the **live API** (not just the docs), so fields stay in sync — the SDK adds keys the official spec still omits (e.g. `rpi_maker_fee`, futures position vouchers, `market_cap`).

## Install

```bash
go get github.com/UnipayFI/go-gate/v4@latest
```

Requires Go 1.27 (see [JSON and timestamps](#json-and-timestamps)).

## Highlights

- One signing/transport core shared by every product; each product line is its own package with a dedicated client.
- Fluent per-endpoint API: `NewXxxService(...).SetFoo(...).Do(ctx)`.
- Amounts as `decimal.Decimal`, timestamps as `time.Time` whose wire format is declared per field with the `format` tag option (`json:"create_time,format:unix"`) — Gate's string- and number-encoded numbers, its heterogeneous second/millisecond timestamps, and `""` "not set" amounts are all decoded for you.
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

## JSON and timestamps

The SDK uses Go 1.27's `encoding/json/v2` (keep the default `jsonv2` GOEXPERIMENT enabled; without it the SDK does not
compile). Every `time.Time` field declares the wire format Gate actually sends with the `format` tag option: a bare
number of seconds is `json:"create_time,format:unix"`, a quoted one adds `,string`
(`json:"create_time,string,format:unix"`), milliseconds are `format:unixmilli`, and the few RFC 3339 strings (account
`tier_expire_time`, fixed-term extra-APR activity times) are `format:RFC3339`. Go 1.27 only honours `format` tags when
the experimental `ExperimentalSupportFormatTag` option is passed; `common.JSONMarshal` / `common.JSONUnmarshal`, which
every service and WebSocket subscription uses, pass it.
Encoding follows the tag exactly: the unix formats keep fractional seconds and milliseconds rather than truncating
them, and `format:RFC3339` writes whole seconds, as Gate does. Decoding applies the tag's standard semantics plus two
Gate quirks: a timestamp is accepted quoted or bare whatever `,string` says, and a
10/13/16/19-digit value is read as seconds/milliseconds/microseconds/nanoseconds even where the tag declares another
unit (such values would otherwise land in early 1970 or thousands of years ahead; Gate does send microseconds in a
seconds field when a futures price-triggered order is cancelled). `null` reads as the zero time; a `0` that Gate sends
for "not set" reads as the Unix epoch (`1970-01-01T00:00:00Z`), not the zero `time.Time`. Decoded times are in UTC —
use `.Equal` to compare and `.In(loc)` / `.Local()` to display.

A `time.Time` without a `format` option is RFC 3339, as in the standard library. That includes your own types passed
through `common.JSONMarshal` / `common.JSONUnmarshal` or `request.SetBody`.

To serialize SDK types yourself, use `common.JSONMarshal` / `common.JSONUnmarshal`. The one exception to round-tripping
is `spot.Candlestick`: Gate sends it as a positional array, which its `UnmarshalJSON` reads, but it marshals to an
object keyed by field name that `UnmarshalJSON` does not read back. Where you only need
`encoding/json/v2` to accept the tags (say, to store or log SDK values), you can pass
`github.com/go-json-experiment/json.ExperimentalSupportFormatTag(true)` yourself; that applies the standard semantics
only (quoting must match the tag, no unit detection, and `""` amounts are rejected), so it cannot reliably read Gate's
wire data. Plain `encoding/json` returns an error for structs with `format` tags (``unsupported `format` tag option``),
and `log/slog`'s JSON handler logs `!ERROR:...` in place of such a value.

## Packages

**Core**

| Package | Scope |
|---------|-------|
| `gate.go` | entry point: `NewSpotClient`/`NewFuturesClient`/… + WebSocket clients |
| `client/` `request/` | REST + WebSocket client, options, HMAC-SHA512 signer, response decode, subscribe framework |
| `common/` | constants, the `encoding/json/v2` codec: `format`-tagged `time.Time` fields, tolerant `decimal.Decimal` |

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
| `stock` | traditional-finance stock spot: symbols / details, order book, exchanges, fee rate, user assets, orders, positions, transactions, fund transfer |
| `p2p` | P2P merchant API: account & payment methods, ads, transactions, chat |
| `otc` | OTC fiat & stablecoin conversion + bank-card management |
| `bot` | strategy bots: spot / futures / margin / infinite grid, spot / contract martingale, portfolio management, AIHub recommendations |
| `announcement` | announcement articles: paginated list with title / tag / category / language / time filters |

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

- **2026-09-29** — Requires Go 1.27: moved to the standard library's `encoding/json/v2` with the experimental `format` tag support (see [JSON and timestamps](#json-and-timestamps)). Every time field's format was re-checked against the live API — public endpoints and WebSocket channels, plus the private endpoints the test account can use (some exercised with small, immediately reverted orders, lends and loans) — with the official docs only as leads. Fixed, each verified against live responses: portfolio-calculator `calculate_time` (seconds, previously read as milliseconds); uni-lending lends, lend records and interest records and margin-uni loans, loan records and interest records `create_time`/`update_time` (milliseconds, previously read as seconds); trail auto-orders (Gate quotes their IDs and times, which failed to decode); futures `get_leverage` (the key is `Lever`, so the leverage always read as zero); rebate agency/broker transaction and commission history and user info (single objects rather than arrays, which failed to decode; these services now return pointers); and futures WebSocket candlesticks (pushes are arrays, which failed to decode; `SubscribeCandlesticksService.Do` now takes a `WsHandler[[]WsFuturesCandlestick]`). Added the response fields Gate has since introduced: futures WS ticker `funding_*` / `change_*` / `price_type` / `t`, WS order-book `l`, spot WS trade `id_market` / `range`, delivery WS ticker `price_type`, index constituents `time` / `price`, options settlement voucher fields, currency-pair `delisting_time`, futures account `pnl_dividend`, TradFi ticker settlement currency and exchange rate, deposit `arrival_timestamp`, small-balance estimates, total-balance `total`, trail-order `default_leverage` and fixed-term `extra_apr_activity_info`. Decoding now also accepts a timestamp quoted or bare and reads 10/13/16/19-digit values by width. Spot candlestick timestamps are now UTC. `spot.Candlestick` and `account.AccountDetail` now carry `format` tags too, so like every other time-bearing type they need `common.JSONMarshal` rather than `encoding/json` or `log/slog`.
- **2026-07-08** — Aligned to v4.106.106. Added the `stock` package — Gate's traditional-finance stock spot module (`/api/v4/stock/*`, 16 endpoints): symbols & details, order book, supported exchanges, fee rate, user assets, orders (open/history/create/modify/cancel/cancel-all), positions & close, transactions & fund transfer. Public read structs reconciled against the live API (`stock/exchanges` requires signing despite being documented public).
- **2026-07-07** — Aligned to v4.106.105. Added five new product packages — `crossex` (cross-exchange margin & contracts), `tradfi` (MT5 stock/forex CFDs), `p2p` (P2P merchant API), `otc` (OTC fiat/stablecoin + bank cards) and `bot` (grid/martingale strategy bots) — plus extensions across existing products: futures trailing & chase auto-orders, BBO orders, split-mode leverage / position mode, `contracts_all`, batch funding rates and positions-timerange; earn auto-invest, fixed-term and dual/staking additions; unified delta-neutral & quick-repayment; rebate partner endpoints; `account/main_keys`; `wallet/getLowCapExchangeList`; options order amend. 150 endpoints added (415 official endpoints now fully covered). Public and account-reachable endpoints reconciled against the live API; capability-gated products (crossex/tradfi/p2p/otc) verified for endpoint + signing correctness.
- **2026-07-01** — Initial release. Full Gate APIv4 coverage: all REST products (spot, futures, delivery, options, margin, unified, wallet, account, sub-account, earn, loan, flash-swap, rebate) and spot/futures/delivery WebSocket public + private channels. Every public and private endpoint reconciled against the live API; order lifecycle (spot + futures, REST + WebSocket) verified with live trades.

## License

[MIT](LICENSE)
