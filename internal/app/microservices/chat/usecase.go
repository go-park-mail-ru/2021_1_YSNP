package chat

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
)

//go:generate mockgen -destination=./mocks/mock_chat_ucase.go -package=mock github.com/go-park-mail-ru/2021_1_YSNP/internal/app/microservices/chat  ChatUsecase

type ChatUsecase interface {
	CreateChat(req *models.ChatCreateReq, userID uint64) (*models.ChatResponse, *errors.Error)
	GetChatById(chatID uint64, userID uint64) (*models.ChatResponse, *errors.Error)
	GetUserChats(userID uint64) ([]*models.ChatResponse, *errors.Error)
	CreateMessage(req *models.CreateMessageReq, userID uint64) (*models.MessageResp, *errors.Error)
	GetLastNMessages(req *models.GetLastNMessagesReq) ([]*models.MessageResp, *errors.Error)
	GetNMessagesBefore(req *models.GetNMessagesBeforeReq) ([]*models.MessageResp, *errors.Error)
}