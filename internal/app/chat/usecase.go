package chat

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
)

type ChatUsecase interface {
	CreateChat(req *models.ChatCreateReq, userID int) (*models.ChatResponse, *errors.Error)
	GetChatById(chatID int, userId int) (*models.ChatResponse, *errors.Error)
	GetUserChats(userID int) ([]*models.ChatResponse, *errors.Error)
	CreateMessage(req *models.CreateMessageReq, userID int) (*models.MessageResp, *errors.Error)
	GetLastNMessages(req *models.GetLastNMessagesReq) ([]*models.MessageResp, *errors.Error)
	GetNMessagesBefore(req *models.GetNMessagesBeforeReq) ([]*models.MessageResp, *errors.Error)
}
