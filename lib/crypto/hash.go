package crypto

import (
	"crypto/sha1"
	"encoding/base64"
)

func GetHash(data []byte) string {
	hasher := sha1.New()
	hasher.Write(data)
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
