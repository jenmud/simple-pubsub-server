package main

import (
	"bufio"
	"log"
	"net"
)

func NewClient(conn net.Conn) *Client {
	return &Client{
		conn:           conn,
		BroadcastChann: make(chan []byte),
		Done:           make(chan bool),
		Topics:         make(map[string]chan []byte),
	}
}

type Client struct {
	conn           net.Conn
	BroadcastChann chan []byte
	Done           chan bool
	Name           string
	ID             string
	Topics         map[string]chan []byte
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

func (c *Client) Subscribe(topicName string) chan<- []byte {
	topicHandler := func(topicChan chan []byte) {
		log.Printf("Client %q subscribed and listening to %q for messages", c.Address(), topicName)
		for {
			msg := <-topicChan
			c.Send(msg)
		}
	}

	if _, ok := c.Topics[topicName]; !ok {
		topicChan := make(chan []byte)
		c.Topics[topicName] = topicChan
		go topicHandler(topicChan)
		return topicChan
	}

	return c.Topics[topicName]
}

func (c *Client) Unsubscribe(topicName string) {
	delete(c.Topics, topicName)
}
