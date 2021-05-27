package models

//easyjson:json
type UserData struct {
	ID         uint64  `json:"id" valid:"numeric"`
	Name       string  `json:"name" valid:"stringlength(1|30)"`
	Surname    string  `json:"surname" valid:"stringlength(1|30)"`
	Sex        string  `json:"sex" valid:"in(male|female|notstated)"`
	Email      string  `json:"email" valid:"email"`
	Telephone  string  `json:"telephone" valid:"phoneNumber"`
	Password   string  `json:"password,omitempty" valid:"password"`
	DateBirth  string  `json:"dateBirth" valid:"-"`
	Latitude   float64 `json:"latitude" valid:"latitude"`
	Longitude  float64 `json:"longitude" valid:"longitude"`
	Radius     uint64  `json:"radius" valid:"numeric"`
	Address    string  `json:"address" valid:"type(string)"`
	LinkImages string  `json:"linkImages" valid:"type(string)"`
	Rating 	   float64 `json:"rating"`
}

//easyjson:json
type ProfileData struct {
	ID         uint64  `json:"id" valid:"numeric"`
	Name       string  `json:"name" valid:"stringlength(1|30)"`
	Surname    string  `json:"surname" valid:"stringlength(1|30)"`
	Sex        string  `json:"sex" valid:"in(male|female|notstated)"`
	Email      string  `json:"email" valid:"email"`
	Telephone  string  `json:"telephone" valid:"phoneNumber"`
	DateBirth  string  `json:"dateBirth" valid:"-"`
	Latitude   float64 `json:"latitude" valid:"latitude"`
	Longitude  float64 `json:"longitude" valid:"longitude"`
	Radius     uint64  `json:"radius" valid:"numeric"`
	Address    string  `json:"address" valid:"type(string)"`
	LinkImages string  `json:"linkImages" valid:"type(string)"`
	Rating     float64 `json:"rating"`
}

//easyjson:json
type SellerData struct {
	ID      uint64 `json:"id" valid:"numeric"`
	Name    string `json:"name" valid:"stringlength(1|30)"`
	Surname string `json:"surname" valid:"stringlength(1|30)"`
	//Sex        string `json:"sex" valid:"in(male|female)"`
	//Email      string `json:"email" valid:"email"`
	Telephone string `json:"telephone" valid:"phoneNumber"`
	//DateBirth  string `json:"dateBirth" valid:"-"`
	LinkImages string `json:"linkImages" valid:"type(string)"`
	Rating  float64   `json:"rating"`
}

//easyjson:json
type LocationRequest struct {
	Latitude  float64 `json:"latitude" valid:"latitude"`
	Longitude float64 `json:"longitude" valid:"longitude"`
	Radius    uint64  `json:"radius" valid:"numeric"`
	Address   string  `json:"address" valid:"type(string)"`
}

//easyjson:json
type SignUpRequest struct {
	Name       string `json:"name" valid:"stringlength(1|30)"`
	Surname    string `json:"surname" valid:"stringlength(1|30)"`
	Sex        string `json:"sex" valid:"in(male|female|notstated)"`
	Email      string `json:"email" valid:"email"`
	Telephone  string `json:"telephone" valid:"phoneNumber"`
	Password1  string `json:"password1" valid:"password, password1"`
	Password2  string `json:"password2" valid:"password, password2"`
	DateBirth  string `json:"dateBirth" valid:"-"`
	LinkImages string `json:"linkImages" valid:"type(string)"`
}

//easyjson:json
type PasswordChangeRequest struct {
	OldPassword  string `json:"oldPassword" valid:"password"`
	NewPassword1 string `json:"newPassword1" valid:"password, password1"`
	NewPassword2 string `json:"newPassword2" valid:"password, password2"`
}

//easyjson:json
type Response struct {
	Response []struct {
		LastName  string `json:"last_name"`
		FirstName string `json:"first_name"`
		Photo     string `json:"photo_max"`
	} `json:"response"`
}

//easyjson:json
type UserOAuthRequest struct {
	ID            uint64  `json:"id"`
	LastName      string  `json:"last_name"`
	FirstName     string  `json:"first_name"`
	Photo         string  `json:"photo_max"`
	UserOAuthID   float64 `json:"user_oauth_id"`
	UserOAuthType string  `json:"user_oauth_type"`
}

type Achievement struct {
	Titie       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
	LinkPic     string `json:"link_pic"`
	Achieved   	bool `json:"achieved"`
}
