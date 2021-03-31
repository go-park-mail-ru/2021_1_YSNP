package models

type UserData struct {
	ID         uint64 `json:"id"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Sex        string `json:"sex"`
	Email      string `json:"email"`
	Telephone  string `json:"telephone"`
	Password   string `json:"password,omitempty"`
	DateBirth  string `json:"dateBirth"`
	LinkImages string `json:"linkImages"`
}

type SignUpRequest struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Sex        string `json:"sex"`
	Email      string `json:"email"`
	Telephone  string `json:"telephone"`
	Password1  string `json:"password1"`
	Password2  string `json:"password2"`
	DateBirth  string `json:"dateBirth"`
	LinkImages string `json:"linkImages"`
}

type PasswordChangeRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}
