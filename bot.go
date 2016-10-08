package main

import (
	"fmt"
	"log"
	"strings"
	"time"

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

				returnMsg := &Message{
					Type:    "message",
					Channel: receivedMsg.Channel,
				}

				if strings.Contains(strings.ToLower(receivedMsg.Text), "good morning") {
					returnMsg.Text = "Good morning! (/◕ヮ◕)/"
				} else if strings.Contains(strings.ToLower(receivedMsg.Text), "meetup") &&
					strings.Contains(strings.ToLower(receivedMsg.Text), "suggest") {
					result, err := meetupSuggestion()
					if err != nil {
						log.Printf("ERROR: Error getting meetup suggestion: %+v", err)
						returnMsg.Text = "Failed getting a suggestion! (ToT)"
					} else {
						returnMsg.Text = fmt.Sprintf("Group *%s* is having meetup *%s*!\n", result.Group.Who, result.Name)
						returnMsg.Text += fmt.Sprintf("It's at *%s* and has *%v* spaces left!\n", time.Unix(result.Time/1000, 0).Format(time.UnixDate), result.RSVPLimit-result.YesRSVPCount)
						returnMsg.Text += fmt.Sprintf("You can RSVP here (^_-)-☆  %v", result.EventURL)
					}
				} else {
					returnMsg.Text = "What do you want from me? (ToT)"
				}

				err = postMessage(ws, returnMsg)
				if err != nil {
					log.Printf("ERROR: Error posting message: %+v", err)

				}
			}
		}
	}
}
