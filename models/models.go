package models

type LoginData struct {
	Telephone  string `json:"telephone"`
	Password   string `json:"password"`
	IsLoggedIn bool   `json:"is_logged_in"`
}

type SignUpData struct {
	ID         uint64   `json:"id"`
	Name       string   `json:"name"`
	Surname    string   `json:"surname"`
	Sex        string   `json:"sex"`
	Email      string   `json:"email"`
	Telephone  string   `json:"telephone"`
	Password   string   `json:"password"`
	DateBirth  string   `json:"date_birth"`
	LinkImages []string `json:"linkImages"`
}

type ProductData struct {
	ID           uint64   `json:"id"`
	Name         string   `json:"name"`
	Date         string   `json:"date"`
	Amount       int      `json:"amount"`
	LinkImages   []string `json:"link_images"`
	Description  string   `json:"description"`
	OwnerID      int      `json:"owner_id"`
	OwnerName    string   `json:"owner_name"`
	OwnerSurname string   `json:"owner_surname"`
	Views        int      `json:"views"`
	Likes        int      `json:"likes"`
}

type ProductListData struct {
	ID         uint64   `json:"id"`
	Name       string   `json:"name"`
	Date       string   `json:"date"`
	Amount     int      `json:"amount"`
	LinkImages []string `json:"link_images"`
}

type Error struct {
	Message string `json:"message"`
}
