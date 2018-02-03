package messages

// Connect is a connection message.
type Connect struct {
	ID    string `xml:"id,attr"`
	Name  string `xml:"name,attr"`
	Topic string `xml:"source,omitempty"`
}
