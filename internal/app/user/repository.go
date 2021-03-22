package user

import "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"

type UserRepository interface {
	Insert(user *models.UserData) error
	SelectByTelephone(telephone string) (*models.UserData, error)
	SelectByID(userID uint64) (*models.UserData, error)
	Update(user *models.UserData) error
}
