package chat

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
)

type ChatRepository interface {
	InsertChat(chat *models.Chat, userID uint64) error
	GetChatById(chatId uint64, userID uint64) (*models.Chat, error)
	GetUserChats(userId uint64) ([]*models.Chat, error)
	InsertMessage(req *models.CreateMessageReq, userID uint64) (*models.Message, error)
	GetLastNMessages(req *models.GetLastNMessagesReq) ([]*models.Message, error)
	GetNMessagesBefore(req *models.GetNMessagesBeforeReq) ([]*models.Message, error)
}
