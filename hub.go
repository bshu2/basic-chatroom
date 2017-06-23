package main

//import "log"

const (
	CHAT_MESSAGE = 0
	SYSTEM_MESSAGE = 1
)

type BroadcastMessage struct {
	Message_type int `json:"message_type"`
	Username string `json:"username"`
	Message_text string `json:"message_text"`
}

type Hub struct {
	clients map[*Client] bool //map of all connected clients
	broadcast chan ClientMessage //channel to broadcast messages
	add chan *Client //channel to add clients
	remove chan *Client //channel to remove clients
}

func NewHub() *Hub {
	hub := new(Hub)
	hub.clients = make(map[*Client] bool)
	hub.broadcast = make(chan ClientMessage)
	hub.add = make(chan *Client)
	hub.remove = make(chan *Client)
	return hub
}

func (hub *Hub) run() {
	for {
		select {
		//add client in add channel from clients map
		case client := <-hub.add:
			hub.clients[client] = true
			//log.Println("client added")

		//remove client in remove channel from clients map
		case client := <-hub.remove:
			_, ok := hub.clients[client]
			if ok {
				delete(hub.clients, client)
				client.ws_conn.Close()
				//log.Println("client removed")
			}

		//broadcast message in broadcast channel to all connected clients
		case client_message := <-hub.broadcast:
			var broadcast_message BroadcastMessage
			if !(client_message.client.joined) {
				broadcast_message.Message_type = SYSTEM_MESSAGE
				broadcast_message.Username = ""
				broadcast_message.Message_text = client_message.message.Username + " has joined the room."
				client_message.client.joined = true
			} else {
				broadcast_message.Message_type = CHAT_MESSAGE
				broadcast_message.Username = client_message.message.Username 
				broadcast_message.Message_text = client_message.message.Message_text
			}
			for client := range hub.clients {
				client.message_buffer <- broadcast_message
			}
		}
	}
}