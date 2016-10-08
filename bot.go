package main

import (
	"fmt"
	"log"
	"strings"

	"golang.org/x/net/websocket"
)

func main() {
	ws, botID, err := connect(apiKeys["slack"])
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Connected with bot ID %q\n", botID)

	err = verifyConnection(ws)
	if err != nil {
		log.Fatal(err)
	}

	listen(ws, botID)
}

func verifyConnection(ws *websocket.Conn) error {
	msg, err := getMessage(ws)
	if err != nil {
		return err
	}

	// First message to receive from the websocket should be "hello"
	if msg.Type != "hello" {
		return fmt.Errorf("Connection error, received: %+v", msg)
	}

	return nil
}

func listen(ws *websocket.Conn, botID botID) {
	for {
		msg, err := getMessage(ws)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Message:%+v\n\n", msg)

		if msg.Type == "message" {
			if strings.Contains(msg.Text, fmt.Sprintf("<@%s>", botID)) {
				log.Printf("Bot mentioned!")
			}
		}
	}
}
