package usecase

import (
	"database/sql"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/microservices/chat/mocks"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
)

func TestChatUsecase_CreateChat_SuccessWithCreate(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatRepo := mock.NewMockChatRepository(ctrl)
	chatUcase := NewChatUsecase(chatRepo)

	chat := &models.Chat{
		PartnerID: uint64(0),
		ProductID: uint64(1),
	}
	var userID uint64 = 2

	chatRepo.EXPECT().CheckChatExist(chat, userID).Return(nil)

	_, err := chatUcase.CreateChat(&models.ChatCreateReq{
		ProductID: 1,
		PartnerID: 0,
	}, userID)

	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestChatUsecase_CreateChat_SuccessWithoutCreate(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatRepo := mock.NewMockChatRepository(ctrl)
	chatUcase := NewChatUsecase(chatRepo)

	chat := &models.Chat{
		PartnerID: uint64(0),
		ProductID: uint64(1),
	}
	var userID uint64 = 2

	chatRepo.EXPECT().CheckChatExist(chat, userID).Return(sql.ErrNoRows)
	chatRepo.EXPECT().InsertChat(chat, userID).Return(nil)

	_, err := chatUcase.CreateChat(&models.ChatCreateReq{
		ProductID: 1,
		PartnerID: 0,
	}, userID)

	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestChatUsecase_CreateChat_CreateErr(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatRepo := mock.NewMockChatRepository(ctrl)
	chatUcase := NewChatUsecase(chatRepo)

	chat := &models.Chat{
		PartnerID: uint64(0),
		ProductID: uint64(1),
	}
	var userID uint64 = 2

	chatRepo.EXPECT().CheckChatExist(chat, userID).Return(sql.ErrNoRows)
	chatRepo.EXPECT().InsertChat(chat, userID).Return(sql.ErrConnDone)

	_, err := chatUcase.CreateChat(&models.ChatCreateReq{
		ProductID: 1,
		PartnerID: 0,
	}, userID)

	assert.Equal(t, err, errors.UnexpectedInternal(sql.ErrConnDone))
}

func TestChatUsecase_CreateChat_Error(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatRepo := mock.NewMockChatRepository(ctrl)
	chatUcase := NewChatUsecase(chatRepo)

	chat := &models.Chat{
		PartnerID: uint64(0),
		ProductID: uint64(1),
	}
	var userID uint64 = 2

	chatRepo.EXPECT().CheckChatExist(chat, userID).Return(sql.ErrConnDone)

	_, err := chatUcase.CreateChat(&models.ChatCreateReq{
		ProductID: 1,
		PartnerID: 0,
	}, userID)

	assert.Equal(t, err, errors.UnexpectedInternal(sql.ErrConnDone))
}

func TestChatUsecase_GetChatById_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatRepo := mock.NewMockChatRepository(ctrl)
	chatUcase := NewChatUsecase(chatRepo)

	var userID uint64 = 0
	var chatID uint64 = 1

	chatRepo.EXPECT().GetChatById(chatID, userID).Return(&models.Chat{}, nil)

	_, err := chatUcase.GetChatById(chatID, userID)

	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestChatUsecase_GetChatById_NoChat(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatRepo := mock.NewMockChatRepository(ctrl)
	chatUcase := NewChatUsecase(chatRepo)

	var userID uint64 = 0
	var chatID uint64 = 1

	chatRepo.EXPECT().GetChatById(chatID, userID).Return(nil, sql.ErrNoRows)

	_, err := chatUcase.GetChatById(chatID, userID)

	assert.Equal(t, err, errors.Cause(errors.ChatNotExist))
}

func TestChatUsecase_GetChatById_Error(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatRepo := mock.NewMockChatRepository(ctrl)
	chatUcase := NewChatUsecase(chatRepo)

	var userID uint64 = 0
	var chatID uint64 = 1

	chatRepo.EXPECT().GetChatById(chatID, userID).Return(nil, sql.ErrConnDone)

	_, err := chatUcase.GetChatById(chatID, userID)

	assert.Equal(t, err, errors.UnexpectedInternal(sql.ErrConnDone))
}

func TestChatUsecase_GetUserChats_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatRepo := mock.NewMockChatRepository(ctrl)
	chatUcase := NewChatUsecase(chatRepo)

	var userID uint64 = 0

	chat := &models.Chat{PartnerID: 0}

	chatRepo.EXPECT().GetUserChats(userID).Return([]*models.Chat{chat}, nil)

	_, err := chatUcase.GetUserChats(userID)

	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestChatUsecase_GetUserChats_Error(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatRepo := mock.NewMockChatRepository(ctrl)
	chatUcase := NewChatUsecase(chatRepo)

	var userID uint64 = 0

	chatRepo.EXPECT().GetUserChats(userID).Return(nil, sql.ErrConnDone)

	_, err := chatUcase.GetUserChats(userID)

	assert.Equal(t, err, errors.UnexpectedInternal(sql.ErrConnDone))
}

func TestChatUsecase_CreateMessage_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatRepo := mock.NewMockChatRepository(ctrl)
	chatUcase := NewChatUsecase(chatRepo)

	var userID uint64 = 0

	req := &models.CreateMessageReq{
		ChatID:  1,
		Content: "ggg",
	}

	chatRepo.EXPECT().InsertMessage(req, userID).Return(&models.Message{}, nil)

	_, err := chatUcase.CreateMessage(req, userID)

	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestChatUsecase_CreateMessage_Error(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatRepo := mock.NewMockChatRepository(ctrl)
	chatUcase := NewChatUsecase(chatRepo)

	var userID uint64 = 0

	req := &models.CreateMessageReq{
		ChatID:  1,
		Content: "ggg",
	}

	chatRepo.EXPECT().InsertMessage(req, userID).Return(nil, sql.ErrConnDone)

	_, err := chatUcase.CreateMessage(req, userID)

	assert.Equal(t, err, errors.UnexpectedInternal(sql.ErrConnDone))
}

func TestChatUsecase_GetLastNMessages_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatRepo := mock.NewMockChatRepository(ctrl)
	chatUcase := NewChatUsecase(chatRepo)

	req := &models.GetLastNMessagesReq{
		UserID: 0,
		ChatID: 1,
		Count:  2,
	}

	msg := &models.Message{ChatID: 1}

	chatRepo.EXPECT().GetLastNMessages(req).Return([]*models.Message{msg}, nil)

	_, err := chatUcase.GetLastNMessages(req)

	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestChatUsecase_GetLastNMessages_Error(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatRepo := mock.NewMockChatRepository(ctrl)
	chatUcase := NewChatUsecase(chatRepo)

	req := &models.GetLastNMessagesReq{
		UserID: 0,
		ChatID: 1,
		Count:  2,
	}

	chatRepo.EXPECT().GetLastNMessages(req).Return(nil, sql.ErrConnDone)

	_, err := chatUcase.GetLastNMessages(req)

	assert.Equal(t, err, errors.UnexpectedInternal(sql.ErrConnDone))
}

func TestChatUsecase_GetNMessagesBefore_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatRepo := mock.NewMockChatRepository(ctrl)
	chatUcase := NewChatUsecase(chatRepo)

	req := &models.GetNMessagesBeforeReq{
		ChatID:        0,
		Count:         12,
		LastMessageID: 1,
	}

	msg := &models.Message{ChatID: 1}

	chatRepo.EXPECT().GetNMessagesBefore(req).Return([]*models.Message{msg}, nil)

	_, err := chatUcase.GetNMessagesBefore(req)

	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestChatUsecase_GetNMessagesBefore_Error(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatRepo := mock.NewMockChatRepository(ctrl)
	chatUcase := NewChatUsecase(chatRepo)

	req := &models.GetNMessagesBeforeReq{
		ChatID:        0,
		Count:         12,
		LastMessageID: 1,
	}

	chatRepo.EXPECT().GetNMessagesBefore(req).Return(nil, sql.ErrConnDone)

	_, err := chatUcase.GetNMessagesBefore(req)

	assert.Equal(t, err, errors.UnexpectedInternal(sql.ErrConnDone))
}
