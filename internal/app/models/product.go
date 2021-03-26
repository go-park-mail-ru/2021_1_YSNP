package models

const Url = "http://89.208.199.170:8080"

//const Url = "http://localhost:8080"

type ProductData struct {
	ID           uint64   `json:"id"`
	Name         string   `json:"name"`
	Date         string   `json:"date"`
	Amount       int      `json:"amount"`
	LinkImages   []string `json:"linkImages"`
	Description  string   `json:"description"`
	Category     string   `json:"category"`
	OwnerID      uint64   `json:"ownerId"`
	OwnerName    string   `json:"ownerName"`
	OwnerSurname string   `json:"ownerSurname"`
	Views        int      `json:"views"`
	Likes        int      `json:"likes"`
}

type ProductListData struct {
	ID         uint64   `json:"id"`
	Name       string   `json:"name"`
	Date       string   `json:"date"`
	Amount     int      `json:"amount"`
	LinkImages []string `json:"linkImages"`
}

type Content struct {
	From uint64 `json:"from"`
	Count uint64 `json:"count"`
}

type OrderType struct {
	Main bool `json:"main"`
}

type Page struct {
	Content Content `json:"content"`
	Order OrderType `order:"order"`
}