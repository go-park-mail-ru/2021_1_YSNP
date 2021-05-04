package websocket

import (
	errs "errors"
	"github.com/gorilla/websocket"
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

	//authTime = 60 * time.Minute
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	userID uint64
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
}


func NewClient(h *Hub, conn *websocket.Conn, userID uint64) *Client {
	return &Client{
		userID: userID,
		hub:    h,
		conn:   conn,
		send:   make(chan []byte, 256),
	}
}

func (c *Client) Register() {
	c.hub.register <- c
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
		c.hub.wg.Done()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		c.hub.log.LogWSError(c.conn.RemoteAddr().String(), c.userID, err.Error())
	}
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.hub.log.LogWSError(c.conn.RemoteAddr().String(), c.userID, err.Error())
			} else {
				c.hub.log.LogWSError(c.conn.RemoteAddr().String(), c.userID, err.Error())
			}
			break
		}
		//message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.hub.toReceive <- &MSG{
			Data:   message,
			UserID: c.userID,
		}

		c.hub.log.LogWSInfo(c.conn.RemoteAddr().String(), c.userID, "New message")
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	//authTimer := time.NewTimer(authTime)

	defer func() {
		ticker.Stop()
		//authTimer.Stop()
		c.conn.Close()
		c.hub.wg.Done()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				c.hub.log.LogWSError(c.conn.RemoteAddr().String(), c.userID, err.Error())
			}
			if err := c.sendMessages(message, ok); err != nil {
				return
			}

		case <-ticker.C:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				c.hub.log.LogWSError(c.conn.RemoteAddr().String(), c.userID, err.Error())
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.hub.log.LogWSError(c.conn.RemoteAddr().String(), c.userID, err.Error())
				return
			}

		//case <-authTimer.C:
		//	if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
		//		//log err
		//	}
		//	//log err
		//	return
		}
	}
}

func (c *Client) sendMessages(message []byte, ok bool) error {
	if !ok {
		if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
			c.hub.log.LogWSError(c.conn.RemoteAddr().String(), c.userID, err.Error())
		}
		c.hub.log.LogWSInfo(c.conn.RemoteAddr().String(), c.userID, "Channel closed")
		return errs.New("")
	}

	if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
		//log err
		return err
	}

	n := len(c.send)
	for i := 0; i < n; i++ {
		if err := c.conn.WriteMessage(websocket.TextMessage, <-c.send); err != nil {
			//log err
			return err
		}
	}

	return nil
}
