package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"gitlab.com/jenmud/myBetex/messages"
	"net"
	"time"
)

func NewClient(conn net.Conn) *Client {
	c := &Client{
		conn:           conn,
		BroadcastChann: make(chan []byte),
		Done:           make(chan bool),
	}

	go c.ListenForBroadcasts()
	return c
}

type Client struct {
	conn           net.Conn
	BroadcastChann chan []byte
	Done           chan bool
	Name           string
	ID             string
}

func (c *Client) Address() string {
	return c.conn.RemoteAddr().String()
}

func (c *Client) ListenForBroadcasts() {
	for {
		select {
		case msg := <-c.BroadcastChann:
			c.Send(msg)
		case <-c.Done:
			c.conn.Close()
		}
	}
}

func (c *Client) Read(delim byte) ([]byte, error) {
	rw := bufio.NewReader(c.conn)
	received, err := rw.ReadString(delim)
	return []byte(received), err
}

func (c *Client) Send(msg []byte) error {
	_, err := c.conn.Write(msg)
	return err
}

func (c *Client) SendError(errType, msg string) error {
	errMsg := messages.Error{
		Type:    errType,
		Tick:    fmt.Sprintf("%s", time.Now()),
		Message: msg,
	}

	output, err := xml.Marshal(errMsg)
	if err != nil {
		return err
	}

	return c.Send(output)
}
