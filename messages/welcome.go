package messages

type Welcome struct {
	Name     string `xml:"name"`
	Datetime string `xml:"datetime"`
}
