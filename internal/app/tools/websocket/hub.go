package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
)

type MSG struct {
	UserID uint64
	Data []byte
}

type Hub struct {
	clients map[uint64][]*Client
	register chan *Client
	unregister chan *Client

	toSend chan *MSG
	toReceive chan *MSG

	wg sync.WaitGroup
	mx sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[uint64][]*Client),
		register:   make(chan *Client, 256),
		unregister: make(chan *Client, 256),

		toSend:     make(chan *MSG, 256),
		toReceive:   make(chan *MSG, 256),

		wg:      sync.WaitGroup{},
		mx:		 sync.Mutex{},
	}
}

func (h *Hub) Run() {
	go h.registerClient()
	go h.unregisterClient()
	go h.sendMsgWorker()
	//go h.receiveMsgWorker()
}

func (h *Hub) Stop() {
	h.mx.Lock()

	for _, userClients := range h.clients {
		for _, cl := range userClients {
			close(cl.send)
		}
	}

	h.mx.Unlock()

	h.wg.Wait()

	close(h.register)
	close(h.unregister)
	close(h.toSend)
	close(h.toReceive)
}

func (h *Hub) RegisterClient(conn *websocket.Conn, userID uint64) {
	client := NewClient(h, conn, userID)
	client.Register()
}

func (h *Hub) registerClient() {
	for c := range h.register {

		h.mx.Lock()
		h.clients[c.userID] = append(h.clients[c.userID], c)
		h.mx.Unlock()

		go c.readPump()
		go c.writePump()

		h.wg.Add(2)
	}
}

func (h *Hub) unregisterClient() {
	for c := range h.unregister {
		close(c.send)

		h.mx.Lock()

		clients, ok := h.clients[c.userID]

		if ok {
			for i, client := range clients {
				if client == c {
					clients[i] = clients[len(clients)-1]
					clients = clients[:len(clients)-1]

					if len(clients) == 0 {
						delete(h.clients, c.userID)
					} else {
						h.clients[c.userID] = clients
					}

					break
				}
			}
		}

		h.mx.Unlock()
	}
}

func (h *Hub) sendMsgWorker() {
	for msg := range h.toSend {
		h.mx.Lock()

		clients, ok := h.clients[msg.UserID]
		if ok {
			for _, c := range clients {
				c.send <- msg.Data
			}
			//log msg
			fmt.Println(msg)
		} else {
			//log no clients
			fmt.Println("no clients")
		}

		h.mx.Unlock()
	}
}