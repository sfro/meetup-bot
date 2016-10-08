package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/websocket"
)

type startResponse struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
	Url   string `json:"url"`
	Self  self   `json:"self"`
}

type botId string

type self struct {
	Id             botId                  `json:"id"`
	Name           string                 `json:"name"`
	Prefs          map[string]interface{} `json:"prefs"`
	Created        int                    `json:"created"`
	ManualPresence string                 `json:"manual_presence"`
}

func connect(apiKey string) (*websocket.Conn, botId, error) {
	resp, err := http.Get(fmt.Sprintf("https://slack.com/api/rtm.start?token=%s", apiKey))
	if err != nil {
		return nil, "", err
	}

	if resp.StatusCode != 200 {
		return nil, "", fmt.Errorf("Expected status 200")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, "", err
	}

	startResponse := startResponse{}
	err = json.Unmarshal(body, &startResponse)
	if err != nil {
		return nil, "", err
	}

	if !startResponse.Ok {
		return nil, "", fmt.Errorf("Start response not ok, got: %+v", startResponse)
	}

	// Connect to the websocket using the URL we received
	webSocket, err := websocket.Dial(startResponse.Url, "", "https://api.slack.com/")
	if err != nil {
		return nil, "", err
	}

	// Return the bot ID, it is the ID that is returned in meesages when the bot is mentioned
	return webSocket, startResponse.Self.Id, nil
}

type Message struct {
	Id      int    `json:"id,omitempty"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

func getMessage(ws *websocket.Conn) (*Message, error) {
	msg := &Message{}
	return msg, websocket.JSON.Receive(ws, &msg)
}

var messageNumber int

func postMessage(ws *websocket.Conn, msg *Message) error {
	messageNumber++
	msg.Id = messageNumber
	return websocket.JSON.Send(ws, msg)
}
