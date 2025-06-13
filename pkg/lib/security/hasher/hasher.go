package hasher

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

func GenerateSHA2(values ...string) (string, error) {
	if len(values) < 2 {
		return "", errors.New("hasher: not enough arguments")
	}

	hash := hmac.New(sha256.New, []byte(values[0]))

	for _, v := range values[1:] {
		_, _ = hash.Write([]byte(v))
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
