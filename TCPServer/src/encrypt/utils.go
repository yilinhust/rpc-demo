package encrypt

import (
	"strings"
	"encoding/base64"
)


// Encode by Base64.URLEncoding.
func Encode(part []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(part), "=")
}

// Decode by Base64.URLEncoding.
func Decode(part string) ([]byte, error) {
	if l := len(part) % 4; l > 0 {
		part += strings.Repeat("=", 4-l)
	}

	return base64.URLEncoding.DecodeString(part)
}

