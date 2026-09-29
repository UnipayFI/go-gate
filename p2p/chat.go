package p2p

import (
	"context"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
)

// GetChatsListService -- POST /api/v4/p2p/merchant/chat/get_chats_list (private)
//
// Returns the chat history for a P2P order, with incremental/backward paging.
type GetChatsListService struct {
	c    *P2PClient
	body map[string]any
}

func (c *P2PClient) NewGetChatsListService() *GetChatsListService {
	return &GetChatsListService{c: c, body: map[string]any{}}
}

// SetTxID selects the order whose chat to fetch (omit or 0 for the latest order
// with chat).
func (s *GetChatsListService) SetTxID(txid int64) *GetChatsListService {
	s.body["txid"] = txid
	return s
}

// SetLastReceived sets the timestamp of the last received message for forward
// incremental fetch.
func (s *GetChatsListService) SetLastReceived(lastReceived int64) *GetChatsListService {
	s.body["lastreceived"] = lastReceived
	return s
}

// SetFirstReceived sets the timestamp of the first received message for paging
// backward.
func (s *GetChatsListService) SetFirstReceived(firstReceived int64) *GetChatsListService {
	s.body["firstreceived"] = firstReceived
	return s
}

func (s *GetChatsListService) Do(ctx context.Context) (*P2PResponse[P2PChatList], error) {
	req := request.Post(ctx, s.c, "/api/v4/p2p/merchant/chat/get_chats_list", s.body).WithSign()
	return doP2P[P2PChatList](req)
}

// P2PChatList is the chat history for one P2P order plus its order-state markers.
type P2PChatList struct {
	Messages    []P2PChatMessage `json:"messages"`
	Memo        string           `json:"memo"`
	HasHistory  bool             `json:"has_history"`
	TxID        int64            `json:"txid"`
	SRVTM       time.Time        `json:"SRVTM,format:unix"`
	OrderStatus string           `json:"order_status"`
}

// P2PChatMessage is a single chat message in a P2P order conversation.
type P2PChatMessage struct {
	IsSell   int           `json:"is_sell"`
	MsgType  int           `json:"msg_type"`
	Msg      string        `json:"msg"`
	Username string        `json:"username"`
	Timest   time.Time     `json:"timest,format:unix"`
	MsgObj   P2PChatMsgObj `json:"msg_obj"`
	UID      string        `json:"uid"`
	Type     int           `json:"type"`
	Pic      string        `json:"pic"`
	FileKey  string        `json:"file_key"`
	FileType string        `json:"file_type"`
	RiskType int           `json:"risk_type"`
	ToastMsg string        `json:"toast_msg"`
}

// P2PChatMsgObj is the structured payload carried by system/template chat
// messages (order status changes, payment-method shares, cancellations).
type P2PChatMsgObj struct {
	Status         string    `json:"status"`
	Text           string    `json:"text"`
	PaymentVoucher []string  `json:"payment_voucher"`
	ReasonID       int       `json:"reason_id"`
	ToastID        int       `json:"toast_id"`
	ReasonMemo     string    `json:"reason_memo"`
	CancelTime     time.Time `json:"cancel_time,format:unix"`
	SellerConfirm  int       `json:"seller_confirm"`
	ID             string    `json:"id"`
	AccountDes     string    `json:"account_des"`
	PayType        string    `json:"pay_type"`
	File           string    `json:"file"`
	FileKey        string    `json:"file_key"`
	Account        string    `json:"account"`
	Memo           string    `json:"memo"`
	Code           string    `json:"code"`
	MemoExt        string    `json:"memo_ext"`
	TradeTips      string    `json:"trade_tips"`
	RealName       string    `json:"real_name"`
	IsDelete       int       `json:"is_delete"`
	PayName        string    `json:"pay_name"`
}

// SendChatMessageService -- POST /api/v4/p2p/merchant/chat/send_chat_message (private)
//
// Sends a text (or file-reference) chat message on a P2P order.
type SendChatMessageService struct {
	c    *P2PClient
	body map[string]any
}

func (c *P2PClient) NewSendChatMessageService(txid int64, message string) *SendChatMessageService {
	return &SendChatMessageService{c: c, body: map[string]any{
		"txid":    txid,
		"message": message,
	}}
}

// SetType sets the message type: 0 text, 1 file (image or video); defaults to 0.
func (s *SendChatMessageService) SetType(msgType int) *SendChatMessageService {
	s.body["type"] = msgType
	return s
}

func (s *SendChatMessageService) Do(ctx context.Context) (*P2PResponse[P2PSendChatResult], error) {
	req := request.Post(ctx, s.c, "/api/v4/p2p/merchant/chat/send_chat_message", s.body).WithSign()
	return doP2P[P2PSendChatResult](req)
}

// P2PSendChatResult is the result of sending a chat message.
type P2PSendChatResult struct {
	SRVTM          time.Time `json:"SRVTM,format:unix"`
	TxID           int64     `json:"txid"`
	ConversationID string    `json:"conversation_id"`
	MsgType        int       `json:"msg_type"`
	RiskType       int       `json:"risk_type"`
	ToastMsg       string    `json:"toast_msg"`
}

// UploadChatFileService -- POST /api/v4/p2p/merchant/chat/upload_chat_file (private)
//
// Uploads a base64-encoded image/video for use in a P2P order chat and returns
// its file key.
type UploadChatFileService struct {
	c    *P2PClient
	body map[string]any
}

func (c *P2PClient) NewUploadChatFileService(imageContentType, base64Img string) *UploadChatFileService {
	return &UploadChatFileService{c: c, body: map[string]any{
		"image_content_type": imageContentType,
		"base64_img":         base64Img,
	}}
}

func (s *UploadChatFileService) Do(ctx context.Context) (*P2PResponse[P2PUploadFileData], error) {
	req := request.Post(ctx, s.c, "/api/v4/p2p/merchant/chat/upload_chat_file", s.body).WithSign()
	return doP2P[P2PUploadFileData](req)
}

// P2PUploadFileData carries the file key of an uploaded chat file.
type P2PUploadFileData struct {
	FileKey string `json:"file_key"`
}
