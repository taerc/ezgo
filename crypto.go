package ezgo

import (
	"crypto/sha256"
	"encoding/hex"
)

/// SHA256 @description SHA256

func SHA256(src string) string {
	h := sha256.New()
	h.Write([]byte(src))
	sum := h.Sum(nil)
	s := hex.EncodeToString(sum)
	return s
}
