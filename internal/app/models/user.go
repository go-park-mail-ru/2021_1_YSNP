package models

type UserData struct {
	ID         uint64  `json:"id" valid:"numeric"`
	Name       string  `json:"name" valid:"stringlength(1|30)"`
	Surname    string  `json:"surname" valid:"stringlength(1|30)"`
	Sex        string  `json:"sex" valid:"in(male|female)"`
	Email      string  `json:"email" valid:"email"`
	Telephone  string  `json:"telephone" valid:"phoneNumber"`
	Password   string  `json:"password,omitempty" valid:"password"`
	DateBirth  string  `json:"dateBirth" valid:"-"`
	Latitude   float64 `json:"latitude" valid:"latitude"`
	Longitude  float64 `json:"longitude" valid:"longitude"`
	Radius     uint64  `json:"radius" valid:"numeric"`
	Address    string  `json:"address" valid:"type(string)"`
	LinkImages string  `json:"linkImages" valid:"type(string)"`
}

type ProfileData struct {
	ID         uint64  `json:"id" valid:"numeric"`
	Name       string  `json:"name" valid:"stringlength(1|30)"`
	Surname    string  `json:"surname" valid:"stringlength(1|30)"`
	Sex        string  `json:"sex" valid:"in(male|female)"`
	Email      string  `json:"email" valid:"email"`
	Telephone  string  `json:"telephone" valid:"phoneNumber"`
	DateBirth  string  `json:"dateBirth" valid:"-"`
	Latitude   float64 `json:"latitude" valid:"latitude"`
	Longitude  float64 `json:"longitude" valid:"longitude"`
	Radius     uint64  `json:"radius" valid:"numeric"`
	Address    string  `json:"address" valid:"type(string)"`
	LinkImages string  `json:"linkImages" valid:"type(string)"`
}

type SellerData struct {
	ID      uint64 `json:"id" valid:"numeric"`
	Name    string `json:"name" valid:"stringlength(1|30)"`
	Surname string `json:"surname" valid:"stringlength(1|30)"`
	//Sex        string `json:"sex" valid:"in(male|female)"`
	//Email      string `json:"email" valid:"email"`
	Telephone string `json:"telephone" valid:"phoneNumber"`
	//DateBirth  string `json:"dateBirth" valid:"-"`
	LinkImages string `json:"linkImages" valid:"type(string)"`
}

type PositionData struct {
	Latitude  float64 `json:"latitude" valid:"latitude"`
	Longitude float64 `json:"longitude" valid:"longitude"`
	Radius    uint64  `json:"radius" valid:"numeric"`
	Address   string  `json:"address" valid:"type(string)"`
}

type SignUpRequest struct {
	Name       string `json:"name" valid:"stringlength(1|30)"`
	Surname    string `json:"surname" valid:"stringlength(1|30)"`
	Sex        string `json:"sex" valid:"in(male|female)"`
	Email      string `json:"email" valid:"email"`
	Telephone  string `json:"telephone" valid:"phoneNumber"`
	Password1  string `json:"password1" valid:"password, password1"`
	Password2  string `json:"password2" valid:"password, password2"`
	DateBirth  string `json:"dateBirth" valid:"-"`
	LinkImages string `json:"linkImages" valid:"type(string)"`
}

type PasswordChangeRequest struct {
	OldPassword  string `json:"oldPassword" valid:"password"`
	NewPassword1 string `json:"newPassword1" valid:"password, password1"`
	NewPassword2 string `json:"newPassword2" valid:"password, password2"`
}
