package tinypng

import "encoding/base64"

func (t *tinyPng) getHeaders() map[string]string {
	key := "api" + ":" + t.apiKey
	authToken := base64.StdEncoding.EncodeToString([]byte(key))

	return map[string]string{
		"Authorization": "Basic " + authToken,
	}
}
