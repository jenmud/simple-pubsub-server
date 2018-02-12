package messages

type Publish struct {
	Topic string `xml:"topic,attr"`
}
