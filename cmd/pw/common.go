package pw

import (
	"crypto/hmac"
	"crypto/sha256"
)

func hmacSha256(pass, data []byte) ([]byte, error) {
	h := hmac.New(sha256.New, pass)
	if _, err := h.Write([]byte(data)); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}
