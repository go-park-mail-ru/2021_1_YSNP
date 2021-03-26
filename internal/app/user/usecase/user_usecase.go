package usecase

import(
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/user"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"golang.org/x/crypto/bcrypt"
	"os"
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		//TODO: создать ошибку
	}

	user.Password = string(hashedPassword)

	err = uu.userRepo.Insert(user)
	if err != nil {
		//TODO: создать ошибку
	}

	return nil
}

func (uu *UserUsecase) GetByTelephone(telephone string) (*models.UserData, *errors.Error) {
	user, err := uu.userRepo.SelectByTelephone(telephone)
	if err != nil {
		//TODO: создать ошибку
	}

	return user, nil
}

func (uu *UserUsecase) GetByID(userID uint64) (*models.UserData, *errors.Error) {
	user, err := uu.userRepo.SelectByID(userID)
	if err != nil {
		//TODO: создать ошибку
	}

	return user, nil
}

func (uu *UserUsecase) UpdateProfile(userID uint64, newUserData *models.UserData) (*models.UserData, *errors.Error) {
	//user, err := uu.userRepo.SelectByID(userID)
	//if err != nil {
	//	//TODO: создать ошибку
	//}

	newUserData.ID = userID

	err := uu.userRepo.Update(newUserData)
	if err != nil {
		//TODO: создать ошибку
	}

	return newUserData, nil
}

func (uu *UserUsecase) UpdateAvatar(userID uint64, newAvatar string) (*models.UserData, *errors.Error) {
	user, err := uu.userRepo.SelectByID(userID)
	if err != nil {
		//TODO: создать ошибку
	}

	oldAvatar := user.LinkImages
	user.LinkImages = newAvatar
	err = uu.userRepo.Update(user)
	if err != nil {
		//TODO: создать ошибку
	}

	if len(oldAvatar) != 0 {
		err := os.Remove(oldAvatar)
		if err != nil {
			//TODO: создать ошибку
		}
	}

	return user, nil
}

func (uu *UserUsecase) CheckPassword(user *models.UserData, password string) *errors.Error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		//TODO: создать ошибку
	}
	return nil
}

func (uu *UserUsecase) UpdatePassword(userID uint64, password string) (*models.UserData, *errors.Error) {
	user, err := uu.userRepo.SelectByID(userID)
	if err != nil {
		//TODO: создать ошибку
	}

	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		//TODO: создать ошибку
	}

	user.Password = string(newHashedPassword)

	err = uu.userRepo.Update(user)
	if err != nil {
		//TODO: создать ошибку
	}

	return user, nil
}

