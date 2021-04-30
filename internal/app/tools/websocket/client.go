package websocket

import (
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024

	authTime = 60 * time.Minute
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { // Note: for tests
		return true
	},
}

type Client struct {
	userId int
	//hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
}


func NewClient(/*h *Hub,*/ conn *websocket.Conn, userId int) *Client {
	return &Client{
		userId: userId,
		//hub:    h,
		conn:   conn,
		send:   make(chan []byte, 256),
	}
}