package event

import "strings"

func buildTopic(topic, lang string) string {
	if strings.HasSuffix(topic, lang) {
		return topic
	}
	return topic + _underscoreDelim + lang
}
