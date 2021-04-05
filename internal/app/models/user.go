package models

type UserData struct {
	ID         uint64 `json:"id" valid:"numeric"`
	Name       string `json:"name" valid:"stringlength(5|30)"`
	Surname    string `json:"surname" valid:"stringlength(5|30)"`
	Sex        string `json:"sex" valid:"in(male|female)"`
	Email      string `json:"email" valid:"email"`
	Telephone  string `json:"telephone" valid:"phoneNumber"`
	Password   string `json:"password,omitempty" valid:"password"`
	DateBirth  string `json:"dateBirth" valid:"-"`
	LinkImages string `json:"linkImages" valid:"type(string)"`
}

type SignUpRequest struct {
	Name       string `json:"name" valid:"stringlength(5|30)"`
	Surname    string `json:"surname" valid:"stringlength(5|30)"`
	Sex        string `json:"sex" valid:"in(male|female)"`
	Email      string `json:"email" valid:"email"`
	Telephone  string `json:"telephone" valid:"phoneNumber"`
	Password1  string `json:"password1" valid:"password"`
	Password2  string `json:"password2" valid:"password"`
	DateBirth  string `json:"dateBirth" valid:"-"`
	LinkImages string `json:"linkImages" valid:"type(string)"`
}

type PasswordChangeRequest struct {
	OldPassword string `json:"oldPassword" valid:"stringlength(6|30)"`
	NewPassword string `json:"newPassword" valid:"stringlength(6|30)"`
	//NewPassword2 string `json:"newPassword2" valid:"stringlength(6|30)"`
}
