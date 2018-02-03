package messages

type Ping struct {
	From string `xml:"from,omitempty"`
}

type Pong struct {
	Tick string `xml:"tick`
	To   string `xml:"from,omitempty"`
}
