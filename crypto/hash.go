package crypto

import (
	"crypto/sha1"
	"encoding/hex"
)

func Sha1String(data string) string {
	sha1 := sha1.New()
	sha1.Write([]byte(data))
	return hex.EncodeToString(sha1.Sum([]byte(nil)))
}

