package messages

type Subscribe struct {
	Topic string `xml:"topic,attr,omitempty"`
}
