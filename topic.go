package main

import (
	"log"
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

func (t *Topic) Close() {
	for _, client := range t.Subscribers {
		t.Unsubscribe(client)
	}
}
