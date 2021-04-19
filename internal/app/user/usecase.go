package user

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
)

//go:generate mockgen -destination=./mocks/mock_user_ucase.go -package=mock github.com/go-park-mail-ru/2021_1_YSNP/internal/app/user UserUsecase

type UserUsecase interface {
	Create(signUp *models.SignUpRequest) (*models.UserData, *errors.Error)
	UpdateAvatar(userID uint64, newAvatar string) (*models.UserData, *errors.Error)

	GetByID(userID uint64) (*models.ProfileData, *errors.Error)
	GetSellerByID(userID uint64) (*models.SellerData, *errors.Error)

	UpdateProfile(userID uint64, changeData *models.ProfileChangeRequest) (*models.UserData, *errors.Error)
	UpdatePassword(userID uint64, password string) (*models.UserData, *errors.Error)
	UpdateLocation(userID uint64, data *models.LocationChangeRequest) (*models.UserData, *errors.Error)

	GetByTelephone(telephone string) (*models.UserData, *errors.Error)
	CheckPassword(user *models.UserData, password string) *errors.Error
}
