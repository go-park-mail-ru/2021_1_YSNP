package websocket

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/logger"
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

func NewWSServer(logger *logger.Logger) *WSServer {
	h := NewHub(logger)
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
	s.hub.log.GetLogger().Info("Start websocket server")
}

func (s *WSServer) Stop() {
	s.hub.Stop()
}

func (s *WSServer) RegisterClient(w http.ResponseWriter, r *http.Request, userID uint64) error {
	Upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.hub.log.LogWSError("server", userID, err.Error())
		return err
	}

	s.hub.RegisterClient(conn, userID)
	s.hub.log.LogWSInfo("server", userID, "Register new client with remote addr:"+conn.RemoteAddr().String())
	return nil
}

func (s *WSServer) handleMessages() {
	for m := range s.toReceive {
		msg := &models.WSMessageReq{}

		if err := json.Unmarshal(m.Data, &msg); err != nil {
			s.hub.log.LogWSError("server", m.UserID, err.Error())
			continue
		}
		msg.UserID = m.UserID

		c := &WSContext{Request: msg}

		handler, ok := s.routes[msg.Type]
		if !ok {
			s.hub.log.LogWSError("server", m.UserID, "No handler for type:"+msg.Type)
			continue
		}

		handler(c)

		s.SendMessage(c.Response, c.Response.UserID)
		if msg.Type == "CreateMessageReq" {
			typeData := &models.CreateMsgAdditData{}
			if err := json.Unmarshal(msg.Data.TypeData, &typeData); err != nil {
				s.hub.log.LogWSError("server", m.UserID, err.Error())
				return
			}
			s.SendMessage(c.Response, typeData.PartnerID)
		}
	}
}

func (s *WSServer) SendMessage(msg *models.WSMessageResp, reciever uint64) {
	defer func() {
		if recover() != nil {
			s.hub.log.LogWSInfo("server", msg.UserID, "Recover...")
		}
	}()

	data, err := json.Marshal(msg)
	if err != nil {
		s.hub.log.LogWSInfo("server", msg.UserID, err.Error())
		return
	}

	s.toSend <- &MSG{
		UserID: reciever,
		Data:   data,
	}
}
