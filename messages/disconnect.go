package messages

type Disconnect struct {
}

type Bye struct {
	Tick string `xml:"tick,omitempty"`
}
