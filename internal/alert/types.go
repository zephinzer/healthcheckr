package alert

type Channel interface {
	GetType() string
	Send(message Message) error
}

type Message struct {
	Title   string
	Body    string
	Details string
	Type    string
}
