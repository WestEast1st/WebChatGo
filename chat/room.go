package main

import (
	"log"
	"net/http"

	"../trace"
	"github.com/gorilla/websocket"
)

type room struct {
	//foward is 他clientからのメッセージを保持するチャンネル
	forward chan []byte
	//join は chat room に join 後の client chanelを入れる
	join chan *client
	//leave はチャットルームから退室しようとしているチャンネル
	leave chan *client
	//client には在室している全てのクライアントが保持
	clients map[*client]bool
	// tracerはチャットルーム場で行われた操作のログを受け取ります
	tracer trace.Tracer
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			//参加
			r.clients[client] = true
			r.tracer.Trace("  [*] new client join")
		case client := <-r.leave:
			//退室
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("  [*] new client off")
		case msg := <-r.forward:
			r.tracer.Trace("    [*] I received a message :", string(msg))
			for client := range r.clients {
				select {
				case client.send <- msg:
					//message send
					r.tracer.Trace("    [*] client submit message.")
				default:
					delete(r.clients, client)
					close(client.send)
					r.tracer.Trace("    [×] client submit message.")
				}
			}
		}
	}
}

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
	}
}
