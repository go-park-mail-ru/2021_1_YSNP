package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/chat"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
	log "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/logger"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/middleware"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/websocket"
)

type ChatHandler struct {
	chatUcase chat.ChatUsecase
}

func NewChatHandler(chatUcase chat.ChatUsecase) *ChatHandler {
	return &ChatHandler{
		chatUcase: chatUcase,
	}
}

func (ch *ChatHandler) Configure(r *mux.Router, mw *middleware.Middleware, server *websocket.WSServer) {
	r.HandleFunc("/chat/new", mw.CheckAuthMiddleware(ch.CreateChat)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/chat/list", mw.SetCSRFToken(mw.CheckAuthMiddleware(ch.GetUserChats))).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/chat/{cid:[0-9]+}", mw.SetCSRFToken(mw.CheckAuthMiddleware(ch.GetChatByID))).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/chat/ws", mw.CheckAuthMiddleware(ch.ServeWs(server))).Methods(http.MethodGet, http.MethodOptions)
}

// CreateChat godoc
// @Summary      Create chat
// @Description  Handler for creating chat
// @Tags         Chat
// @Accept       json
// @Produce      json
// @Param        body body models.ChatCreateReq true "Chat"
// @Success      200 {object} models.ChatResponse
// @Failure      400  {object}  errors.Error
// @Failure      404  {object}  errors.Error
// @Failure      500  {object}  errors.Error
// @Router      /chat/new [post]
func (ch *ChatHandler) CreateChat(w http.ResponseWriter, r *http.Request) {
	logger, ok := r.Context().Value(middleware.ContextLogger).(*logrus.Entry)
	if !ok {
		logger = log.GetDefaultLogger()
		logger.Warn("no logger")
	}
	defer r.Body.Close()

	userID, ok := r.Context().Value(middleware.ContextUserID).(uint64)
	if !ok {
		errE := errors.Cause(errors.UserUnauthorized)
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
	logger.Info("user id ", userID)

	req := &models.ChatCreateReq{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Error(err)
		errE := errors.UnexpectedBadRequest(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
	logger.Info("chat data ", req)

	resp, errE := ch.chatUcase.CreateChat(req, userID)
	if errE != nil {
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	body, err := json.Marshal(resp)
	if err != nil {
		logger.Error(err)
		errE := errors.UnexpectedInternal(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)

}

// GetUserChats godoc
// @Summary      Get user chats
// @Description  Handler for getting user chat
// @Tags         Chat
// @Accept       json
// @Produce      json
// @Success      200 {object} []models.ChatResponse
// @Failure      400  {object}  errors.Error
// @Failure      404  {object}  errors.Error
// @Failure      500  {object}  errors.Error
// @Router      /chat/list [get]
func (ch *ChatHandler) GetUserChats(w http.ResponseWriter, r *http.Request) {
	logger, ok := r.Context().Value(middleware.ContextLogger).(*logrus.Entry)
	if !ok {
		logger = log.GetDefaultLogger()
		logger.Warn("no logger")
	}
	defer r.Body.Close()

	userID, ok := r.Context().Value(middleware.ContextUserID).(uint64)
	if !ok {
		errE := errors.Cause(errors.UserUnauthorized)
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
	logger.Info("user id ", userID)

	resp, errE := ch.chatUcase.GetUserChats(userID)
	if errE != nil {
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	body, err := json.Marshal(resp)
	if err != nil {
		logger.Error(err)
		errE := errors.UnexpectedInternal(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

// GetChatByID godoc
// @Summary      Get chat by id
// @Description  Handler for getting chat
// @Tags         Chat
// @Accept       json
// @Produce      json
// @Param        cid path int64 true "Chat ID"
// @Success      200 {object} models.ChatResponse
// @Failure      400  {object}  errors.Error
// @Failure      404  {object}  errors.Error
// @Failure      500  {object}  errors.Error
// @Router      /chat/{cid} [get]
func (ch *ChatHandler) GetChatByID(w http.ResponseWriter, r *http.Request) {
	logger, ok := r.Context().Value(middleware.ContextLogger).(*logrus.Entry)
	if !ok {
		logger = log.GetDefaultLogger()
		logger.Warn("no logger")
	}
	defer r.Body.Close()

	userID, ok := r.Context().Value(middleware.ContextUserID).(uint64)
	if !ok {
		errE := errors.Cause(errors.UserUnauthorized)
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
	logger.Info("user id ", userID)

	vars := mux.Vars(r)
	chatID, _ := strconv.ParseUint(vars["cid"], 10, 64)
	logger.Info("product id ", chatID)

	resp, errE := ch.chatUcase.GetChatById(chatID, userID)
	if errE != nil {
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	body, err := json.Marshal(resp)
	if err != nil {
		logger.Error(err)
		errE := errors.UnexpectedInternal(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (ch *ChatHandler) ServeWs(srv *websocket.WSServer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger, ok := r.Context().Value(middleware.ContextLogger).(*logrus.Entry)
		if !ok {
			logger = log.GetDefaultLogger()
			logger.Warn("no logger")
		}
		defer r.Body.Close()

		userID, ok := r.Context().Value(middleware.ContextUserID).(uint64)
		if !ok {
			errE := errors.Cause(errors.UserUnauthorized)
			logger.Error(errE.Message)
			w.WriteHeader(errE.HttpError)
			w.Write(errors.JSONError(errE))
			return
		}
		//var userID uint64 = 2
		logger.Info("user id ", userID)

		if err := srv.RegisterClient(w, r, userID); err != nil {
			errE := errors.UnexpectedInternal(err)
			logger.Error(errE.Message)
			w.WriteHeader(errE.HttpError)
			w.Write(errors.JSONError(errE))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
