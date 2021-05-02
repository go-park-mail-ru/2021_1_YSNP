package models

import (
	proto "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/proto/chat"
	"github.com/golang/protobuf/ptypes"
	"time"
)

type Chat struct {
	ID uint64

	CreationTime time.Time
	LastMsgID uint64
	LastMsgContent string
	LastMsgTime time.Time

	PartnerID uint64
	PartnerName string
	PartnerSurname string
	PartnerAvatarLink string

	ProductID uint64
	ProductName string
	ProductAmount int
	ProductAvatarLink string

	LastReadMsgId uint64
	NewMessages int
}



type ChatCreateReq struct {
	ProductID uint64  `json:"productID"`
	PartnerID uint64 `json:"partnerID"`
}



type ChatResponse struct {
	ID uint64 `json:"id"`
	CreationTime time.Time `json:"creation_time"`
	LastMsgContent string `json:"last_msg_content"`
	LastMsgTime time.Time  `json:"last_msg_time"`

	PartnerName string `json:"partner_name"`
	PartnerSurname string `json:"partner_surname"`
	PartnerAvatarLink string `json:"partner_avatar"`

	ProductName string `json:"product_name"`
	ProductAmount int `json:"product_amount"`
	ProductAvatarLink string `json:"product_avatar_link"`

	NewMessages int `json:"new_messages"`
}

type Message struct {
	ID uint64
	Content string
	CreationTime time.Time
	ChatID uint64
	UserID uint64
}

type CreateMessageReq struct {
	ChatID uint64  `json:"chat_id"`
	Content string  `json:"content"`
}

type GetLastNMessagesReq struct {
	UserID uint64 `json:"-"`
	ChatID uint64	`json:"chat_id"`
	Count int `json:"count"`
}

type GetNMessagesBeforeReq struct {
	ChatID uint64  `json:"chat_id"`
	Count int `json:"count"`
	LastMessageID uint64  `json:"message_id"`
}

type MessageResp struct {
	ID uint64  `json:"id"`
	Content string  `json:"content"`
	CreationTime time.Time  `json:"time"`
	ChatID uint64 `json:"chat_id"`
	UserID uint64 `json:"user_id"`
}

func ModelChatRespToGRPC(chatModel *ChatResponse) *proto.ChatResp {
	creationTime, _ := ptypes.TimestampProto(chatModel.CreationTime)
	lastMsgTime, _ := ptypes.TimestampProto(chatModel.LastMsgTime)

	return &proto.ChatResp{
		ID:                int64(chatModel.ID),
		CreationTime:      creationTime,
		LastMsgContent:    chatModel.LastMsgContent,
		LastMsgTime:       lastMsgTime,
		PartnerName:       chatModel.PartnerName,
		PartnerSurname:    chatModel.PartnerSurname,
		PartnerAvatarLink: chatModel.PartnerAvatarLink,
		ProductName:       chatModel.ProductName,
		ProductAmount: int32(chatModel.ProductAmount),
		ProductAvatarLink: chatModel.ProductAvatarLink,
		NewMessages:       int32(chatModel.NewMessages),
	}
}

func ModelMessageRespToGRPC(msgResp *MessageResp) *proto.MessageResp {
	creationTime, _ := ptypes.TimestampProto(msgResp.CreationTime)

	return &proto.MessageResp{
		ID:           int64(msgResp.ID),
		Content:      msgResp.Content,
		CreationTime: creationTime,
		ChatID:       int64(msgResp.ChatID),
		UserID:       int64(msgResp.UserID),
	}
}
