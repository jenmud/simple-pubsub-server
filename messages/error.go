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

type NoPublishers struct {
	Tick    string `xml:"tick"`
	Message string `xml:"message"`
}

func (np NoPublishers) Error() string {
	return fmt.Sprint("%s: %s\n", np.Tick, np.Message)
}

type AlreadyBeingPublished struct {
	Tick    string `xml:"tick"`
	Message string `xml:"message"`
}

func (abp AlreadyBeingPublished) Error() string {
	return fmt.Sprint("%s: %s\n", abp.Tick, abp.Message)
}
