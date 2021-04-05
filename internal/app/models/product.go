package models

type ProductData struct {
	ID           uint64   `json:"id" valid:"numeric"`
	Name         string   `json:"name" valid:"stringlength(1|100)"`
	Date         string   `json:"date" valid:"-"`
	Amount       int      `json:"amount" valid:"numeric"`
	LinkImages   []string `json:"linkImages" valid:"stringArray"`
	Description  string   `json:"description" valid:"stringlength(10|4000)"`
	Category     string   `json:"category" valid:"type(string)"`
	Address		 string	  `json:"address"`
	Longitude    string   `json:"longitude"`
	Latitude     string   `json:"latitude"`
	OwnerID      uint64   `json:"ownerId" valid:"numeric"`
	OwnerName    string   `json:"ownerName" valid:"stringlength(5|30)"`
	OwnerSurname string   `json:"ownerSurname" valid:"stringlength(5|30)"`
	Views        int      `json:"views" valid:"numeric"`
	Likes        int      `json:"likes" valid:"numeric"`
}

type ProductListData struct {
	ID         uint64   `json:"id" valid:"numeric"`
	Name       string   `json:"name" valid:"stringlength(1|100)"`
	Date       string   `json:"date" valid:"-"`
	Amount     int      `json:"amount" valid:"numeric"`
	LinkImages []string `json:"linkImages" valid:"stringArray"`
}

type Content struct {
	From  uint64 `json:"from" valid:"numeric"`
	Count uint64 `json:"count" valid:"numeric"`
}

type OrderType struct {
	Main bool `json:"main" valid:"type(bool)"`
}

type Page struct {
	Content Content   `json:"content" valid:"-"`
	Order   OrderType `json:"order" valid:"-"`
}
