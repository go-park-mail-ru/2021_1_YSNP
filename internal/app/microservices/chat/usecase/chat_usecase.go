package usecase

import (
	"database/sql"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/microservices/chat"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
)

type ChatUsecase struct {
	chatRepo chat.ChatRepository
}

func NewChatUsecase(repo chat.ChatRepository) chat.ChatUsecase {
	return &ChatUsecase{
		chatRepo: repo,
	}
}

func (c *ChatUsecase) CreateChat(req *models.ChatCreateReq, userID uint64) (*models.ChatResponse, *errors.Error) {
	chat := &models.Chat{
		PartnerID: req.PartnerID,
		ProductID: req.ProductID,
	}
	err := c.chatRepo.CheckChatExist(chat, userID)
	switch {
	case err == sql.ErrNoRows:
		err := c.chatRepo.InsertChat(chat, userID)
		if err != nil {
			return nil, errors.UnexpectedInternal(err)
		}
	case err != nil:
		return nil, errors.UnexpectedInternal(err)
	}

	return &models.ChatResponse{
		ID:                chat.ID,
		CreationTime:      chat.CreationTime,
		LastMsgContent:    chat.LastMsgContent,
		LastMsgTime:       chat.LastMsgTime,
		PartnerID:         chat.PartnerID,
		PartnerName:       chat.PartnerName,
		PartnerSurname:    chat.PartnerSurname,
		PartnerAvatarLink: chat.PartnerAvatarLink,
		ProductID:         chat.ProductID,
		ProductName:       chat.ProductName,
		ProductAmount:     chat.ProductAmount,
		ProductAvatarLink: chat.ProductAvatarLink,
		NewMessages:       chat.NewMessages,
	}, nil
}

func (c *ChatUsecase) GetChatById(chatID uint64, userID uint64) (*models.ChatResponse, *errors.Error) {
	chat, err := c.chatRepo.GetChatById(chatID, userID)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.Cause(errors.ChatNotExist)
	case err != nil:
		return nil, errors.UnexpectedInternal(err)
	}

	return &models.ChatResponse{
		ID:                chat.ID,
		CreationTime:      chat.CreationTime,
		LastMsgContent:    chat.LastMsgContent,
		LastMsgTime:       chat.LastMsgTime,
		PartnerID:         chat.PartnerID,
		PartnerName:       chat.PartnerName,
		PartnerSurname:    chat.PartnerSurname,
		PartnerAvatarLink: chat.PartnerAvatarLink,
		ProductID:         chat.ProductID,
		ProductName:       chat.ProductName,
		ProductAmount:     chat.ProductAmount,
		ProductAvatarLink: chat.ProductAvatarLink,
		NewMessages:       chat.NewMessages,
	}, nil
}

func (c *ChatUsecase) GetUserChats(userID uint64) ([]*models.ChatResponse, *errors.Error) {
	chats, err := c.chatRepo.GetUserChats(userID)
	if err != nil {
		return nil, errors.UnexpectedInternal(err)
	}

	chatResp := []*models.ChatResponse{}
	for _, chat := range chats {
		chatResp = append(chatResp, &models.ChatResponse{
			ID:                chat.ID,
			CreationTime:      chat.CreationTime,
			LastMsgContent:    chat.LastMsgContent,
			LastMsgTime:       chat.LastMsgTime,
			PartnerID:         chat.PartnerID,
			PartnerName:       chat.PartnerName,
			PartnerSurname:    chat.PartnerSurname,
			PartnerAvatarLink: chat.PartnerAvatarLink,
			ProductID:         chat.ProductID,
			ProductName:       chat.ProductName,
			ProductAmount:     chat.ProductAmount,
			ProductAvatarLink: chat.ProductAvatarLink,
			NewMessages:       chat.NewMessages,
		})
	}

	return chatResp, nil
}

func (c *ChatUsecase) CreateMessage(req *models.CreateMessageReq, userID uint64) (*models.MessageResp, *errors.Error) {
	msg, err := c.chatRepo.InsertMessage(req, userID)
	if err != nil {
		return nil, errors.UnexpectedInternal(err)
	}

	return &models.MessageResp{
		ID:           msg.ID,
		Content:      msg.Content,
		CreationTime: msg.CreationTime,
		ChatID:       msg.ChatID,
		UserID:       msg.UserID,
	}, nil
}

func (c *ChatUsecase) GetLastNMessages(req *models.GetLastNMessagesReq) ([]*models.MessageResp, *errors.Error) {
	msgs, err := c.chatRepo.GetLastNMessages(req)
	if err != nil {
		return nil, errors.UnexpectedInternal(err)
	}

	msgsResp := []*models.MessageResp{}
	for _, msg := range msgs {
		msgsResp = append(msgsResp, &models.MessageResp{
			ID:           msg.ID,
			Content:      msg.Content,
			CreationTime: msg.CreationTime,
			ChatID:       msg.ChatID,
			UserID:       msg.UserID,
		})
	}

	return msgsResp, nil
}

func (c *ChatUsecase) GetNMessagesBefore(req *models.GetNMessagesBeforeReq) ([]*models.MessageResp, *errors.Error) {
	msgs, err := c.chatRepo.GetNMessagesBefore(req)
	if err != nil {
		return nil, errors.UnexpectedInternal(err)
	}

	msgsResp := []*models.MessageResp{}
	for _, msg := range msgs {
		msgsResp = append(msgsResp, &models.MessageResp{
			ID:           msg.ID,
			Content:      msg.Content,
			CreationTime: msg.CreationTime,
			ChatID:       msg.ChatID,
			UserID:       msg.UserID,
		})
	}

	return msgsResp, nil
}
