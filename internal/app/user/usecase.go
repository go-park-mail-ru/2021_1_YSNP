package user

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	errors2 "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
	"mime/multipart"
)

//go:generate mockgen -destination=./mocks/mock_user_ucase.go -package=mock github.com/go-park-mail-ru/2021_1_YSNP/internal/app/user UserUsecase

type UserUsecase interface {
	Create(user *models.UserData) *errors2.Error
	GetByTelephone(telephone string) (*models.UserData, *errors2.Error)
	GetByID(userID uint64) (*models.ProfileData, *errors2.Error)
	GetSellerByID(userID uint64) (*models.SellerData, *errors2.Error)
	UpdateProfile(userID uint64, newUserData *models.UserData) (*models.UserData, *errors2.Error)
	UpdateAvatar(userID uint64, fileHeader *multipart.FileHeader) (*models.UserData, *errors2.Error)
	CheckPassword(user *models.UserData, password string) *errors2.Error
	UpdatePassword(userID uint64, password string) (*models.UserData, *errors2.Error)
	UpdateLocation(userID uint64, data *models.LocationRequest) (*models.UserData, *errors2.Error)
}
