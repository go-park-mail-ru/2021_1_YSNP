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

type ProfileData struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Sex        string `json:"sex"`
	Email      string `json:"email"`
	Telephone  string `json:"telephone"`
	DateBirth  string `json:"dateBirth"`
	LinkImages string `json:"linkImages"`
}

type SignUpRequest struct {
	Name       string `json:"name" valid:"stringlength(5|30)"`
	Surname    string `json:"surname" valid:"stringlength(5|30)"`
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
