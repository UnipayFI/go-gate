export const meta = {
  name: 'impl-gate-rest',
  description: 'Implement a Gate API v4 REST package from the official gateapi-go spec + live responses, in go-bitget house style',
  phases: [{ title: 'Implement', detail: 'one agent per file-group; official spec + live-verified' }],
}

// args: {
//   pkg, dir, refDir, clientType, apiFile,
//   sharedTypesNote,          // string: enums/types already defined in types.go + who-owns-what
//   testHelpersNote,          // string: which _test.go defines testClient/testPublicClient
//   groups: [{ file, label, clientCtor, endpoints, typesOwned, testHints }]
// }
const A = typeof args === 'string' ? JSON.parse(args) : args
const DIR = A.dir
const REF = A.refDir

const STYLE = `You are adding endpoints to an existing Go SDK for the Gate.com (Gate API v4) exchange. Module: github.com/UnipayFI/go-gate. Package: ${A.pkg}. Working dir: ${DIR}. Client type: ${A.clientType}.

Your code MUST look like it was written by the same author as the rest of the SDK. BEFORE writing, Read these house-style anchors in full:
- ${DIR}/spot/market.go       — the canonical Service pattern (public GET): one Service struct per endpoint, NewXxxService constructor on the client, chainable SetFoo setters for OPTIONAL params, Do(ctx) calling request.Get/Do[T].
- ${DIR}/spot/account.go      — a PRIVATE GET Service using .WithSign().
- ${DIR}/spot/client.go       — the client type + SyncServerTime.
- ${DIR}/spot/types.go        — how enums are declared.
- ${DIR}/spot/spot_test.go    — the live-test pattern with testutil helpers + assertCovers.
- ${DIR}/request/request.go   — request.Get/Post/Put/Delete, params map, .WithSign(); ${DIR}/request/send.go — request.Do[T]/DoRaw (Gate has NO envelope: 2xx body IS the payload).
- ${DIR}/common/json.go       — global tolerant DECIMAL codec; time is NOT global (per-field format tags).

GROUND TRUTH (two sources, reconcile BOTH):
1. Official Gate SDK (authoritative for PATH, HTTP method, path/query/body PARAMS, and field names/types): read the assigned methods in ${REF}/${A.apiFile} and the model_*.go structs they return (grep for the return type, e.g. \`grep -n "func (a \\*XxxApiService) MethodName" ${REF}/${A.apiFile}\` then read that method for the path + params, and \`cat ${REF}/model_<snake>.go\` for the response struct). The official docs markdown is in ${REF}/docs/*.md if you need semantics.
2. LIVE API (authoritative for FIELDS — the official model can lag the live response): curl every PUBLIC endpoint with /usr/bin/curl against https://api.gateio.ws (e.g. \`/usr/bin/curl -sS 'https://api.gateio.ws/api/v4/spot/trades?currency_pair=BTC_USDT&limit=1'\`). Add EVERY key the live response contains, even if the official model omits it. For PARAMS trust the official spec; where doc and live differ on FIELDS trust live.

HOUSE STYLE — follow precisely:
- One Service type per endpoint:
    // XxxService -- <HTTP METHOD> <full /api/v4 path> (<"private" if signed>)
    //
    // <one-line purpose>.
    type XxxService struct { c *${A.clientType}; params map[string]string }   // or body map[string]any for POST/PUT
    func (c *${A.clientType}) NewXxxService(<REQUIRED args>) *XxxService { return &XxxService{c: c, params: map[string]string{...}} }
    func (s *XxxService) SetOpt(v T) *XxxService { s.params["opt"] = ...; return s }   // OPTIONAL params only
    func (s *XxxService) Do(ctx context.Context) (*Xxx, error) { req := request.Get(ctx, s.c, "<path>", s.params); return request.Do[Xxx](req) }
  - REQUIRED params (incl. path params like {currency_pair}, {order_id}, {settle}) -> constructor arguments. Path params are formatted into the path string (e.g. "/api/v4/spot/orders/"+orderID). OPTIONAL params -> chainable SetXxx setters.
  - Endpoint returns a JSON ARRAY -> Do returns ([]Xxx, error): resp, err := request.Do[[]Xxx](req); if err != nil { return nil, err }; return *resp, nil
  - PRIVATE endpoints (anything requiring auth: accounts, orders, trades, positions, wallet, etc.) MUST chain .WithSign(). PUBLIC market endpoints must NOT.
  - POST/PUT/DELETE with a body: build body map[string]any; required -> constructor, optional -> setters; for array/nested bodies use req.SetBody(v). DELETE with only query -> request.Delete(ctx, c, path, params).
  - Full paths always include the /api/v4 prefix, e.g. "/api/v4/spot/orders".
  - Use strconv for numeric setters: s.params["limit"] = strconv.Itoa(limit); unix-second times strconv.FormatInt(t.Unix(), 10).

TYPE MAPPING (map by the OFFICIAL Go type = the wire encoding, and the field SEMANTICS):
- Monetary/quantity/price/rate/ratio/fee/pnl/volume values (official type string OR float) -> github.com/shopspring/decimal Decimal, plain tag \`json:"key"\`. The global codec tolerates ""/null/number/string.
- TIMESTAMPS -> time.Time with a FORMAT TAG chosen by (a) wire encoding and (b) unit:
    * official type string, name WITHOUT _ms (e.g. create_time) -> \`json:"create_time,string,format:unix"\`
    * official type string, name WITH _ms          -> \`json:"create_time_ms,string,format:unixmilli"\`
    * official type int/int64/float, seconds (e.g. futures create_time, or *_time in seconds) -> \`json:"create_time,format:unix"\`
    * official type int/int64, milliseconds (e.g. server_time, current, update, *_ms as number) -> \`json:"...,format:unixmilli"\`
  Decide sec vs ms from the field name/description ("in milliseconds" or _ms => unixmilli; else unix seconds). Only apply this to fields that are SEMANTICALLY a moment in time — NOT ids/versions/counts/durations.
- IDs: match the official type EXACTLY — spot order/trade id is string; futures/delivery order id is int64; sequence_id/trade_id may be string. Getting number-vs-string wrong makes the live decode FAIL.
- Precision / count / *_num / level / interval-seconds / user_id -> int or int64 (Gate sends BARE integers here; do NOT use string).
- Enums (side, type, status, time_in_force, account, ...) -> the typed const from ${A.pkg}/types.go when it exists; else string.
- booleans -> bool.
- RESPONSE structs: Go-exported CamelCase names expanding abbreviations per ~/Downloads/go-okx-字段命名说明.md (Id->ID, Url->URL, Api->API, Usdt->USDT, Pnl->PnL where it reads well), but the json tag MUST equal the real snake_case key EXACTLY. Do NOT put ,omitempty on response-struct fields (it breaks the assertCovers round-trip). Nested objects -> nested named structs.

${A.sharedTypesNote}

OUTPUT:
- Write ONE self-contained Go file: ${'`'}${DIR}/${'${g.file}'}${'`'} with \`package ${A.pkg}\`, correct imports (only what you use — unused imports break the build), all Services + response structs for YOUR assigned endpoints.
- ALSO write a companion test file (same name, _test.go) in package ${A.pkg} with ONE func Test<Label>(t *testing.T) that live-tests EVERY endpoint you implemented, using the existing helpers ${A.testHelpersNote} plus testutil.Ctx(t), testutil.FetchRawGet/FetchRawPost(t,c,cx,path,params,sign), testutil.AssertCovers(t,label,raw,resp), testutil.Tolerable(t,label,err), testutil.WriteEnabled(). Rules:
    * PUBLIC endpoint: call it, log the result, fetch raw and AssertCovers. Use realistic params (currency_pair=BTC_USDT, settle=usdt, contract=BTC_USDT, limit=2).
    * PRIVATE READ: call signed; if it errors, treat testutil.Tolerable as a pass (log + return); else AssertCovers. Never t.Fatal on a tolerable capability/empty error.
    * STATE-CHANGING (place/amend/cancel order, transfer): guard the whole block behind if !testutil.WriteEnabled() { t.Skip(...) }; use TINY amounts on BTC_USDT (spot: ~5 USDT notional min; qty small), place -> query -> cancel so it is reversible. WITHDRAWAL endpoints: implement but DO NOT call in tests.
    * Read creds via the existing testClient(t) (skips when unset).
- Run \`cd ${DIR} && gofmt -w <yourfile> <yourtestfile>\` (do NOT run go build/go test — sibling agents write other files in this package concurrently).
- Edit ONLY your two assigned files. Do not touch client.go, types.go, or any sibling's file.

Return a structured report.`

const SCHEMA = {
  type: 'object',
  additionalProperties: false,
  required: ['file', 'endpoints', 'issues'],
  properties: {
    file: { type: 'string' },
    endpoints: {
      type: 'array',
      items: {
        type: 'object',
        additionalProperties: false,
        required: ['name', 'method', 'path', 'private'],
        properties: {
          name: { type: 'string' },
          method: { type: 'string' },
          path: { type: 'string' },
          private: { type: 'boolean' },
          note: { type: 'string' },
        },
      },
    },
    typesDefined: { type: 'array', items: { type: 'string' }, description: 'exported response struct/type names defined in this file' },
    issues: { type: 'string', description: 'doc/live mismatches, ambiguous fields, endpoints not curlable, or types referenced-but-owned-elsewhere' },
  },
}

phase('Implement')
const reports = await parallel(A.groups.map((g) => () =>
  agent(
    `${STYLE}\n\n=== YOUR ASSIGNMENT ===\nFile: ${DIR}/${g.file}\nCompanion test: same path with _test.go, func Test${g.label}.\nClient constructor in tests: ${g.clientCtor || A.testHelpersNote}\nYou OWN and must DEFINE these response types (do not let another file define them): ${(g.typesOwned || []).join(', ') || '(name them yourself)'}\nEndpoints to implement (read each method in ${REF}/${A.apiFile} for exact path/params, then the returned model_*.go, then curl live where public):\n${g.endpoints}\n${g.testHints ? '\nTest hints: ' + g.testHints : ''}`,
    { label: g.label, phase: 'Implement', schema: SCHEMA }
  )
))

return reports.filter(Boolean)
