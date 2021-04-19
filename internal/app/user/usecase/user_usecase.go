package usecase

import (
	"database/sql"
	"os"

	"golang.org/x/crypto/bcrypt"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/user"
)

type UserUsecase struct {
	userRepo user.UserRepository
}

func NewUserUsecase(repo user.UserRepository) user.UserUsecase {
	return &UserUsecase{
		userRepo: repo,
	}
}

func (uu *UserUsecase) Create(signUpData *models.SignUpRequest) (*models.UserData, *errors.Error) {
	user := &models.UserData{
		Name:       signUpData.Name,
		Surname:    signUpData.Surname,
		Sex:        signUpData.Sex,
		Email:      signUpData.Email,
		Telephone:  signUpData.Telephone,
		Password:   signUpData.Password1,
		DateBirth:  signUpData.DateBirth,
		LinkImages: signUpData.LinkImages,
	}

	if _, err := uu.GetByTelephone(user.Telephone); err == nil {
		return nil, errors.Cause(errors.TelephoneAlreadyExists)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.UnexpectedInternal(err)
	}
	user.Password = string(hashedPassword)

	err = uu.userRepo.Insert(user)
	if err != nil {
		return nil, errors.UnexpectedInternal(err)
	}

	return user, nil
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

func (uu *UserUsecase) GetSellerByID(userID uint64) (*models.SellerData, *errors.Error) {
	user, err := uu.userRepo.SelectByID(userID)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.Cause(errors.UserNotExist)
	case err != nil:
		return nil, errors.UnexpectedInternal(err)
	}

	seller := &models.SellerData{
		ID:         user.ID,
		Name:       user.Name,
		Surname:    user.Surname,
		Telephone:  user.Telephone,
		LinkImages: user.LinkImages,
	}

	return seller, nil
}

func (uu *UserUsecase) UpdateProfile(userID uint64, profileData *models.ProfileChangeRequest) (*models.UserData, *errors.Error) {
	user := &models.UserData{
		Name:      profileData.Name,
		Surname:   profileData.Surname,
		Sex:       profileData.Sex,
		Email:     profileData.Email,
		DateBirth: profileData.DateBirth,
	}

	oldUser, err := uu.userRepo.SelectByID(userID)
	if err != nil {
		return nil, errors.Cause(errors.UserNotExist)
	}

	user.ID = userID
	user.Telephone = oldUser.Telephone
	user.Password = oldUser.Password
	user.LinkImages = oldUser.LinkImages

	err = uu.userRepo.Update(user)
	if err != nil {
		return nil, errors.UnexpectedInternal(err)
	}

	return user, nil
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

func (uu *UserUsecase) UpdateLocation(userID uint64, locationData *models.LocationChangeRequest) (*models.UserData, *errors.Error) {
	user, err := uu.userRepo.SelectByID(userID)
	if err != nil {
		return nil, errors.Cause(errors.UserNotExist)
	}

	user.Latitude = locationData.Latitude
	user.Longitude = locationData.Longitude
	user.Radius = locationData.Radius
	user.Address = locationData.Address

	err = uu.userRepo.Update(user)
	if err != nil {
		return nil, errors.UnexpectedInternal(err)
	}

	return user, nil
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

func (uu *UserUsecase) CheckPassword(user *models.UserData, password string) *errors.Error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return errors.Cause(errors.WrongPassword)
	}

	return nil
}
