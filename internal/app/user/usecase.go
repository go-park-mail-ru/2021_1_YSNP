package user

import (
	"mime/multipart"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
)

//go:generate mockgen -destination=./mocks/mock_user_ucase.go -package=mock github.com/go-park-mail-ru/2021_1_YSNP/internal/app/user UserUsecase

type UserUsecase interface {
	Create(user *models.UserData) *errors.Error
	GetByTelephone(telephone string) (*models.UserData, *errors.Error)
	GetByID(userID uint64) (*models.ProfileData, *errors.Error)
	GetSellerByID(userID uint64) (*models.SellerData, *errors.Error)
	UpdateProfile(userID uint64, newUserData *models.UserData) (*models.UserData, *errors.Error)
	UpdateAvatar(userID uint64, fileHeader *multipart.FileHeader) (*models.UserData, *errors.Error)
	CheckPassword(user *models.UserData, password string) *errors.Error
	UpdatePassword(userID uint64, password string) (*models.UserData, *errors.Error)
	UpdateLocation(userID uint64, data *models.LocationRequest) (*models.UserData, *errors.Error)
	Delete(userID uint64) *errors.Error

	CreateOrLogin(userOAuth *models.UserOAuthRequest) *errors.Error
}
