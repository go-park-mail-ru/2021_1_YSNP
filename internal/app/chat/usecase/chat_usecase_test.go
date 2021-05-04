package usecase

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	mock "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/chat/mocks"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/proto/chat"
)

func TestChatClient_CreateChat(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	chatClient := mock.NewMockChatClient(ctrl)
	cl := &ChatClient{client: chatClient}

	chatClient.EXPECT().CreateChat(context.Background(), &chat.ChatCreateReq{
		UserID:    int64(0),
		PartnerID: int64(1),
		ProductID: int64(2),
	}).Return(&chat.ChatResp{}, nil)
	_, err := cl.CreateChat(&models.ChatCreateReq{
		ProductID: 2,
		PartnerID: 1,
	}, uint64(0))

	assert.Equal(t, err, (*errors.Error)(nil))

	//error
	chatClient.EXPECT().CreateChat(context.Background(), &chat.ChatCreateReq{
		UserID:    int64(0),
		PartnerID: int64(1),
		ProductID: int64(2),
	}).Return(nil, grpc.ErrClientConnClosing)
	_, err = cl.CreateChat(&models.ChatCreateReq{
		ProductID: 2,
		PartnerID: 1,
	}, uint64(0))

	assert.Equal(t, err, errors.GRPCError(grpc.ErrClientConnClosing))
}

func TestChatClient_GetChatById(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	chatClient := mock.NewMockChatClient(ctrl)
	cl := &ChatClient{client: chatClient}

	chatClient.EXPECT().GetChatByID(context.Background(), &chat.GetChatByIDReq{
		UserID: int64(1),
		ChatID: int64(2),
	}).Return(&chat.ChatResp{}, nil)

	_, err := cl.GetChatById(2, 1)
	assert.Equal(t, err, (*errors.Error)(nil))

	//error
	chatClient.EXPECT().GetChatByID(context.Background(), &chat.GetChatByIDReq{
		UserID: int64(1),
		ChatID: int64(2),
	}).Return(nil, grpc.ErrClientConnClosing)

	_, err = cl.GetChatById(2, 1)
	assert.Equal(t, err, errors.GRPCError(grpc.ErrClientConnClosing))
}

func TestChatClient_GetUserChats(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	chatClient := mock.NewMockChatClient(ctrl)
	cl := &ChatClient{client: chatClient}

	chatClient.EXPECT().GetUserChats(context.Background(), &chat.UserID{UserID: int64(1)}).Return(&chat.ChatRespArray{}, nil)

	_, err := cl.GetUserChats(1)
	assert.Equal(t, err, (*errors.Error)(nil))

	//error
	chatClient.EXPECT().GetUserChats(context.Background(), &chat.UserID{UserID: int64(1)}).Return(nil, grpc.ErrClientConnClosing)

	_, err = cl.GetUserChats(1)
	assert.Equal(t, err, errors.GRPCError(grpc.ErrClientConnClosing))
}

func TestChatClient_CreateMessage(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	chatClient := mock.NewMockChatClient(ctrl)
	cl := &ChatClient{client: chatClient}

	chatClient.EXPECT().CreateMessage(context.Background(), &chat.CreateMessageReq{
		UserID:  int64(1),
		ChatID:  int64(2),
		Content: "fhdjhg",
	}).Return(&chat.MessageResp{}, nil)

	_, err := cl.CreateMessage(&models.CreateMessageReq{
		ChatID:  2,
		Content: "fhdjhg",
	}, 1)
	assert.Equal(t, err, (*errors.Error)(nil))

	//error
	chatClient.EXPECT().CreateMessage(context.Background(), &chat.CreateMessageReq{
		UserID:  int64(1),
		ChatID:  int64(2),
		Content: "fhdjhg",
	}).Return(nil, grpc.ErrClientConnClosing)

	_, err = cl.CreateMessage(&models.CreateMessageReq{
		ChatID:  2,
		Content: "fhdjhg",
	}, 1)
	assert.Equal(t, err, errors.GRPCError(grpc.ErrClientConnClosing))
}

func TestChatClient_GetLastNMessages(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	chatClient := mock.NewMockChatClient(ctrl)
	cl := &ChatClient{client: chatClient}

	chatClient.EXPECT().GetLastNMessages(context.Background(), &chat.GetLastNMessagesReq{
		UserID: int64(1),
		ChatID: int64(2),
		Count:  int32(3),
	}).Return(&chat.MessageRespArray{}, nil)

	_, err := cl.GetLastNMessages(&models.GetLastNMessagesReq{
		UserID: 1,
		ChatID: 2,
		Count:  3,
	})
	assert.Equal(t, err, (*errors.Error)(nil))

	//error
	chatClient.EXPECT().GetLastNMessages(context.Background(), &chat.GetLastNMessagesReq{
		UserID: int64(1),
		ChatID: int64(2),
		Count:  int32(3),
	}).Return(nil, grpc.ErrClientConnClosing)

	_, err = cl.GetLastNMessages(&models.GetLastNMessagesReq{
		UserID: 1,
		ChatID: 2,
		Count:  3,
	})
	assert.Equal(t, err, errors.GRPCError(grpc.ErrClientConnClosing))
}

func TestChatClient_GetNMessagesBefore(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	chatClient := mock.NewMockChatClient(ctrl)
	cl := &ChatClient{client: chatClient}

	chatClient.EXPECT().GetNMessagesBefore(context.Background(), &chat.GetNMessagesReq{
		ChatID:        int64(1),
		Count:         int32(2),
		LastMessageId: int64(3),
	}).Return(&chat.MessageRespArray{}, nil)

	_, err := cl.GetNMessagesBefore(&models.GetNMessagesBeforeReq{
		ChatID:        1,
		Count:         2,
		LastMessageID: 3,
	})
	assert.Equal(t, err, (*errors.Error)(nil))

	//error
	chatClient.EXPECT().GetNMessagesBefore(context.Background(), &chat.GetNMessagesReq{
		ChatID:        int64(1),
		Count:         int32(2),
		LastMessageId: int64(3),
	}).Return(nil, grpc.ErrClientConnClosing)

	_, err = cl.GetNMessagesBefore(&models.GetNMessagesBeforeReq{
		ChatID:        1,
		Count:         2,
		LastMessageID: 3,
	})
	assert.Equal(t, err, errors.GRPCError(grpc.ErrClientConnClosing))
}
