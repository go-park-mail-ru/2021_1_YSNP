package websocket

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/chat"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/middleware"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/websocket"
	"github.com/gorilla/mux"
	"net/http"
)

type ChatWSHandler struct {
	chatUcase chat.ChatUsecase
}

func NewChatWSHandler(chatUcase chat.ChatUsecase) *ChatWSHandler {
	return &ChatWSHandler{
		chatUcase: chatUcase,
	}
}

func (ch *ChatWSHandler) Configure(r *mux.Router, mw *middleware.Middleware, server *websocket.WSServer) {
	server.SetHandlerFunc("CreateMessageReq", ch.CreateMessage)
	server.SetHandlerFunc("GetLastNMessagesReq", ch.GetLastNMessages)
	server.SetHandlerFunc("GetNMessagesBeforeReq", ch.GetNMessagesBefore)
}

func (ch *ChatWSHandler) CreateMessage(ctx *websocket.WSContext) {
	userID := ctx.Request.UserID
	req := &models.CreateMessageReq{}
	err := json.Unmarshal(ctx.Request.Data, req)
	if err != nil {
		//log err
		errE := errors.UnexpectedBadRequest(err)
		ctx.WriteResponse(errE.HttpError, userID, ctx.Request.Type, errors.JSONError(errE))
		return
	}

	resp, errE := ch.chatUcase.CreateMessage(req, userID)
	if errE != nil {
		//log err
		ctx.WriteResponse(errE.HttpError, userID, ctx.Request.Type, errors.JSONError(errE))
		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
			//log err
		errE := errors.UnexpectedInternal(err)
		ctx.WriteResponse(errE.HttpError, userID, ctx.Request.Type, errors.JSONError(errE))
		return
	}

	ctx.WriteResponse(http.StatusOK, userID, ctx.Request.Type, data)
}

func (ch *ChatWSHandler) GetLastNMessages(ctx *websocket.WSContext){
	userID := ctx.Request.UserID
	req := &models.GetLastNMessagesReq{}
	err := json.Unmarshal(ctx.Request.Data, req)
	if err != nil {
		//log err
		errE := errors.UnexpectedBadRequest(err)
		ctx.WriteResponse(errE.HttpError, userID, ctx.Request.Type, errors.JSONError(errE))
		return
	}

	resp, errE := ch.chatUcase.GetLastNMessages(req)
	if errE != nil {
		//log err
		ctx.WriteResponse(errE.HttpError, userID, ctx.Request.Type, errors.JSONError(errE))
		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		//log err
		errE := errors.UnexpectedInternal(err)
		ctx.WriteResponse(errE.HttpError, userID, ctx.Request.Type, errors.JSONError(errE))
		return
	}

	ctx.WriteResponse(http.StatusOK, userID, ctx.Request.Type, data)
}

func (ch *ChatWSHandler) GetNMessagesBefore(ctx *websocket.WSContext){
	userID := ctx.Request.UserID
	req := &models.GetNMessagesBeforeReq{}
	err := json.Unmarshal(ctx.Request.Data, req)
	if err != nil {
		//log err
		errE := errors.UnexpectedBadRequest(err)
		ctx.WriteResponse(errE.HttpError, userID, ctx.Request.Type, errors.JSONError(errE))
		return
	}

	resp, errE := ch.chatUcase.GetNMessagesBefore(req)
	if errE != nil {
		//log err
		ctx.WriteResponse(errE.HttpError, userID, ctx.Request.Type, errors.JSONError(errE))
		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		//log err
		errE := errors.UnexpectedInternal(err)
		ctx.WriteResponse(errE.HttpError, userID, ctx.Request.Type, errors.JSONError(errE))
		return
	}

	ctx.WriteResponse(http.StatusOK, userID, ctx.Request.Type, data)
}