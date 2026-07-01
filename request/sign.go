package request

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"

	"github.com/UnipayFI/go-gate/common"
)

// SignFn mirrors client.SignFn: it turns the prehash string into the SIGN
// header value, given the configured secret.
type SignFn = func(secret, prehash string) (signature string, err error)

// HMACSign is Gate's default request signer:
//
//	SIGN = hex( HMAC-SHA512( secretKey, prehash ) )
//
// where prehash = method + "\n" + path + "\n" + query + "\n" +
// hex(SHA512(body)) + "\n" + timestamp (see Request.prepare).
func HMACSign(secret, prehash string) (string, error) {
	mac := hmac.New(sha512.New, common.StringToBytes(secret))
	if _, err := mac.Write(common.StringToBytes(prehash)); err != nil {
		return "", err
	}
	return hex.EncodeToString(mac.Sum(nil)), nil
}

// hashPayload returns hex(SHA512(body)) — the request-body component of the
// signature prehash. An empty body hashes the empty string, matching Gate.
func hashPayload(body string) string {
	sum := sha512.Sum512(common.StringToBytes(body))
	return hex.EncodeToString(sum[:])
}
