package push

type Message struct {
	UserID int
	Token  string
	Data   map[string]string
}
