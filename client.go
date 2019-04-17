package main

import (
	"github.com/gorilla/websocket"
)

//client is one chat user
type client struct {
	//client is websocket
	socket *websocket.Conn
	//send is message channel
	send chan []byte
	//chan is client chat room
	room *room
}

//method read
func (c *client) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			c.room.forward <- msg
		} else {
			break
		}
	}
}

//method write
func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
