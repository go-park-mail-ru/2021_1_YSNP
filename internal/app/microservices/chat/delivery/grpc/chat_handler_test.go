package grpc

import (
	"context"
	"database/sql"
	mock "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/microservices/chat/mocks"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
	protoChat "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/proto/chat"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChatHandler_CreateChat_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatUcase := mock.NewMockChatUsecase(ctrl)

	req := &models.ChatCreateReq{ProductID: uint64(2), PartnerID: uint64(1)}

	chatHandler := NewChatServer(chatUcase)
	chatUcase.EXPECT().CreateChat(req, uint64(1)).Return(&models.ChatResponse{}, nil)

	chatHandler.CreateChat(context.Background(), &protoChat.ChatCreateReq{
		PartnerID: int64(1),
		ProductID: int64(2),
		UserID: int64(1),
	})

	//error
	chatUcase.EXPECT().CreateChat(req, uint64(1)).Return(nil, errors.UnexpectedInternal(sql.ErrConnDone))

	_, err := chatHandler.CreateChat(context.Background(), &protoChat.ChatCreateReq{
		PartnerID: int64(1),
		ProductID: int64(2),
		UserID: int64(1),
	})
	assert.Error(t, err)
}

func TestChatServer_GetChatByID_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatUcase := mock.NewMockChatUsecase(ctrl)

	chatHandler := NewChatServer(chatUcase)
	chatUcase.EXPECT().GetChatById(uint64(2), uint64(1)).Return(&models.ChatResponse{}, nil)

	chatHandler.GetChatByID(context.Background(), &protoChat.GetChatByIDReq{
		UserID: int64(1),
		ChatID: int64(2),
	})

	//error
	chatUcase.EXPECT().GetChatById(uint64(2), uint64(1)).Return(nil, errors.UnexpectedInternal(sql.ErrConnDone))

	_, err := chatHandler.GetChatByID(context.Background(), &protoChat.GetChatByIDReq{
		UserID: int64(1),
		ChatID: int64(2),
	})
	assert.Error(t, err)
}

func TestChatServer_GetUserChats_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatUcase := mock.NewMockChatUsecase(ctrl)

	chatHandler := NewChatServer(chatUcase)

	chat := &models.ChatResponse{PartnerID: 1}

	chatUcase.EXPECT().GetUserChats(uint64(1)).Return([]*models.ChatResponse{chat},nil)

	chatHandler.GetUserChats(context.Background(), &protoChat.UserID{UserID: int64(1)})

	//error
	chatUcase.EXPECT().GetUserChats(uint64(1)).Return(nil, errors.UnexpectedInternal(sql.ErrConnDone))

	_, err := chatHandler.GetUserChats(context.Background(), &protoChat.UserID{UserID: int64(1)})
	assert.Error(t, err)
}

func TestChatServer_CreateMessage_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatUcase := mock.NewMockChatUsecase(ctrl)

	chatHandler := NewChatServer(chatUcase)
	chatUcase.EXPECT().CreateMessage(&models.CreateMessageReq{
		ChatID:  uint64(1),
		Content: "dfsfsd",
	}, uint64(0)).Return(&models.MessageResp{}, nil)

	chatHandler.CreateMessage(context.Background(), &protoChat.CreateMessageReq{
		UserID:  0,
		ChatID:  1,
		Content: "dfsfsd",
	})

	//error
	chatUcase.EXPECT().CreateMessage(&models.CreateMessageReq{
		ChatID:  uint64(1),
		Content: "dfsfsd",
	}, uint64(0)).Return(nil, errors.UnexpectedInternal(sql.ErrConnDone))

	_, err := chatHandler.CreateMessage(context.Background(), &protoChat.CreateMessageReq{
		UserID:  0,
		ChatID:  1,
		Content: "dfsfsd",
	})
	assert.Error(t, err)
}

func TestChatServer_GetLastNMessages_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatUcase := mock.NewMockChatUsecase(ctrl)

	chatHandler := NewChatServer(chatUcase)

	msg := &models.MessageResp{ChatID: 1}

	chatUcase.EXPECT().GetLastNMessages(&models.GetLastNMessagesReq{
		UserID: uint64(0),
		ChatID: uint64(1),
		Count:  2,
	}).Return([]*models.MessageResp{msg}, nil)

	chatHandler.GetLastNMessages(context.Background(), &protoChat.GetLastNMessagesReq{
		UserID: 0,
		ChatID: 1,
		Count:  2,
	})

	//error
	chatUcase.EXPECT().GetLastNMessages(&models.GetLastNMessagesReq{
		UserID: uint64(0),
		ChatID: uint64(1),
		Count:  2,
	}).Return([]*models.MessageResp{msg}, errors.UnexpectedInternal(sql.ErrConnDone))

	_, err := chatHandler.GetLastNMessages(context.Background(), &protoChat.GetLastNMessagesReq{
		UserID: 0,
		ChatID: 1,
		Count:  2,
	})
	assert.Error(t, err)
}

func TestChatServer_GetNMessagesBefore_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatUcase := mock.NewMockChatUsecase(ctrl)

	chatHandler := NewChatServer(chatUcase)

	msg := &models.MessageResp{ChatID: 1}

	chatUcase.EXPECT().GetNMessagesBefore(&models.GetNMessagesBeforeReq{
		ChatID:        uint64(0),
		Count:         1,
		LastMessageID: uint64(2),
	}).Return([]*models.MessageResp{msg}, nil)

	chatHandler.GetNMessagesBefore(context.Background(), &protoChat.GetNMessagesReq{
		ChatID:        0,
		Count:         1,
		LastMessageId: 2,
	})

	//error
	chatUcase.EXPECT().GetNMessagesBefore(&models.GetNMessagesBeforeReq{
		ChatID:        uint64(0),
		Count:         1,
		LastMessageID: uint64(2),
	}).Return(nil, errors.UnexpectedInternal(sql.ErrConnDone))

	_, err := chatHandler.GetNMessagesBefore(context.Background(), &protoChat.GetNMessagesReq{
		ChatID:        0,
		Count:         1,
		LastMessageId: 2,
	})
	assert.Error(t, err)
}
