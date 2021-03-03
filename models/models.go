package models

type SignInData struct {
	Telephone   string `json:"telephone"`
	Password    string `json:"password"`
	IsLoggedIn  bool   `json:"is_logged_in"`
}

type SignUpData struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Surname     string   `json:"surname"`
	Sex         string   `json:"sex"`
	Email       string   `json:"email"`
	Telephone   string   `json:"telephone"`
	Password    string   `json:"password"`
	DateBirth   int      `json:"date_birth"`
	Day         string   `json:"day"`
	Month       string   `json:"month"`
	Year        string   `json:"year"`
	LinkImages  []string `json:"linkImages"`
}

type ProductData struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Date        int		 `json:"date"`
	Amount      int      `json:"amount"`
	LinkImages  []string `json:"linkImages"`
	Description string   `json:"description"`
	OwnerID     int		 `json:"owner_id"`
	Views       int		 `json:"views"`
	Likes       int		 `json:"likes"`
}