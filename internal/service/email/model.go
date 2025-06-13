package email

const (
	_subject   = "subject"
	_userEmail = "userEmail"
	_text      = "text"
)

type Email struct {
	Body map[string]string
}

const (
	_defaultMimeTextPlain = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	_defaultTmplHeader    = `<!doctype html><html><head><meta name="viewport" content="width=device-width" /></head><body>`
	_defaultTmplFooter    = "</body></html>"
)
