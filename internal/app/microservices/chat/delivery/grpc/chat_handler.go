package grpc

import (
	"context"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/microservices/chat"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	protoChat "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/proto/chat"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ChatServer struct {
	chatUcase chat.ChatUsecase
}

func NewChatServer(cu chat.ChatUsecase) *ChatServer {
	return &ChatServer{
		chatUcase: cu,
	}
}

func (cs *ChatServer) CreateChat(ctx context.Context, req *protoChat.ChatCreateReq) (*protoChat.ChatResp, error) {
	chatReq := &models.ChatCreateReq{
		ProductID: uint64(req.GetProductID()),
		PartnerID: uint64(req.GetPartnerID()),
	}

	chatResp, err := cs.chatUcase.CreateChat(chatReq, uint64(req.UserID))
	if err != nil {
		return nil, status.Error(codes.Code(err.ErrorCode), err.Message)
	}

	return models.ModelChatRespToGRPC(chatResp), nil
}

func (cs *ChatServer) GetChatByID(ctx context.Context, req *protoChat.GetChatByIDReq) (*protoChat.ChatResp, error) {
	chatResp, err := cs.chatUcase.GetChatById(uint64(req.GetChatID()), uint64(req.GetUserID()))
	if err != nil {
		return nil, status.Error(codes.Code(err.ErrorCode), err.Message)
	}

	return models.ModelChatRespToGRPC(chatResp), nil
}

func (cs *ChatServer) GetUserChats(ctx context.Context, userID *protoChat.UserID) (*protoChat.ChatRespArray, error) {
	chatResp, err := cs.chatUcase.GetUserChats(uint64(userID.GetUserID()))
	if err != nil {
		return nil, status.Error(codes.Code(err.ErrorCode), err.Message)
	}

	resps := &protoChat.ChatRespArray{}

	for _, resp := range chatResp {
		resps.Chats = append(resps.Chats, models.ModelChatRespToGRPC(resp))
	}

	return resps, nil
}

func (cs *ChatServer) CreateMessage(ctx context.Context, req *protoChat.CreateMessageReq) (*protoChat.MessageResp, error) {
	msgReq := &models.CreateMessageReq{
		ChatID:  uint64(req.GetChatID()),
		Content: req.GetContent(),
	}

	resp, err := cs.chatUcase.CreateMessage(msgReq, uint64(req.GetUserID()))
	if err != nil {
		return nil, status.Error(codes.Code(err.ErrorCode), err.Message)
	}

	return models.ModelMessageRespToGRPC(resp), nil
}

func (cs *ChatServer) GetLastNMessages(ctx context.Context, req *protoChat.GetLastNMessagesReq) (*protoChat.MessageRespArray, error) {
	msgReq := &models.GetLastNMessagesReq{
		UserID: uint64(req.UserID),
		ChatID: uint64(req.GetChatID()),
		Count:  int(req.GetCount()),
	}

	msgResp, err := cs.chatUcase.GetLastNMessages(msgReq)
	if err != nil {
		return nil, status.Error(codes.Code(err.ErrorCode), err.Message)
	}

	resps := &protoChat.MessageRespArray{}

	for _, resp := range msgResp {
		resps.Messages = append(resps.Messages, models.ModelMessageRespToGRPC(resp))
	}

	return resps, nil
}

func (cs *ChatServer) GetNMessagesBefore(ctx context.Context, req *protoChat.GetNMessagesReq) (*protoChat.MessageRespArray, error) {
	msgReq := &models.GetNMessagesBeforeReq{
		ChatID:        uint64(req.GetChatID()),
		Count:         int(req.GetCount()),
		LastMessageID: uint64(req.LastMessageId),
	}

	msgResp, err := cs.chatUcase.GetNMessagesBefore(msgReq)
	if err != nil {
		return nil, status.Error(codes.Code(err.ErrorCode), err.Message)
	}

	resps := &protoChat.MessageRespArray{}

	for _, resp := range msgResp {
		resps.Messages = append(resps.Messages, models.ModelMessageRespToGRPC(resp))
	}

	return resps, nil
}