package http

import (
	"context"
	"encoding/json"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/chat"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	errors "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
	log "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/logger"
	middleware2 "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/middleware"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/websocket"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net/http"
	"strconv"
	"time"
)

type ChatHandler struct {
	chatUcase chat.ChatUsecase
	logrusLogger *logrus.Entry
}

func NewChatHandler(chatUcase chat.ChatUsecase) *ChatHandler {
	return &ChatHandler{
		chatUcase: chatUcase,
	}
}

func (ch *ChatHandler) NewLogger(logger *logrus.Entry) {
	ch.logrusLogger = logger
}

func (ch *ChatHandler) AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	ch.logrusLogger = ch.logrusLogger.WithFields(logrus.Fields{
		"method":  info.FullMethod,
		"request": req,
		"work_id": uuid.New(),
	})
	ch.logrusLogger.Info("Get connection")

	start := time.Now()


	reply, err := handler(ctx, req)

	if err != nil {
		ch.logrusLogger.WithFields(logrus.Fields{
			"req": req,
			"error": err.Error(),
		}).Error("Error occurred")
	}

	ch.logrusLogger.WithFields(logrus.Fields{
		"work_time": time.Since(start),
	}).Info("Fulfilled connection")
	return reply, err
}

func (ch *ChatHandler) Configure(r *mux.Router, mw *middleware2.Middleware, server *websocket.WSServer) {
	r.HandleFunc("/chat/new", mw.CheckAuthMiddleware(ch.CreateChat)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/chat/list", mw.SetCSRFToken(mw.CheckAuthMiddleware(ch.GetUserChats))).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/chat/{cid:[0-9]+}", mw.SetCSRFToken(mw.CheckAuthMiddleware(ch.GetChatByID))).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/chat/ws", mw.CheckAuthMiddleware(ch.ServeWs(server))).Methods(http.MethodGet, http.MethodOptions)
}

func (ch *ChatHandler) CreateChat(w http.ResponseWriter, r *http.Request) {
	logger, ok := r.Context().Value(middleware2.ContextLogger).(*logrus.Entry)
	if !ok {
		logger = log.GetDefaultLogger()
		logger.Warn("no logger")
	}
	defer r.Body.Close()

	userID, ok := r.Context().Value(middleware2.ContextUserID).(uint64)
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

func (ch *ChatHandler) GetUserChats(w http.ResponseWriter, r *http.Request) {
	logger, ok := r.Context().Value(middleware2.ContextLogger).(*logrus.Entry)
	if !ok {
		logger = log.GetDefaultLogger()
		logger.Warn("no logger")
	}
	defer r.Body.Close()

	userID, ok := r.Context().Value(middleware2.ContextUserID).(uint64)
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

func (ch *ChatHandler) GetChatByID(w http.ResponseWriter, r *http.Request) {
	logger, ok := r.Context().Value(middleware2.ContextLogger).(*logrus.Entry)
	if !ok {
		logger = log.GetDefaultLogger()
		logger.Warn("no logger")
	}
	defer r.Body.Close()

	userID, ok := r.Context().Value(middleware2.ContextUserID).(uint64)
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
		logger, ok := r.Context().Value(middleware2.ContextLogger).(*logrus.Entry)
		if !ok {
			logger = log.GetDefaultLogger()
			logger.Warn("no logger")
		}
		defer r.Body.Close()

		userID, ok := r.Context().Value(middleware2.ContextUserID).(uint64)
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
