package uid

import (
	"encoding/base64"
	"math/rand"
	"strings"
	"time"
)

func ShortUID() string {
	b := make([]byte, 6)
	_, _ = rand.New(rand.NewSource(time.Now().Unix())).Read(b)
	return removeChars(base64.URLEncoding.EncodeToString(b))
}

var _uidReplacer = strings.NewReplacer("_", "", "-", "", "=", "")

func removeChars(s string) string {
	var sb strings.Builder
	sb.Grow(len(s))
	sb.WriteString(_uidReplacer.Replace(s))
	return sb.String()
}
