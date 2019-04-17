package main

type room struct {
	//foward is 他clientからのメッセージを保持するチャンネル
	foward chan []byte
	//join は chat room に join 後の client chanelを入れる
	join chan *client
	//leave はチャットルームから退室しようとしているチャンネル
	leave chan *client
	//client には在室している全てのクライアントが保持
	clients map[*client]bool
}

func (r *room) rum() {
	for {
		select {
		case client := <-r.join:
			//参加
			r.clients[client] = true
		case client := <-r.leave:
			//退室
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.foward:
			for client := range r.clients {
				select {
				case client.send <- msg:
					//message send
				default:
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}
