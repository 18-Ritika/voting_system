package models

import "github.com/gorilla/websocket"

type Session struct {
	SessionID string
	Votes     map[string]string
	Clients   []*Client
}

type Client struct {
	ID   string
	Conn *websocket.Conn
}
