package user

import "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"

//go:generate mockgen -destination=./mocks/mock_user_repo.go -package=mock github.com/go-park-mail-ru/2021_1_YSNP/internal/app/user UserRepository

type UserRepository interface {
	Insert(user *models.UserData) error
	SelectByTelephone(telephone string) (*models.UserData, error)
	SelectByID(userID uint64) (*models.UserData, error)
	Update(user *models.UserData) error

	InsertOAuth(userOAuth *models.UserOAuthRequest) error
	SelectByOAuthID(userOAuthID float64) uint64
}
