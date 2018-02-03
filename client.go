package main

import (
	"bufio"
	"net"
)

func NewClient() *Client {
	return &Client{
		BroadcastChann: make(chan []byte),
		Done:           make(chan bool),
	}
}

type Client struct {
	conn           net.Conn
	BroadcastChann chan []byte
	Done           chan bool
	Name           string
	ID             string
	Topics         []<-chan string
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
