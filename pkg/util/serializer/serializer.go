package serializer

import (
	"net/http"

	"github.com/bytedance/sonic"
)

func BodyToJSON(r *http.Request, value any) error {
	return sonic.ConfigDefault.NewDecoder(r.Body).Decode(&value)
}
