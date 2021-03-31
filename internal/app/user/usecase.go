package user

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
)

type UserUsecase interface {
	Create(user *models.UserData) *errors.Error
	GetByTelephone(telephone string) (*models.UserData, *errors.Error)
	GetByID(userID uint64) (*models.UserData, *errors.Error)
	UpdateProfile(userID uint64, newUserData *models.UserData) (*models.UserData, *errors.Error)
	UpdateAvatar(userID uint64, newAvatar string) (*models.UserData, *errors.Error)
	CheckPassword(user *models.UserData, password string) *errors.Error
	UpdatePassword(userID uint64, password string) (*models.UserData, *errors.Error)
}
