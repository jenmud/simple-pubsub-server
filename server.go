package main

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"gitlab.com/jenmud/myBetex/messages"
	"log"
	"net"
	"time"
)

const heartbeatInterval = 10 * time.Second

func NewServer() *Server {
	server := Server{
		Clients:        make(map[string]*Client),
		BroadcastChann: make(chan []byte),
		Topics:         make(map[string]*Topic),
		publishers:     make(map[string][]*Topic),
	}

	return &server
}

type Server struct {
	Clients        map[string]*Client
	BroadcastChann chan []byte
	Topics         map[string]*Topic
	publishers     map[string][]*Topic
}

func (s *Server) AddClient(client *Client) error {
	name := client.Address()

	if _, ok := s.Clients[name]; ok {
		msg := fmt.Sprintf("Client %q already exists, please use a different name", name)
		return errors.New(msg)
	}

	s.Clients[name] = client
	log.Printf("Registered clients %d", len(s.Clients))
	return nil
}

func (s *Server) RemoveClient(client *Client) error {
	name := client.Address()
	delete(s.Clients, name)
	log.Printf("Registered clients %d", len(s.Clients))
	return nil
}

func (s *Server) HasConnected(client *Client) bool {
	_, ok := s.Clients[client.Address()]
	return ok
}

func (s *Server) HandleBroadcasts() {
	for {
		msg := <-s.BroadcastChann
		for c := range s.Clients {
			s.Clients[c].BroadcastChann <- msg
		}
	}
}

func (s *Server) Broadcast(msg []byte) {
	s.BroadcastChann <- msg
}

func (s *Server) HandlePingMsg(replyChan chan<- []byte) error {
	pong := messages.Pong{
		Tick: fmt.Sprint("%s", time.Now()),
	}

	output, err := xml.Marshal(pong)
	if err != nil {
		return err
	}

	replyChan <- output
	return nil
}

func (s *Server) HandleConnectMsg(client *Client, msg []byte) error {
	var err error
	var connMsg messages.Connect

	if err := xml.Unmarshal(msg, &connMsg); err != nil {
		return err
	}

	client.ID = connMsg.ID
	client.Name = connMsg.Name

	if err := s.AddClient(client); err != nil {
		return err
	}

	count := 0
	topics := make([]string, len(s.Topics))
	for topic, _ := range s.Topics {
		topics[count] = topic
		count++
	}

	welcomeMsg := messages.Welcome{
		Name:     connMsg.Name,
		Address:  client.Address(),
		Datetime: time.Now().String(),
		Topics:   topics,
	}

	welcomeMsgBytes, err := xml.Marshal(welcomeMsg)
	if err != nil {
		log.Printf("(Client: %s) Welcome MSG marshal error: %s", client.Address(), err)
		return err
	}

	if err := client.Send(welcomeMsgBytes); err != nil {
		log.Printf("(Client: %s) Send welcome MSG error: %s", client.Address(), err)
		return err
	}

	return err
}

func (s *Server) HandleSubscribeMsg(client *Client, msg []byte) error {
	var err error
	var subscribe messages.Subscribe

	if !s.HasConnected(client) {
		return client.SendError("NotConnected", "Not connected, please first connect.")
	}

	if err := xml.Unmarshal(msg, &subscribe); err != nil {
		return err
	}

	if subscribe.Topic != "" {
		topic, ok := s.Topics[subscribe.Topic]

		if !ok {
			msg := fmt.Sprintf("There are no publishers for topic %s", subscribe.Topic)
			return client.SendError("NoPublishers", msg)
		}

		topic.Subscribe(client)
	}

	return err
}

func (s *Server) HandlePublishMsg(client *Client, msg []byte) error {
	var publishMsg messages.Publish

	if err := xml.Unmarshal(msg, &publishMsg); err != nil {
		return err
	}

	// Check if we already have a client publishing to the topic
	topic, ok := s.Topics[publishMsg.Topic]
	if ok {
		msg := fmt.Sprintf("Topic %s is already being published by client %s", topic.Name, topic.Publisher.Address())
		return client.SendError("AlreadyBeingPublished", msg)
	}

	// Create a topic and make the client the publisher
	topic = NewTopic(publishMsg.Topic, client)
	s.Topics[topic.Name] = topic
	s.publishers[client.Address()] = append(s.publishers[client.Address()], topic)
	return nil
}

func (s *Server) HandlePublisherDisconnect(client *Client) error {
	log.Printf("Disconnecting publisher %s and closing any topics it has open", client.Address())
	topics, ok := s.publishers[client.Address()]
	if ok {
		for _, topic := range topics {
			topic.Close()
			log.Printf("Removing topic %s", topic.Name)
			delete(s.Topics, topic.Name)
		}
	}

	return nil
}

func (s *Server) HandleDisconnectMsg(client *Client) error {
	log.Printf("Client requested to disconnect, removing %s", client.Address())
	var err error

	disconnectTime := fmt.Sprintf("%s", time.Now())
	err = s.RemoveClient(client)
	if err != nil {
		return err
	}

	output, err := xml.Marshal(messages.Bye{Tick: disconnectTime})
	if err != nil {
		return err
	}

	err = client.Send(output)
	if err != nil {
		return err
	}

	client.Done <- true
	return err
}

func (s *Server) Dispatch(client *Client, msg []byte) error {
	var err error

	decoder := xml.NewDecoder(bytes.NewReader(msg))
	token, err := decoder.Token()
	if err != nil {
		return err
	}

	// This should have a low overhead as we are not parsing the entire
	// message, but only dealing with the first token. NOTE: this is only
	// a theory and not tested!
	switch t := token.(type) {
	case xml.StartElement:
		switch t.Name.Local {
		case "connect":
			if err := s.HandleConnectMsg(client, msg); err != nil {
				return err
			}
		case "ping":
			if err := s.HandlePingMsg(client.BroadcastChann); err != nil {
				return err
			}
		case "subscribe":
			if err := s.HandleSubscribeMsg(client, msg); err != nil {
				return err
			}
		case "publish":
			if err := s.HandlePublishMsg(client, msg); err != nil {
				return err
			}
		case "disconnect":
			if err := s.HandlePublisherDisconnect(client); err != nil {
				return err
			}

			if err := s.HandleDisconnectMsg(client); err != nil {
				return err
			}
		}
	}

	return err
}

func (s *Server) HandleClientConnection(client *Client) {
	defer s.HandlePublisherDisconnect(client)
	defer s.RemoveClient(client)
	log.Printf("New client connection from %s", client.Address())

	for {
		msg, err := client.Read('\n')
		if err != nil {
			log.Printf("(Client: %s) Read error: %s", client.Address(), err)
			return
		}

		if err := s.Dispatch(client, msg); err != nil {
			log.Printf("(Client: %s) Message parsing error: %s", client.Address(), err)
			return
		}
	}
}

func (s *Server) HandleHeartbeats(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		tick := <-ticker.C
		heartbeat := messages.Heartbeat{
			Tick: fmt.Sprintf("%s", tick),
		}

		output, err := xml.Marshal(heartbeat)
		if err != nil {
			log.Printf("Error marshaling heartbeat: %s", err)
			continue
		}

		s.Broadcast(output)
	}
	defer ticker.Stop()
}

func (s *Server) ListenAndAccept(port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	log.Printf("Server running and accepting connections on port %d", port)
	defer listener.Close()

	go s.HandleHeartbeats(heartbeatInterval)
	go s.HandleBroadcasts()
	defer close(s.BroadcastChann)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accepting error: %s", err)
			continue
		}

		log.Printf("Client connection accepted %s", conn.RemoteAddr().String())
		defer conn.Close()

		go s.HandleClientConnection(NewClient(conn))
	}
}
