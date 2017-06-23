package main

import (
	"log"
	"time"
	"net/http"
	"github.com/gorilla/websocket"
)

const (
	PING_INTERVAL = 30 * time.Second
	PONG_WAIT = 35 * time.Second
	WRITE_DELAY = 5 * time.Second
)

type ClientMessage struct {
	client *Client
	message Message
}

type Message struct {
	Username string `json:"username"`
	Message_text string `json:"message_text"`
}

type Client struct {
	hub *Hub //pointer to hub that this client is connected to
	ws_conn *websocket.Conn //websocket connection
	message_buffer chan BroadcastMessage
	joined bool //whether or not the user has joined the chat
}

func NewClient(hub *Hub, ws_conn *websocket.Conn) *Client {
	client := new(Client)
	client.hub = hub
	client.ws_conn = ws_conn
	client.message_buffer = make(chan BroadcastMessage)
	client.joined = false
	return client
}

var upgrader = websocket.Upgrader{}

/*
Reads messages from websocket and send to hub's broadcast channel
*/
func runClient(hub *Hub, w http.ResponseWriter, r *http.Request) {
	ws_conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer ws_conn.Close()

	client := NewClient(hub, ws_conn)
	hub.add <- client
	go client.hub_to_ws()
	client.ws_to_hub()
}

/*
Directs messages from websocket to hub
*/
func (client *Client) ws_to_hub() {
	defer func() {
		client.hub.remove <- client
	}()

	client.ws_conn.SetReadDeadline(time.Now().Add(PONG_WAIT))
	client.ws_conn.SetPongHandler(func(string) error {
		client.ws_conn.SetReadDeadline(time.Now().Add(PONG_WAIT))
		return nil
	})
	for {
		var message Message
		err := client.ws_conn.ReadJSON(&message)
		if err != nil {
			log.Printf("ReadJSON error: %v", err)
			return
		}
		client_message := ClientMessage{client: client, message: message}
		client.hub.broadcast <- client_message
	}
}

/*
Directs messages from hub to websocket
*/
func (client *Client) hub_to_ws() {
	timer := time.NewTicker(PING_INTERVAL)
	defer func() {
		timer.Stop()
		client.hub.remove <- client
	}()
	for {
		select { //write broadcast message to client's websocket
		case broadcast_message := <-client.message_buffer:
			client.ws_conn.SetWriteDeadline(time.Now().Add(WRITE_DELAY))
			err := client.ws_conn.WriteJSON(&broadcast_message)
			if err != nil {
				log.Printf("WriteJSON error: %v", err)
				return
			}
		case <-timer.C: //send a ping
			client.ws_conn.SetWriteDeadline(time.Now().Add(WRITE_DELAY))
			err := client.ws_conn.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				log.Printf("WriteMessage error: %v", err)
				return
			}
		}
	}
}