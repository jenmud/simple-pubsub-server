package main

import (
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
		Topics:         make(map[string][]chan<- string),
	}

	return &server
}

type Server struct {
	Clients        map[string]*Client
	BroadcastChann chan []byte
	Topics         map[string][]chan<- string
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

func (s *Server) HandleBroadcasts() {
	for {
		select {
		case msg := <-s.BroadcastChann:
			for c := range s.Clients {
				s.Clients[c].BroadcastChann <- msg
			}
		}
	}
}

func (s *Server) Broadcast(msg []byte) {
	s.BroadcastChann <- msg
}

func (s *Server) SendPing(replyChan chan<- []byte) error {
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

func (s *Server) HandleMessage(client *Client, msg []byte) error {
	var connMsg messages.Connect

	switch {
	case xml.Unmarshal(msg, &connMsg) == nil:
		client.ID = connMsg.ID
		client.Name = connMsg.Name

		if err := s.AddClient(client); err != nil {
			return err
		}

		welcomeMsg := messages.Welcome{Name: connMsg.Name, Datetime: time.Now().String()}
		welcomeMsgBytes, err := xml.Marshal(welcomeMsg)
		if err != nil {
			log.Printf("(Client: %s) Welcome MSG marshal error: %s", client.Address(), err)
			return err
		}

		if err := client.Send(welcomeMsgBytes); err != nil {
			log.Printf("(Client: %s) Send welcome MSG error: %s", client.Address(), err)
			return err
		}
	default:
		return nil
	}

	return nil
}

func (s *Server) HandleConnection(client *Client) {
	defer s.RemoveClient(client)
	log.Printf("New client connection from %s", client.Address())

	for {
		input, err := client.Read('\n')
		if err != nil {
			log.Printf("(Client: %s) Read error: %s", client.Address(), err)
			return
		}

		go client.ListenForBroadcasts()
		if err := s.HandleMessage(client, input); err != nil {
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

		client := NewClient()
		client.conn = conn

		go s.HandleConnection(client)
	}
}
