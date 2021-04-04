package models

type UserData struct {
	ID         uint64 `json:"id" valid:"numeric"`
	Name       string `json:"name" valid:"stringlength(5|30)"`
	Surname    string `json:"surname" valid:"stringlength(5|30)"`
	Sex        string `json:"sex" valid:"in(male|female)"`
	Email      string `json:"email" valid:"email"`
	Telephone  string `json:"telephone" valid:"stringlength(|)"`
	Password   string `json:"password,omitempty" valid:"stringlength(|))"`
	DateBirth  string `json:"dateBirth" valid:"-"`
	LinkImages string `json:"linkImages" valid:"type(string)"`
}

type SignUpRequest struct {
	Name       string `json:"name" valid:"stringlength(5|30)"`
	Surname    string `json:"surname" valid:"stringlength(5|30)"`
	Sex        string `json:"sex" valid:"in(male|female)"`
	Email      string `json:"email" valid:"email"`
	Telephone  string `json:"telephone" valid:"stringlength(10|13)"`
	Password1  string `json:"password1" valid:"stringlength(6|30)"`
	Password2  string `json:"password2" valid:"stringlength(6|30)"`
	DateBirth  string `json:"dateBirth" valid:"-"`
	LinkImages string `json:"linkImages" valid:"type(string)"`
}

type PasswordChangeRequest struct {
	OldPassword string `json:"oldPassword" valid:"stringlength(6|30)"`
	NewPassword string `json:"newPassword" valid:"stringlength(6|30)"`
}
