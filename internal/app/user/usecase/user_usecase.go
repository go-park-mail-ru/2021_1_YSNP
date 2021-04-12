package usecase

import (
	"database/sql"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/user"
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
	if _, err := uu.GetByTelephone(user.Telephone); err == nil {
		return errors.Cause(errors.TelephoneAlreadyExists)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.UnexpectedInternal(err)
	}
	user.Password = string(hashedPassword)

	err = uu.userRepo.Insert(user)
	if err != nil {
		return errors.UnexpectedInternal(err)
	}

	return nil
}

func (uu *UserUsecase) GetByTelephone(telephone string) (*models.UserData, *errors.Error) {
	user, err := uu.userRepo.SelectByTelephone(telephone)

	switch {
	case err == sql.ErrNoRows:
		return nil, errors.Cause(errors.UserNotExist)
	case err != nil:
		return nil, errors.UnexpectedInternal(err)
	}

	return user, nil
}

func (uu *UserUsecase) GetByID(userID uint64) (*models.ProfileData, *errors.Error) {
	user, err := uu.userRepo.SelectByID(userID)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.Cause(errors.UserNotExist)
	case err != nil:
		return nil, errors.UnexpectedInternal(err)
	}

	profile := &models.ProfileData{
		Name:       user.Name,
		Surname:    user.Surname,
		Sex:        user.Sex,
		Email:      user.Email,
		Telephone:  user.Telephone,
		DateBirth:  user.DateBirth,
		Latitude:   user.Latitude,
		Longitude:  user.Longitude,
		Radius:     user.Radius,
		Address:    user.Address,
		LinkImages: user.LinkImages,
	}

	return profile, nil
}

func (uu *UserUsecase) UpdateProfile(userID uint64, newUserData *models.UserData) (*models.UserData, *errors.Error) {
	oldUser, err := uu.userRepo.SelectByID(userID)
	if err != nil {
		return nil, errors.Cause(errors.UserNotExist)
	}

	newUserData.ID = userID
	newUserData.Password = oldUser.Password

	err = uu.userRepo.Update(newUserData)
	if err != nil {
		return nil, errors.UnexpectedInternal(err)
	}

	return newUserData, nil
}

func (uu *UserUsecase) UpdateAvatar(userID uint64, newAvatar string) (*models.UserData, *errors.Error) {
	user, err := uu.userRepo.SelectByID(userID)
	if err != nil {
		return nil, errors.Cause(errors.UserNotExist)
	}

	oldAvatar := user.LinkImages
	user.LinkImages = newAvatar
	err = uu.userRepo.Update(user)
	if err != nil {
		return nil, errors.UnexpectedInternal(err)
	}

	if oldAvatar != "" {
		err := os.Remove(oldAvatar)
		if err != nil {
			return nil, errors.UnexpectedInternal(err)
		}
	}

	return user, nil
}

func (uu *UserUsecase) CheckPassword(user *models.UserData, password string) *errors.Error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return errors.Cause(errors.WrongPassword)
	}

	return nil
}

func (uu *UserUsecase) UpdatePassword(userID uint64, password string) (*models.UserData, *errors.Error) {
	user, err := uu.userRepo.SelectByID(userID)
	if err != nil {
		return nil, errors.Cause(errors.UserNotExist)
	}

	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.UnexpectedInternal(err)
	}
	user.Password = string(newHashedPassword)

	err = uu.userRepo.Update(user)
	if err != nil {
		return nil, errors.UnexpectedInternal(err)
	}

	return user, nil
}

func (uu *UserUsecase) UpdatePosition(userID uint64, data *models.PositionData) (*models.UserData, *errors.Error) {
	user, err := uu.userRepo.SelectByID(userID)
	if err != nil {
		return nil, errors.Cause(errors.UserNotExist)
	}

	user.Latitude = data.Latitude
	user.Longitude = data.Longitude
	user.Radius = data.Radius
	user.Address = data.Address

	err = uu.userRepo.Update(user)
	if err != nil {
		return nil, errors.UnexpectedInternal(err)
	}

	return user, nil
}
