package otc

import (
	"context"

	"github.com/UnipayFI/go-gate/v4/request"
)

// CreatePreUploadService -- POST /api/v4/otc/upload/pre_upload (private)
//
// Issues an S3 POST Policy for uploading a file to the temporary bucket. Upload
// the file directly to the returned URL with the returned fields as form-data
// (the file field last, success is HTTP 204), then pass the returned base64
// file_key unchanged to the business endpoint that consumes it — do not decode
// it. Unsubmitted files stay in the temporary bucket until the lifecycle rules
// reclaim them.
//
// contentType is the base64 of the file's MIME type, not the MIME type itself;
// a plaintext value containing "/" may be blocked by the gateway. The supported
// values are "aW1hZ2UvcG5n" (image/png), "aW1hZ2UvanBlZw==" (image/jpeg),
// "aW1hZ2UvanBn" (image/jpg) and "YXBwbGljYXRpb24vcGRm" (application/pdf).
type CreatePreUploadService struct {
	c    *OTCClient
	body map[string]any
}

func (c *OTCClient) NewCreatePreUploadService(contentType string) *CreatePreUploadService {
	return &CreatePreUploadService{c: c, body: map[string]any{
		"content_type": contentType,
	}}
}

// SetScene selects the business scene, which decides the temporary path and the
// production directory the file is moved to: "general" (the default, fiat buy
// payment receipts), "bank" (adding a card and bank-card supplements),
// "assessment" (professional verification) or "credit" (credit limit increase).
func (s *CreatePreUploadService) SetScene(scene string) *CreatePreUploadService {
	s.body["scene"] = scene
	return s
}

func (s *CreatePreUploadService) Do(ctx context.Context) (*OTCPreUploadResponse, error) {
	req := request.Post(ctx, s.c, "/api/v4/otc/upload/pre_upload", s.body).WithSign()
	return request.Do[OTCPreUploadResponse](req)
}

// OTCPreUploadResponse is the envelope returned by the pre-upload endpoint. Code
// is 0 on success and 10010400 for a parameter error.
type OTCPreUploadResponse struct {
	Code      int              `json:"code"`
	Message   string           `json:"message"`
	Data      OTCPreUploadData `json:"data"`
	Timestamp int64            `json:"timestamp"`
}

// OTCPreUploadData carries the temporary object key and the S3 direct-upload
// parameters. FileKey is the base64 object path to hand back to the business
// endpoint unchanged; ExpiresIn is the credential's lifetime in seconds
// (currently 5400, matching the expiration inside Fields.Policy). The policy
// caps the upload at 10 MB for every scene.
type OTCPreUploadData struct {
	FileKey   string                   `json:"file_key"`
	URL       string                   `json:"url"`
	Fields    OTCPreUploadPolicyFields `json:"fields"`
	ExpiresIn int                      `json:"expires_in"`
}

// OTCPreUploadPolicyFields are the S3 POST Policy form fields. Send every one of
// them unchanged during the direct upload. Key is the plaintext object path,
// identical to the base64 decoding of OTCPreUploadData.FileKey.
type OTCPreUploadPolicyFields struct {
	Key            string `json:"key"`
	ContentType    string `json:"Content-Type"`
	XAmzCredential string `json:"X-Amz-Credential"`
	XAmzAlgorithm  string `json:"X-Amz-Algorithm"`
	XAmzDate       string `json:"X-Amz-Date"`
	Policy         string `json:"Policy"`
	XAmzSignature  string `json:"X-Amz-Signature"`
}
