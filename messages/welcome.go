package messages

type Welcome struct {
	Name     string   `xml:"name,attr"`
	Address  string   `xml:"address,attr"`
	Datetime string   `xml:"datetime,attr"`
	Topics   []string `xml:"topics>topic"`
}
