package usecase

import (
	"database/sql"
	errors2 "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
	"golang.org/x/crypto/bcrypt"
	"mime/multipart"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/upload"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/user"
)

type UserUsecase struct {
	userRepo   user.UserRepository
	uploadRepo upload.UploadRepository
}

func NewUserUsecase(repo user.UserRepository, uploadRepo upload.UploadRepository) user.UserUsecase {
	return &UserUsecase{
		userRepo:   repo,
		uploadRepo: uploadRepo,
	}
}

func (uu *UserUsecase) Create(user *models.UserData) *errors2.Error {
	if _, err := uu.GetByTelephone(user.Telephone); err == nil {
		return errors2.Cause(errors2.TelephoneAlreadyExists)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors2.UnexpectedInternal(err)
	}
	user.Password = string(hashedPassword)

	err = uu.userRepo.Insert(user)
	if err != nil {
		return errors2.UnexpectedInternal(err)
	}

	return nil
}

func (uu *UserUsecase) GetByTelephone(telephone string) (*models.UserData, *errors2.Error) {
	user, err := uu.userRepo.SelectByTelephone(telephone)

	switch {
	case err == sql.ErrNoRows:
		return nil, errors2.Cause(errors2.UserNotExist)
	case err != nil:
		return nil, errors2.UnexpectedInternal(err)
	}

	return user, nil
}

func (uu *UserUsecase) GetByID(userID uint64) (*models.ProfileData, *errors2.Error) {
	user, err := uu.userRepo.SelectByID(userID)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors2.Cause(errors2.UserNotExist)
	case err != nil:
		return nil, errors2.UnexpectedInternal(err)
	}

	profile := &models.ProfileData{
		ID: 		userID,
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

func (uu *UserUsecase) GetSellerByID(userID uint64) (*models.SellerData, *errors2.Error) {
	user, err := uu.userRepo.SelectByID(userID)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors2.Cause(errors2.UserNotExist)
	case err != nil:
		return nil, errors2.UnexpectedInternal(err)
	}

	profile := &models.SellerData{
		ID:         user.ID,
		Name:       user.Name,
		Surname:    user.Surname,
		Telephone:  user.Telephone,
		LinkImages: user.LinkImages,
	}

	return profile, nil
}

func (uu *UserUsecase) UpdateProfile(userID uint64, newUserData *models.UserData) (*models.UserData, *errors2.Error) {
	oldUser, err := uu.userRepo.SelectByID(userID)
	if err != nil {
		return nil, errors2.Cause(errors2.UserNotExist)
	}

	newUserData.ID = userID
	newUserData.Password = oldUser.Password
	newUserData.LinkImages = oldUser.LinkImages
	err = uu.userRepo.Update(newUserData)
	if err != nil {
		return nil, errors2.UnexpectedInternal(err)
	}

	return newUserData, nil
}

func (uu *UserUsecase) UpdateAvatar(userID uint64, fileHeader *multipart.FileHeader) (*models.UserData, *errors2.Error) {
	user, err := uu.userRepo.SelectByID(userID)
	if err != nil {
		return nil, errors2.Cause(errors2.UserNotExist)
	}

	imgUrl, err := uu.uploadRepo.InsertPhoto(fileHeader, "static/avatar/")
	if err != nil {
		return nil, errors2.UnexpectedInternal(err)
	}

	oldAvatar := user.LinkImages
	user.LinkImages = imgUrl
	err = uu.userRepo.Update(user)
	if err != nil {
		return nil, errors2.UnexpectedInternal(err)
	}

	err = uu.uploadRepo.RemovePhoto(oldAvatar)
	if err != nil {
		return nil, errors2.UnexpectedInternal(err)
	}

	return user, nil
}

func (uu *UserUsecase) CheckPassword(user *models.UserData, password string) *errors2.Error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return errors2.Cause(errors2.WrongPassword)
	}

	return nil
}

func (uu *UserUsecase) UpdatePassword(userID uint64, password string) (*models.UserData, *errors2.Error) {
	user, err := uu.userRepo.SelectByID(userID)
	if err != nil {
		return nil, errors2.Cause(errors2.UserNotExist)
	}

	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors2.UnexpectedInternal(err)
	}
	user.Password = string(newHashedPassword)

	err = uu.userRepo.Update(user)
	if err != nil {
		return nil, errors2.UnexpectedInternal(err)
	}

	return user, nil
}

func (uu *UserUsecase) UpdateLocation(userID uint64, data *models.LocationRequest) (*models.UserData, *errors2.Error) {
	user, err := uu.userRepo.SelectByID(userID)
	if err != nil {
		return nil, errors2.Cause(errors2.UserNotExist)
	}

	user.Latitude = data.Latitude
	user.Longitude = data.Longitude
	user.Radius = data.Radius
	user.Address = data.Address

	err = uu.userRepo.Update(user)
	if err != nil {
		return nil, errors2.UnexpectedInternal(err)
	}

	return user, nil
}
