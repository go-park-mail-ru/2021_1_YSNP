package models

import "time"

type Chat struct {
	ID int

	CreationTime time.Time
	LastMsgID int
	LastMsgContent string
	LastMsgTime time.Time

	PartnerID int
	PartnerName string
	PartnerSurname string
	PartnerAvatarLink string

	ProductID int
	ProductName string
	ProductAvatarLink string

	LastReadMsgId int
	NewMessages int
}



type ChatCreateReq struct {
	ProductID int  `json:"productID"`
	PartnerID int `json:"partnerID"`
}



type ChatResponse struct {
	ID int `json:"id"`
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
	ID int
	Content string
	CreationTime time.Time
	ChatID int
	UserID int
}

type CreateMessageReq struct {
	ChatID int  `json:"chat_id"`
	Content string  `json:"content"`
}

type GetLastNMessagesReq struct {
	ChatID int	`json:"chat_id"`
	Count int `json:"messages"`
}

type GetNMessagesBeforeReq struct {
	ChatID int  `json:"chat_id"`
	Count int `json:"n_messages"`
	LastMessageID int  `json:"message_id"`
}

type MessageResp struct {
	ID int  `json:"id"`
	Content string  `json:"content"`
	CreationTime time.Time  `json:"time"`
	ChatID int `json:"chat_id"`
	UserID int `json:"user_id"`
}

