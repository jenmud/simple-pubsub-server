package messages

type Welcome struct {
	Name     string `xml:"name"`
	Address  string `xml:"address"`
	Datetime string `xml:"datetime"`
}
