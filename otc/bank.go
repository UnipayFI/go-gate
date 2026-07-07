package otc

import (
	"context"

	"github.com/UnipayFI/go-gate/v4/request"
)

// ListBanksService -- GET /api/v4/otc/bank/list (private)
//
// Returns the authenticated user's saved bank cards; each card's id is required
// when placing a fiat order.
type ListBanksService struct {
	c *OTCClient
}

func (c *OTCClient) NewListBanksService() *ListBanksService {
	return &ListBanksService{c: c}
}

func (s *ListBanksService) Do(ctx context.Context) (*OTCBankListResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/otc/bank/list").WithSign()
	return request.Do[OTCBankListResponse](req)
}

// OTCBankListResponse is the envelope returned by the bank-card list endpoint.
type OTCBankListResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    OTCBankList `json:"data"`
}

// OTCBankList holds the user's bank-card list and the server timestamp.
type OTCBankList struct {
	Lists     []OTCBankListItem `json:"lists"`
	Timestamp int64             `json:"timestamp"`
}

// OTCBankListItem is a single saved bank card. is_default is 1 when the card is
// the default and 0 otherwise; submit_time / update_time are formatted datetime
// strings.
type OTCBankListItem struct {
	ID                      string `json:"id"`
	BankAccountName         string `json:"bank_account_name"`
	BankName                string `json:"bank_name"`
	BankCountry             string `json:"bank_country"`
	BankAddress             string `json:"bank_address"`
	BankCode                string `json:"bank_code"`
	BranchCode              string `json:"branch_code"`
	IBAN                    string `json:"iban"`
	Swift                   string `json:"swift"`
	RemittanceLineNumber    string `json:"remittance_line_number"`
	AgentBankName           string `json:"agent_bank_name"`
	AgentBankSwift          string `json:"agent_bank_swift"`
	SubmitTime              string `json:"submit_time"`
	UpdateTime              string `json:"update_time"`
	Status                  string `json:"status"`
	DocumentationFileType   string `json:"documentation_file_type"`
	Memo                    string `json:"memo"`
	IsDefault               int    `json:"is_default"`
	BankID                  string `json:"bank_id"`
	DocumentationFileKeyURL string `json:"documentation_file_key_url"`
}

// CreateBankService -- POST /api/v4/otc/bank/create (private)
//
// Creates a bank card. documentation_file carries the account-opening proof file
// content (jpg/jpeg/png/pdf, sent as a base64/binary string).
type CreateBankService struct {
	c    *OTCClient
	body map[string]any
}

func (c *OTCClient) NewCreateBankService(bankAccountName, bankName, bankCountry, bankAddress, iban, swift, documentationFile string) *CreateBankService {
	return &CreateBankService{c: c, body: map[string]any{
		"bank_account_name":  bankAccountName,
		"bank_name":          bankName,
		"bank_country":       bankCountry,
		"bank_address":       bankAddress,
		"iban":               iban,
		"swift":              swift,
		"documentation_file": documentationFile,
	}}
}

// SetRemittanceLineNumber sets the optional remittance routing number.
func (s *CreateBankService) SetRemittanceLineNumber(remittanceLineNumber string) *CreateBankService {
	s.body["remittance_line_number"] = remittanceLineNumber
	return s
}

// SetAgentBankName sets the optional correspondent bank name.
func (s *CreateBankService) SetAgentBankName(agentBankName string) *CreateBankService {
	s.body["agent_bank_name"] = agentBankName
	return s
}

// SetAgentBankSwift sets the optional correspondent bank SWIFT code.
func (s *CreateBankService) SetAgentBankSwift(agentBankSwift string) *CreateBankService {
	s.body["agent_bank_swift"] = agentBankSwift
	return s
}

func (s *CreateBankService) Do(ctx context.Context) (*OTCBankCreateResponse, error) {
	req := request.Post(ctx, s.c, "/api/v4/otc/bank/create", s.body).WithSign()
	return request.Do[OTCBankCreateResponse](req)
}

// OTCBankCreateResponse is the envelope returned when creating a bank card.
type OTCBankCreateResponse struct {
	Code      int                 `json:"code"`
	Message   string              `json:"message"`
	Data      OTCBankCreateResult `json:"data"`
	Timestamp int64               `json:"timestamp"`
}

// OTCBankCreateResult carries the new bank card's primary key and review status.
type OTCBankCreateResult struct {
	BankID int64 `json:"bank_id"`
	Status int   `json:"status"`
}

// DeleteBankService -- POST /api/v4/otc/bank/delete (private)
//
// Deletes one of the user's saved bank cards by id.
type DeleteBankService struct {
	c    *OTCClient
	body map[string]any
}

func (c *OTCClient) NewDeleteBankService(bankID string) *DeleteBankService {
	return &DeleteBankService{c: c, body: map[string]any{
		"bank_id": bankID,
	}}
}

func (s *DeleteBankService) Do(ctx context.Context) (*OTCAckResponse, error) {
	req := request.Post(ctx, s.c, "/api/v4/otc/bank/delete", s.body).WithSign()
	return request.Do[OTCAckResponse](req)
}

// SetDefaultBankService -- POST /api/v4/otc/bank/set_default (private)
//
// Marks one of the user's saved bank cards as the default.
type SetDefaultBankService struct {
	c    *OTCClient
	body map[string]any
}

func (c *OTCClient) NewSetDefaultBankService(bankID string) *SetDefaultBankService {
	return &SetDefaultBankService{c: c, body: map[string]any{
		"bank_id": bankID,
	}}
}

func (s *SetDefaultBankService) Do(ctx context.Context) (*OTCAckResponse, error) {
	req := request.Post(ctx, s.c, "/api/v4/otc/bank/set_default", s.body).WithSign()
	return request.Do[OTCAckResponse](req)
}

// GetBankSupplementChecklistService -- GET /api/v4/otc/bank/bank_supplement_checklist (private)
//
// Returns the checklist of supplementary materials still required for a bank card.
type GetBankSupplementChecklistService struct {
	c      *OTCClient
	params map[string]string
}

func (c *OTCClient) NewGetBankSupplementChecklistService(bankID string) *GetBankSupplementChecklistService {
	return &GetBankSupplementChecklistService{c: c, params: map[string]string{
		"bank_id": bankID,
	}}
}

func (s *GetBankSupplementChecklistService) Do(ctx context.Context) (*OTCBankSupplementChecklistResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/otc/bank/bank_supplement_checklist", s.params).WithSign()
	return request.Do[OTCBankSupplementChecklistResponse](req)
}

// OTCBankSupplementChecklistResponse is the envelope returned by the bank-card
// supplement checklist endpoint.
type OTCBankSupplementChecklistResponse struct {
	Code    int                        `json:"code"`
	Message string                     `json:"message"`
	Data    OTCBankSupplementChecklist `json:"data"`
}

// OTCBankSupplementChecklist is the required-materials checklist for a bank card.
// user_type is "personal" or "enterprise".
type OTCBankSupplementChecklist struct {
	UserType  string                           `json:"user_type"`
	Items     []OTCBankSupplementChecklistItem `json:"items"`
	Timestamp int64                            `json:"timestamp"`
}

// OTCBankSupplementChecklistItem is one required material in the checklist. code
// is the material item code (the top-level key of relationship_proof).
type OTCBankSupplementChecklistItem struct {
	Code     string `json:"code"`
	Zh       string `json:"zh"`
	En       string `json:"en"`
	Required bool   `json:"required"`
}

// SubmitPersonalBankSupplementService -- POST /api/v4/otc/bank/personal/bank_supplement (private)
//
// Submits supplementary materials for a personal bank card (ID document front and
// back, plus a proof of address), each sent as a base64/binary string.
type SubmitPersonalBankSupplementService struct {
	c    *OTCClient
	body map[string]any
}

func (c *OTCClient) NewSubmitPersonalBankSupplementService(bankID, idDocumentFront, idDocumentBack, addressProof string) *SubmitPersonalBankSupplementService {
	return &SubmitPersonalBankSupplementService{c: c, body: map[string]any{
		"bank_id":           bankID,
		"id_document_front": idDocumentFront,
		"id_document_back":  idDocumentBack,
		"address_proof":     addressProof,
	}}
}

func (s *SubmitPersonalBankSupplementService) Do(ctx context.Context) (*OTCAckResponse, error) {
	req := request.Post(ctx, s.c, "/api/v4/otc/bank/personal/bank_supplement", s.body).WithSign()
	return request.Do[OTCAckResponse](req)
}

// SubmitEnterpriseBankSupplementService -- POST /api/v4/otc/bank/enterprise/bank_supplement (private)
//
// Submits supplementary materials for an enterprise bank card (business license,
// register of shareholders, legal-representative passport and ownership-structure
// chart, with optional proof of funds and additional materials), each sent as a
// base64/binary string.
type SubmitEnterpriseBankSupplementService struct {
	c    *OTCClient
	body map[string]any
}

func (c *OTCClient) NewSubmitEnterpriseBankSupplementService(bankID, certificate, shareHolders, passport, shareHoldingStructure string) *SubmitEnterpriseBankSupplementService {
	return &SubmitEnterpriseBankSupplementService{c: c, body: map[string]any{
		"bank_id":                 bankID,
		"certificate":             certificate,
		"share_holders":           shareHolders,
		"passport":                passport,
		"share_holding_structure": shareHoldingStructure,
	}}
}

// SetUID sets the optional user id the supplement is submitted for.
func (s *SubmitEnterpriseBankSupplementService) SetUID(uid string) *SubmitEnterpriseBankSupplementService {
	s.body["uid"] = uid
	return s
}

// SetFundsStatement sets the optional proof-of-funds file content.
func (s *SubmitEnterpriseBankSupplementService) SetFundsStatement(fundsStatement string) *SubmitEnterpriseBankSupplementService {
	s.body["funds_statement"] = fundsStatement
	return s
}

// SetAdditional sets the optional additional supplementary material file content.
func (s *SubmitEnterpriseBankSupplementService) SetAdditional(additional string) *SubmitEnterpriseBankSupplementService {
	s.body["additional"] = additional
	return s
}

func (s *SubmitEnterpriseBankSupplementService) Do(ctx context.Context) (*OTCAckResponse, error) {
	req := request.Post(ctx, s.c, "/api/v4/otc/bank/enterprise/bank_supplement", s.body).WithSign()
	return request.Do[OTCAckResponse](req)
}
