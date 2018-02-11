package main

import (
	"fmt"
	"gitlab.com/jenmud/myBetex/messages"
	"log"
	"time"
)

func NewTopic(name string, publisher *Client) *Topic {
	return &Topic{
		Name:        name,
		Publisher:   publisher,
		Subscribers: make(map[string]*Client),
	}
}

type Topic struct {
	Name        string
	Publisher   *Client
	Subscribers map[string]*Client
}

func (t *Topic) Publish(msg []byte) {
	for _, client := range t.Subscribers {
		client.BroadcastChann <- msg
	}
}

func (t *Topic) SubscribePublisher(client *Client) error {
	if t.Publisher != nil {
		return messages.NoPublishers{
			Tick:    fmt.Sprint("%s", time.Now()),
			Message: fmt.Sprintf("Topic %s is already served by client %s", t.Name, t.Publisher.Address()),
		}
	}

	t.Publisher = client
	log.Printf("Client %s Publishing from topic %s", client.Address(), t.Name)
	return nil
}

func (t *Topic) Subscribe(client *Client) {
	if _, ok := t.Subscribers[client.Address()]; !ok {
		t.Subscribers[client.Address()] = client
		log.Printf("Subscribed client %s to topic %s", client.Address(), t.Name)
	}
}

func (t *Topic) Unsubscribe(client *Client) {
	delete(t.Subscribers, client.Address())
	log.Printf("Unsubscribed client %s from topic %s", client.Address(), t.Name)
}
