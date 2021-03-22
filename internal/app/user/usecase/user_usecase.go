package usecase

import(
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/user"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
)

type UserUsecase struct {
	userRepo user.UserRepository
}

func NewUserUsecase(repo user.UserRepository) user.UserUsecase {
	return &UserUsecase{
			userRepo: repo,
	}
}

func (uu *UserUsecase) Create(user *models.UserData) *errors.Error {
	panic("implement me")
}

func (uu *UserUsecase) GetByTelephone(telephone string) (*models.UserData, *errors.Error) {
	panic("implement me")
}

func (uu *UserUsecase) GetByID(userID uint64) (*models.UserData, *errors.Error) {
	panic("implement me")
}

func (uu *UserUsecase) UpdateProfile(userID uint64, newUserData *models.UserData) (*models.UserData, *errors.Error) {
	panic("implement me")
}

func (uu *UserUsecase) UpdateAvatar(userID uint64, newAvatar string) (*models.UserData, *errors.Error) {
	panic("implement me")
}

func (uu *UserUsecase) CheckPassword(user *models.UserData, password string) *errors.Error {
	panic("implement me")
}



