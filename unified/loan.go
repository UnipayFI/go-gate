package unified

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// ListUnifiedLoansService -- GET /api/v4/unified/loans (private)
//
// Returns the outstanding borrowings of the authenticated unified account,
// optionally filtered by currency and loan type.
type ListUnifiedLoansService struct {
	c      *UnifiedClient
	params map[string]string
}

func (c *UnifiedClient) NewListUnifiedLoansService() *ListUnifiedLoansService {
	return &ListUnifiedLoansService{c: c, params: map[string]string{}}
}

// SetCurrency narrows the result to a single currency (e.g. USDT).
func (s *ListUnifiedLoansService) SetCurrency(currency string) *ListUnifiedLoansService {
	s.params["currency"] = currency
	return s
}

// SetPage selects the page number.
func (s *ListUnifiedLoansService) SetPage(page int) *ListUnifiedLoansService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit caps the number of records returned (default 100, max 100).
func (s *ListUnifiedLoansService) SetLimit(limit int) *ListUnifiedLoansService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetType filters by loan type: "platform" (platform borrowing) or "margin"
// (margin borrowing).
func (s *ListUnifiedLoansService) SetType(loanType string) *ListUnifiedLoansService {
	s.params["type"] = loanType
	return s
}

func (s *ListUnifiedLoansService) Do(ctx context.Context) ([]UnifiedLoan, error) {
	req := request.Get(ctx, s.c, "/api/v4/unified/loans", s.params).WithSign()
	resp, err := request.Do[[]UnifiedLoan](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// UnifiedLoan is a single outstanding borrowing of the unified account.
type UnifiedLoan struct {
	Currency     string          `json:"currency"`
	CurrencyPair string          `json:"currency_pair"`
	Amount       decimal.Decimal `json:"amount"`
	Type         string          `json:"type"`
	CreateTime   time.Time       `json:"create_time,format:unix"`
	UpdateTime   time.Time       `json:"update_time,format:unix"`
}

// CreateUnifiedLoanService -- POST /api/v4/unified/loans (private)
//
// Borrows or repays a currency on the unified account. When borrowing, the
// amount must sit between the currency's minimum and the platform/user borrow
// limit; for repayment set repaid_all to clear the full outstanding balance.
type CreateUnifiedLoanService struct {
	c    *UnifiedClient
	body map[string]any
}

func (c *UnifiedClient) NewCreateUnifiedLoanService(currency string, amount decimal.Decimal, loanType string) *CreateUnifiedLoanService {
	return &CreateUnifiedLoanService{c: c, body: map[string]any{
		"currency": currency,
		"amount":   amount.String(),
		"type":     loanType,
	}}
}

// SetRepaidAll, when true, repays the full outstanding amount and overrides the
// amount field. Only used for repayment operations.
func (s *CreateUnifiedLoanService) SetRepaidAll(repaidAll bool) *CreateUnifiedLoanService {
	s.body["repaid_all"] = repaidAll
	return s
}

// SetText attaches a user-defined custom ID.
func (s *CreateUnifiedLoanService) SetText(text string) *CreateUnifiedLoanService {
	s.body["text"] = text
	return s
}

func (s *CreateUnifiedLoanService) Do(ctx context.Context) (*UnifiedLoanResult, error) {
	req := request.Post(ctx, s.c, "/api/v4/unified/loans", s.body).WithSign()
	return request.Do[UnifiedLoanResult](req)
}

// UnifiedLoanResult is the borrow/repay acknowledgement of the unified account.
type UnifiedLoanResult struct {
	TranID int64 `json:"tran_id"`
}

// ListUnifiedLoanRecordsService -- GET /api/v4/unified/loan_records (private)
//
// Returns the borrow/repay history of the authenticated unified account.
type ListUnifiedLoanRecordsService struct {
	c      *UnifiedClient
	params map[string]string
}

func (c *UnifiedClient) NewListUnifiedLoanRecordsService() *ListUnifiedLoanRecordsService {
	return &ListUnifiedLoanRecordsService{c: c, params: map[string]string{}}
}

// SetType filters by record type: "borrow" (borrowing) or "repay" (repayment).
func (s *ListUnifiedLoanRecordsService) SetType(recordType string) *ListUnifiedLoanRecordsService {
	s.params["type"] = recordType
	return s
}

// SetCurrency narrows the result to a single currency (e.g. USDT).
func (s *ListUnifiedLoanRecordsService) SetCurrency(currency string) *ListUnifiedLoanRecordsService {
	s.params["currency"] = currency
	return s
}

// SetPage selects the page number.
func (s *ListUnifiedLoanRecordsService) SetPage(page int) *ListUnifiedLoanRecordsService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit caps the number of records returned (default 100, max 100).
func (s *ListUnifiedLoanRecordsService) SetLimit(limit int) *ListUnifiedLoanRecordsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *ListUnifiedLoanRecordsService) Do(ctx context.Context) ([]UnifiedLoanRecord, error) {
	req := request.Get(ctx, s.c, "/api/v4/unified/loan_records", s.params).WithSign()
	resp, err := request.Do[[]UnifiedLoanRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// UnifiedLoanRecord is a single borrow or repayment history entry.
type UnifiedLoanRecord struct {
	ID            int64           `json:"id"`
	Type          string          `json:"type"`
	RepaymentType string          `json:"repayment_type"`
	BorrowType    string          `json:"borrow_type"`
	CurrencyPair  string          `json:"currency_pair"`
	Currency      string          `json:"currency"`
	Amount        decimal.Decimal `json:"amount"`
	CreateTime    time.Time       `json:"create_time,format:unix"`
}

// ListUnifiedLoanInterestRecordsService -- GET /api/v4/unified/interest_records (private)
//
// Returns the interest-deduction history of the authenticated unified account.
type ListUnifiedLoanInterestRecordsService struct {
	c      *UnifiedClient
	params map[string]string
}

func (c *UnifiedClient) NewListUnifiedLoanInterestRecordsService() *ListUnifiedLoanInterestRecordsService {
	return &ListUnifiedLoanInterestRecordsService{c: c, params: map[string]string{}}
}

// SetCurrency narrows the result to a single currency (e.g. USDT).
func (s *ListUnifiedLoanInterestRecordsService) SetCurrency(currency string) *ListUnifiedLoanInterestRecordsService {
	s.params["currency"] = currency
	return s
}

// SetPage selects the page number.
func (s *ListUnifiedLoanInterestRecordsService) SetPage(page int) *ListUnifiedLoanInterestRecordsService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit caps the number of records returned (default 100, max 100).
func (s *ListUnifiedLoanInterestRecordsService) SetLimit(limit int) *ListUnifiedLoanInterestRecordsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetFrom sets the start time of the query window.
func (s *ListUnifiedLoanInterestRecordsService) SetFrom(from time.Time) *ListUnifiedLoanInterestRecordsService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo sets the end time of the query window (defaults to now server-side).
func (s *ListUnifiedLoanInterestRecordsService) SetTo(to time.Time) *ListUnifiedLoanInterestRecordsService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

// SetType filters by loan type: "platform" (platform borrowing) or "margin"
// (margin borrowing). Defaults to margin server-side.
func (s *ListUnifiedLoanInterestRecordsService) SetType(loanType string) *ListUnifiedLoanInterestRecordsService {
	s.params["type"] = loanType
	return s
}

func (s *ListUnifiedLoanInterestRecordsService) Do(ctx context.Context) ([]UnifiedLoanInterestRecord, error) {
	req := request.Get(ctx, s.c, "/api/v4/unified/interest_records", s.params).WithSign()
	resp, err := request.Do[[]UnifiedLoanInterestRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// UnifiedLoanInterestRecord is a single interest-deduction history entry.
type UnifiedLoanInterestRecord struct {
	Currency     string          `json:"currency"`
	CurrencyPair string          `json:"currency_pair"`
	ActualRate   decimal.Decimal `json:"actual_rate"`
	Interest     decimal.Decimal `json:"interest"`
	Status       int             `json:"status"`
	Type         string          `json:"type"`
	CreateTime   time.Time       `json:"create_time,format:unix"`
}
