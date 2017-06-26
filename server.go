package main

type Server struct {
	hubs map[string] *Hub
}

func NewServer() *Server {
	server := new(Server)
	server.hubs = make(map[string] *Hub)
	return server
}