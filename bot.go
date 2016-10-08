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
		receivedMsg, err := getMessage(ws)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Message:%+v", receivedMsg)

		if receivedMsg.Type == "message" {
			if strings.Contains(receivedMsg.Text, fmt.Sprintf("<@%s>", botID)) {
				log.Printf("Bot mentioned!")

				if strings.Contains(strings.ToLower(receivedMsg.Text), "good morning") {
					returnMsg := &Message{
						Type:    "message",
						Channel: receivedMsg.Channel,
						Text:    "Good morning! (/◕ヮ◕)/",
					}

					err = postMessage(ws, returnMsg)
					if err != nil {
						log.Printf("ERROR: Error posting message: %+v", returnMsg)
					}
				}
			}
		}
	}
}
