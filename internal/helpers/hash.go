package helpers

import (
	"crypto/sha1"
	"encoding/hex"
)

func Hash(original string) string {
	h := sha1.New()
	h.Write([]byte(original))
	return hex.EncodeToString(h.Sum(nil))
}
