package models

import "time"

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
	ChatID uint64	`json:"chat_id"`
	Count int `json:"messages"`
}

type GetNMessagesBeforeReq struct {
	ChatID uint64  `json:"chat_id"`
	Count int `json:"n_messages"`
	LastMessageID uint64  `json:"message_id"`
}

type MessageResp struct {
	ID uint64  `json:"id"`
	Content string  `json:"content"`
	CreationTime time.Time  `json:"time"`
	ChatID uint64 `json:"chat_id"`
	UserID uint64 `json:"user_id"`
}

