package chat

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
)

type ChatRepository interface {
	InsertChat(chat *models.Chat, userID int) error
	GetChatById(chatId int, userID int) (*models.Chat, error)
	GetUserChats(userId int) ([]*models.Chat, error)
	InsertMessage(req *models.CreateMessageReq, userId int) (*models.Message, error)
	GetLastNMessages(req *models.GetLastNMessagesReq) ([]*models.Message, error)
	GetNMessagesBefore(req *models.GetNMessagesBeforeReq) ([]*models.Message, error)
}
