package models

type ProductData struct {
	ID           uint64   `json:"id" valid:"numeric"`
	Name         string   `json:"name" valid:"stringlength(1|100)"`
	Date         string   `json:"date" valid:"-"`
	Amount       int      `json:"amount" valid:"numeric"`
	LinkImages   []string `json:"linkImages" valid:"stringArray"`
	Description  string   `json:"description" valid:"stringlength(10|4000)"`
	Category     string   `json:"category" valid:"type(string)"`
	Address      string   `json:"address" valid:"type(string)"`
	Longitude    float64   `json:"longitude" valid:"longitude"`
	Latitude     float64   `json:"latitude" valid:"latitude"`
	OwnerID      uint64   `json:"ownerId" valid:"numeric"`
	OwnerName    string   `json:"ownerName" valid:"stringlength(1|30)"`
	OwnerSurname string   `json:"ownerSurname" valid:"stringlength(1|30)"`
	Views        int      `json:"views" valid:"numeric"`
	Likes        int      `json:"likes" valid:"numeric"`
	Tariff       int      `json:"tariff" valid:"numeric"`
}

type ProductListData struct {
	ID         uint64   `json:"id" valid:"numeric"`
	Name       string   `json:"name" valid:"stringlength(1|100)"`
	Date       string   `json:"date" valid:"-"`
	Amount     int      `json:"amount" valid:"numeric"`
	LinkImages []string `json:"linkImages" valid:"stringArray"`
	UserLiked  bool     `json:"userLiked" valid:"type(bool)"`
	Tariff     int      `json:"tariff" valid:"numeric"`
}

type Page struct {
	From  uint64 `valid:"numeric"`
	Count uint64 `valid:"numeric"`
}

type Category struct {
	Title string `json:"title"`
}
