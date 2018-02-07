package messages

import (
	"fmt"
)

type UnsubscribedError struct {
	Tick    string `xml:"tick"`
	Message string `xml:"message"`
}

func (ue UnsubscribedError) Error() string {
	return fmt.Sprint("%s: %s\n", ue.Tick, ue.Message)
}

type NotConnected struct {
	Tick    string `xml:"tick"`
	Message string `xml:"message"`
}

func (ne NotConnected) Error() string {
	return fmt.Sprint("%s: %s\n", ne.Tick, ne.Message)
}
