package models

type LoginData struct {
	Telephone string `json:"telephone"`
	Password  string `json:"password"`
}

type SignUpData struct {
	ID         uint64   `json:"id"`
	Name       string   `json:"name"`
	Surname    string   `json:"surname"`
	Sex        string   `json:"sex"`
	Email      string   `json:"email"`
	Telephone  string   `json:"telephone"`
	Password   string   `json:"password,omitempty"`
	DateBirth  string   `json:"dateBirth"`
	LinkImages []string `json:"linkImages"`
}

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

type Error struct {
	Message string `json:"message"`
}

type PasswordChange struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}
