package telegram

type Message struct {
	ChatID int64
	Bot    string
	Text   string
}

const _defaultParseMode = "markdown"
