package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"net/http"
)

type WSContext struct {
	Request  *models.WSMessageReq
	Response *models.WSMessageResp
	handler  func(ctx *WSContext)
}

func (c *WSContext) WriteResponse(status int, userID uint64, msgType string, data []byte) {
	c.Response = &models.WSMessageResp{
		UserID: userID,
		Status: status,
		Type:   msgType,
		Data:   data,
	}
}

type WSServer struct {
	hub *Hub
	*WSRouter
	toReceive chan *MSG
	toSend    chan *MSG
}

type WSRouter struct {
	routes map[string]func(c *WSContext)
}

func NewWSRouter() *WSRouter {
	return &WSRouter{
		routes: make(map[string]func(c *WSContext)),
	}
}

func (r *WSRouter) SetHandlerFunc(msgType string, handler func(c *WSContext)) {
	r.routes[msgType] = handler
}

func NewWSServer() *WSServer {
	h := NewHub()
	return &WSServer{
		hub:       h,
		WSRouter:  NewWSRouter(),
		toReceive: h.toReceive,
		toSend:    h.toSend,
	}
}

func (s *WSServer) Run() {
	s.hub.Run()
	go s.handleMessages()
}

func (s *WSServer) Stop() {
	s.hub.Stop()
}

func (s *WSServer) RegisterClient(w http.ResponseWriter, r *http.Request, userID uint64) error {
	Upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		//log err
		fmt.Println(err)
		return err
	}

	s.hub.RegisterClient(conn, userID)

	return nil
}

func (s *WSServer) handleMessages() {
	for m := range s.toReceive {
		msg := &models.WSMessageReq{}

		if err := json.Unmarshal(m.Data, &msg); err != nil {
			//log err
			fmt.Println(err)
			fmt.Println("pltcm")
			continue
		}
		msg.UserID = m.UserID

		c := &WSContext{Request: msg}

		handler, ok := s.routes[msg.Type]
		if !ok {
			//log no handler for type
			fmt.Println("no handler")
			continue
		}

		handler(c)

		s.SendMessage(c.Response)
	}
}

func (s *WSServer) SendMessage(msg *models.WSMessageResp) {
	defer func() {
		if recover() != nil {
			//log recover
		}
	}()

	data, err := json.Marshal(msg)
	if err != nil {
		//log err
		fmt.Println(err)
		return
	}

	s.toSend <- &MSG{
		UserID: msg.UserID,
		Data:   data,
	}
}
