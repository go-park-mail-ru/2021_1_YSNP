package models

import "time"

//easyjson:json
type ProductData struct {
	ID              uint64   `json:"id" valid:"numeric"`
	Name            string   `json:"name" valid:"stringlength(1|100)"`
	Date            string   `json:"date" valid:"-"`
	Amount          int      `json:"amount" valid:"numeric"`
	LinkImages      []string `json:"linkImages" valid:"stringArray"`
	Description     string   `json:"description" valid:"stringlength(10|4000)"`
	Category        string   `json:"category" valid:"type(string)"`
	Address         string   `json:"address" valid:"type(string)"`
	Longitude       float64  `json:"longitude" valid:"longitude"`
	Latitude        float64  `json:"latitude" valid:"latitude"`
	Views           int      `json:"views" valid:"numeric"`
	Likes           int      `json:"likes" valid:"numeric"`
	Tariff          int      `json:"tariff" valid:"numeric"`
	OwnerID         uint64   `json:"ownerId" valid:"numeric"`
	OwnerName       string   `json:"ownerName" valid:"stringlength(1|30)"`
	OwnerSurname    string   `json:"ownerSurname" valid:"stringlength(1|30)"`
	OwnerLinkImages string   `json:"ownerLinkImages" valid:"type(string)"`
	OwnerRating     float64  `json:"owner_rating"`
	Close           bool     `json:"close" valid:"type(bool)"`
}

//easyjson:json
type ProductListData struct {
	ID         uint64   `json:"id" valid:"numeric"`
	Name       string   `json:"name" valid:"stringlength(1|100)"`
	Date       string   `json:"date" valid:"-"`
	Amount     int      `json:"amount" valid:"numeric"`
	LinkImages []string `json:"linkImages" valid:"stringArray"`
	UserLiked  bool     `json:"userLiked" valid:"type(bool)"`
	Tariff     int      `json:"tariff" valid:"numeric"`
	Close      bool     `json:"close" valid:"type(bool)"`
}

type Page struct {
	From  uint64 `valid:"numeric"`
	Count uint64 `valid:"numeric"`
}

type PageWithSort struct {
	From  uint64 `valid:"numeric"`
	Count uint64 `valid:"numeric"`
	Sort  string `valid:"stringlength(3|5)"`
}

//easyjson:json
type Category struct {
	Title string `json:"title"`
}

//easyjson:json
type Review struct {
	ID             uint64    `json:"id"`
	Content        string    `json:"content"`
	Rating         float32   `json:"rating"`
	ReviewerID     uint64    `json:"reviewer_id"`
	ReviewerName   string    `json:"reviewer_name"`
	ReviewerAvatar string    `json:"reviewer_avatar"`
	ProductID      uint64    `json:"product_id"`
	ProductName    string    `json:"product_name"`
	ProductImage   string    `json:"product_image"`
	TargetID       uint64    `json:"target_id"`
	Type           string    `json:"type"`
	CreationTime   time.Time `json:"creation_time"`
}

//easyjson:json
type WaitingReview struct {
	ProductID    uint64 `json:"product_id"`
	ProductName  string `json:"product_name"`
	ProductImage string `json:"product_image"`
	TargetID     uint64 `json:"target_id"`
	TargetName   string `json:"target_name"`
	TargetAvatar string `json:"target_avatar"`
	Type         string `json:"type"`
}

type SetProductBuyerRequest struct {
	Buyer_id uint64 `json:"buyer_id"`
}
