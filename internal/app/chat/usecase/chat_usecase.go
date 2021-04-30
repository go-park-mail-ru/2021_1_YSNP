package usecase

import (
	"context"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/proto/chat"
	"google.golang.org/grpc"
)

type ChatClient struct {
	client chat.ChatClient
}

func NewChatClient (conn grpc.ClientConnInterface) *ChatClient {
	return &ChatClient{
		client: chat.NewChatClient(conn),
	}
}

func (cc *ChatClient) CreateChat(req *models.ChatCreateReq, userID uint64) (*models.ChatResponse, *errors.Error) {
	resp, err := cc.client.CreateChat(context.Background(), &chat.ChatCreateReq{
		UserID:    int64(userID),
		PartnerID: int64(req.PartnerID),
		ProductID: int64(req.ProductID),
	})

	if err != nil {
		return nil, errors.GRPCError(err)
	}

	return &models.ChatResponse{
		ID:                uint64(resp.ID),
		CreationTime:      resp.CreationTime.AsTime(),
		LastMsgContent:    resp.LastMsgContent,
		LastMsgTime:       resp.LastMsgTime.AsTime(),
		PartnerName:       resp.PartnerName,
		PartnerSurname:    resp.PartnerSurname,
		PartnerAvatarLink: resp.PartnerAvatarLink,
		ProductName: 	   resp.ProductName,
		ProductAvatarLink: resp.ProductAvatarLink,
		NewMessages:       int(resp.NewMessages),
	}, nil
}

func (cc *ChatClient) GetChatById(chatID uint64, userID uint64) (*models.ChatResponse, *errors.Error) {
	resp, err := cc.client.GetChatById(context.Background(), &chat.GetChatByIDReq{
		UserID: int64(userID),
		ChatID: int64(chatID),
	})

	if err != nil {
		return nil, errors.GRPCError(err)
	}

	return &models.ChatResponse{
		ID:                uint64(resp.ID),
		CreationTime:      resp.CreationTime.AsTime(),
		LastMsgContent:    resp.LastMsgContent,
		LastMsgTime:       resp.LastMsgTime.AsTime(),
		PartnerName:       resp.PartnerName,
		PartnerSurname:    resp.PartnerSurname,
		PartnerAvatarLink: resp.PartnerAvatarLink,
		ProductName: 	   resp.ProductName,
		ProductAvatarLink: resp.ProductAvatarLink,
		NewMessages:       int(resp.NewMessages),
	}, nil
}

func (cc *ChatClient) GetUserChats(userID uint64) ([]*models.ChatResponse, *errors.Error) {
	resps, err := cc.client.GetUserChats(context.Background(), &chat.UserID{UserID: int64(userID)})
	if err != nil {
		return nil, errors.GRPCError(err)
	}

	chatResps := []*models.ChatResponse{}
	for _, resp := range resps.Chats {
		chatResps = append(chatResps, &models.ChatResponse{
			ID:                uint64(resp.ID),
			CreationTime:      resp.CreationTime.AsTime(),
			LastMsgContent:    resp.LastMsgContent,
			LastMsgTime:       resp.LastMsgTime.AsTime(),
			PartnerName:       resp.PartnerName,
			PartnerSurname:    resp.PartnerSurname,
			PartnerAvatarLink: resp.PartnerAvatarLink,
			ProductName: 	   resp.ProductName,
			ProductAvatarLink: resp.ProductAvatarLink,
			NewMessages:       int(resp.NewMessages),
		})
	}

	return chatResps, nil
}

func (cc *ChatClient) CreateMessage(req *models.CreateMessageReq, userID uint64) (*models.MessageResp, *errors.Error){
	resp, err := cc.client.CreateMessage(context.Background(), &chat.CreateMessageReq{
		UserID:  int64(userID),
		ChatID:  int64(req.ChatID),
		Content: req.Content,
	})

	if err != nil {
		return nil, errors.GRPCError(err)
	}
	return &models.MessageResp{
		ID:           uint64(resp.ID),
		Content:      resp.Content,
		CreationTime: resp.CreationTime.AsTime(),
		ChatID:       uint64(resp.ChatID),
		UserID: 	  uint64(resp.UserID),
	}, nil

}

func (cc *ChatClient) GetLastNMessages(req *models.GetLastNMessagesReq) ([]*models.MessageResp, *errors.Error) {
	resps, err := cc.client.GetLastNMessages(context.Background(), &chat.GetLastNMessagesReq{
		ChatID: int64(req.ChatID),
		Count:  int32(req.Count),
	})

	if err != nil {
		return nil, errors.GRPCError(err)
	}

	msgResps := []*models.MessageResp{}
	for _, resp := range resps.Messages {
		msgResps = append(msgResps, &models.MessageResp{
			ID:           uint64(resp.ID),
			Content:      resp.Content,
			CreationTime: resp.CreationTime.AsTime(),
			ChatID:       uint64(resp.ChatID),
			UserID: 	  uint64(resp.UserID),
		})
	}

	return msgResps, nil
}

func (cc *ChatClient) GetNMessagesBefore(req *models.GetNMessagesBeforeReq) ([]*models.MessageResp, *errors.Error) {
	resps, err := cc.client.GetNMessagesBefore(context.Background(), &chat.GetNMessagesReq{
		ChatID:        int64(req.ChatID),
		Count:         int32(req.Count),
		LastMessageId: int64(req.LastMessageID),
	})
	if err != nil {
		return nil, errors.GRPCError(err)
	}

	msgResps := []*models.MessageResp{}
	for _, resp := range resps.Messages {
		msgResps = append(msgResps, &models.MessageResp{
			ID:           uint64(resp.ID),
			Content:      resp.Content,
			CreationTime: resp.CreationTime.AsTime(),
			ChatID:       uint64(resp.ChatID),
			UserID: 	  uint64(resp.UserID),
		})
	}

	return msgResps, nil
}