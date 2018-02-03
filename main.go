package main

import (
	"log"
)

func main() {
	server := NewServer()
	if err := server.ListenAndAccept(3003); err != nil {
		log.Fatal(err)
	}
}
