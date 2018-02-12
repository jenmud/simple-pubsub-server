package messages

import (
	"fmt"
)

type Error struct {
	Tick    string `xml:"tick"`
	Type    string `xml:"type"`
	Message string `xml:"message"`
}

func (err Error) Error() string {
	return fmt.Sprint("%s %s: %s\n", err.Tick, err.Type, err.Message)
}
