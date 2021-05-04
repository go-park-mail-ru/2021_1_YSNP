package usecase

import (
	"database/sql"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	errors "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
	uMock "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/upload/mocks"
	mock "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/user/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"testing"
)

var userTest = &models.UserData{
	ID:         0,
	Name:       "Максим",
	Surname:    "Торжков",
	Sex:        "male",
	Email:      "a@a.ru",
	Telephone:  "+79169230768",
	Password:   "Qwerty12",
	DateBirth:  "2021-03-08",
	LinkImages: "/static/avatar/test-avatar1.jpg",
}

func TestUserUsecase_Create_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockUserRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	userUcase := NewUserUsecase(userRepo, uploadRepo)

	userRepo.EXPECT().SelectByTelephone(gomock.Eq(userTest.Telephone)).Return(nil, sql.ErrNoRows)
	userRepo.EXPECT().Insert(gomock.Eq(userTest)).Return(nil)

	err := userUcase.Create(userTest)
	assert.Equal(t, err, (*errors.Error)(nil))

	//error
	userRepo.EXPECT().SelectByTelephone(gomock.Eq(userTest.Telephone)).Return(nil, sql.ErrNoRows)
	userRepo.EXPECT().Insert(gomock.Eq(userTest)).Return(sql.ErrConnDone)

	err = userUcase.Create(userTest)
	assert.Equal(t, err, errors.UnexpectedInternal(sql.ErrConnDone))
}

func TestUserUsecase_Create_TelephoneAlreadyExists(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockUserRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	userUcase := NewUserUsecase(userRepo, uploadRepo)

	userRepo.EXPECT().SelectByTelephone(gomock.Eq(userTest.Telephone)).Return(userTest, nil)

	err := userUcase.Create(userTest)
	assert.Equal(t, err, errors.Cause(errors.TelephoneAlreadyExists))
}

func TestUserUsecase_GetByID_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockUserRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	userUcase := NewUserUsecase(userRepo, uploadRepo)

	userTestProfile := &models.ProfileData{
		Name:       "Максим",
		Surname:    "Торжков",
		Sex:        "male",
		Email:      "a@a.ru",
		Telephone:  "+79169230768",
		DateBirth:  "2021-03-08",
		LinkImages: "",
	}

	userRepo.EXPECT().SelectByID(gomock.Eq(userTest.ID)).Return(userTest, nil)

	user, err := userUcase.GetByID(userTest.ID)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, user, userTestProfile)

	//error
	userRepo.EXPECT().SelectByID(gomock.Eq(userTest.ID)).Return(nil, sql.ErrConnDone)

	_, err = userUcase.GetByID(userTest.ID)
	assert.Equal(t, err, errors.UnexpectedInternal(sql.ErrConnDone))
}

func TestUserUsecase_GetByID_UserNotExist(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockUserRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	userUcase := NewUserUsecase(userRepo, uploadRepo)

	userRepo.EXPECT().SelectByID(gomock.Eq(userTest.ID)).Return(nil, sql.ErrNoRows)

	user, err := userUcase.GetByID(userTest.ID)
	assert.Equal(t, err, errors.Cause(errors.UserNotExist))
	assert.Equal(t, user, (*models.ProfileData)(nil))
}

func TestUserUsecase_GetByTelephone_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockUserRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	userUcase := NewUserUsecase(userRepo, uploadRepo)

	userRepo.EXPECT().SelectByTelephone(gomock.Eq(userTest.Telephone)).Return(userTest, nil)

	user, err := userUcase.GetByTelephone(userTest.Telephone)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, user, userTest)

	//error
	userRepo.EXPECT().SelectByTelephone(gomock.Eq(userTest.Telephone)).Return(nil, sql.ErrConnDone)

	_, err = userUcase.GetByTelephone(userTest.Telephone)
	assert.Equal(t, err, errors.UnexpectedInternal(sql.ErrConnDone))
}

func TestUserUsecase_GetByTelephone_UserNotExist(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockUserRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	userUcase := NewUserUsecase(userRepo, uploadRepo)

	userRepo.EXPECT().SelectByTelephone(gomock.Eq(userTest.Telephone)).Return(nil, sql.ErrNoRows)

	user, err := userUcase.GetByTelephone(userTest.Telephone)
	assert.Equal(t, err, errors.Cause(errors.UserNotExist))
	assert.Equal(t, user, (*models.UserData)(nil))
}

func TestUserUsecase_UpdateProfile_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockUserRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	userUcase := NewUserUsecase(userRepo, uploadRepo)

	userRepo.EXPECT().SelectByID(gomock.Eq(userTest.ID)).Return(userTest, nil)
	userRepo.EXPECT().Update(gomock.Eq(userTest)).Return(nil)

	_, err := userUcase.UpdateProfile(userTest.ID, userTest)
	assert.Equal(t, err, (*errors.Error)(nil))

	//error
	userRepo.EXPECT().SelectByID(gomock.Eq(userTest.ID)).Return(userTest, nil)
	userRepo.EXPECT().Update(gomock.Eq(userTest)).Return(sql.ErrConnDone)

	_, err = userUcase.UpdateProfile(userTest.ID, userTest)
	assert.Equal(t, err, errors.UnexpectedInternal(sql.ErrConnDone))
}

func TestUserUsecase_UpdateProfile_UserNotExist(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockUserRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	userUcase := NewUserUsecase(userRepo, uploadRepo)

	userRepo.EXPECT().SelectByID(gomock.Eq(userTest.ID)).Return(nil, sql.ErrNoRows)

	_, err := userUcase.UpdateProfile(userTest.ID, userTest)
	assert.Equal(t, err, errors.Cause(errors.UserNotExist))
}

func TestUserUsecase_UpdatePassword_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockUserRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	userUcase := NewUserUsecase(userRepo, uploadRepo)

	userRepo.EXPECT().SelectByID(gomock.Eq(userTest.ID)).Return(userTest, nil)
	userRepo.EXPECT().Update(gomock.Eq(userTest)).Return(nil)

	_, err := userUcase.UpdatePassword(userTest.ID, "random_pass")
	assert.Equal(t, err, (*errors.Error)(nil))

	//error
	userRepo.EXPECT().SelectByID(gomock.Eq(userTest.ID)).Return(userTest, nil)
	userRepo.EXPECT().Update(gomock.Eq(userTest)).Return(sql.ErrConnDone)

	_, err = userUcase.UpdatePassword(userTest.ID, "random_pass")
	assert.Equal(t, err, errors.UnexpectedInternal(sql.ErrConnDone))
}

func TestUserUsecase_UpdatePassword_UserNotExist(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockUserRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	userUcase := NewUserUsecase(userRepo, uploadRepo)

	userRepo.EXPECT().SelectByID(gomock.Eq(userTest.ID)).Return(nil, sql.ErrNoRows)

	_, err := userUcase.UpdatePassword(userTest.ID, "random_password")
	assert.Equal(t, err, errors.Cause(errors.UserNotExist))
}

func TestUserUsecase_UpdateAvatar(t *testing.T) {
	//t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockUserRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	userUcase := NewUserUsecase(userRepo, uploadRepo)

	userRepo.EXPECT().SelectByID(gomock.Eq(userTest.ID)).Return(userTest, nil)
	uploadRepo.EXPECT().InsertPhoto(&multipart.FileHeader{}, "static/avatar/").Return("", nil)
	userRepo.EXPECT().Update(gomock.Eq(userTest)).Return(nil)
	uploadRepo.EXPECT().RemovePhoto(gomock.Any()).Return(nil)

	userTest.LinkImages = ""

	_, err := userUcase.UpdateAvatar(userTest.ID, &multipart.FileHeader{})
	assert.Equal(t, err, (*errors.Error)(nil))

	//error
	userRepo.EXPECT().SelectByID(gomock.Eq(userTest.ID)).Return(userTest, nil)
	uploadRepo.EXPECT().InsertPhoto(&multipart.FileHeader{}, "static/avatar/").Return("", nil)
	userRepo.EXPECT().Update(gomock.Eq(userTest)).Return(sql.ErrConnDone)

	userTest.LinkImages = ""

	_, err = userUcase.UpdateAvatar(userTest.ID, &multipart.FileHeader{})
	assert.Equal(t, err, errors.UnexpectedInternal(sql.ErrConnDone))

	//another err
	userRepo.EXPECT().SelectByID(gomock.Eq(userTest.ID)).Return(userTest, nil)
	uploadRepo.EXPECT().InsertPhoto(&multipart.FileHeader{}, "static/avatar/").Return("", nil)
	userRepo.EXPECT().Update(gomock.Eq(userTest)).Return(nil)
	uploadRepo.EXPECT().RemovePhoto(gomock.Any()).Return(sql.ErrConnDone)

	userTest.LinkImages = ""

	_, err = userUcase.UpdateAvatar(userTest.ID, &multipart.FileHeader{})
	assert.Equal(t, err, errors.UnexpectedInternal(sql.ErrConnDone))
}

func TestUserUsecase_UpdateAvatar_UserNotExist(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockUserRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	userUcase := NewUserUsecase(userRepo, uploadRepo)

	userRepo.EXPECT().SelectByID(gomock.Eq(userTest.ID)).Return(nil, sql.ErrNoRows)

	_, err := userUcase.UpdateAvatar(userTest.ID,  &multipart.FileHeader{})
	assert.Equal(t, err, errors.Cause(errors.UserNotExist))
}

func TestUserUsecase_UpdateAvatar_Error(t *testing.T) {
	//t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockUserRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	userUcase := NewUserUsecase(userRepo, uploadRepo)

	userRepo.EXPECT().SelectByID(gomock.Eq(userTest.ID)).Return(userTest, nil)
	uploadRepo.EXPECT().InsertPhoto(gomock.Any(), "static/avatar/").Return("", sql.ErrConnDone)

	_, err := userUcase.UpdateAvatar(userTest.ID,  &multipart.FileHeader{})
	assert.Equal(t, err.ErrorCode, errors.InternalError)
}

func TestUserUsecase_CheckPassword_WrongPassword(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockUserRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	userUcase := NewUserUsecase(userRepo, uploadRepo)

	err := userUcase.CheckPassword(userTest, "password")
	assert.Equal(t, err, errors.Cause(errors.WrongPassword))
}

func TestUserUsecase_UpdatePosition_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockUserRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	userUcase := NewUserUsecase(userRepo, uploadRepo)

	position := &models.LocationRequest{
		Latitude:  1,
		Longitude: 1,
		Radius:    1,
		Address:   "address",
	}

	userLocalTest := &models.UserData{
		ID:         0,
		Name:       "Максим",
		Surname:    "Торжков",
		Sex:        "male",
		Email:      "a@a.ru",
		Telephone:  "+79169230768",
		Password:   "Qwerty12",
		DateBirth:  "2021-03-08",
		LinkImages: "",
	}

	userWithPosit := &models.UserData{
		ID:         0,
		Name:       "Максим",
		Surname:    "Торжков",
		Sex:        "male",
		Email:      "a@a.ru",
		Telephone:  "+79169230768",
		Password:   "Qwerty12",
		DateBirth:  "2021-03-08",
		LinkImages: "",
		Latitude:   1,
		Longitude:  1,
		Radius:     1,
		Address:    "address",
	}

	userRepo.EXPECT().SelectByID(gomock.Eq(userTest.ID)).Return(userLocalTest, nil)
	userRepo.EXPECT().Update(gomock.Eq(userWithPosit)).Return(nil)

	_, err := userUcase.UpdateLocation(userTest.ID, position)
	assert.Equal(t, err, (*errors.Error)(nil))

	//error
	userRepo.EXPECT().SelectByID(gomock.Eq(userTest.ID)).Return(userLocalTest, nil)
	userRepo.EXPECT().Update(gomock.Eq(userWithPosit)).Return(sql.ErrConnDone)

	_, err = userUcase.UpdateLocation(userTest.ID, position)
	assert.Equal(t, err, errors.UnexpectedInternal(sql.ErrConnDone))
}

func TestUserUsecase_UpdatePosition_UserNotExist(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockUserRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	userUcase := NewUserUsecase(userRepo, uploadRepo)

	position := &models.LocationRequest{
		Latitude:  1,
		Longitude: 1,
		Radius:    1,
		Address:   "address",
	}

	userRepo.EXPECT().SelectByID(gomock.Eq(userTest.ID)).Return(nil, sql.ErrNoRows)

	_, err := userUcase.UpdateLocation(userTest.ID, position)
	assert.Equal(t, err, errors.Cause(errors.UserNotExist))
}

func TestUserUsecase_GetSellerByID_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockUserRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	userUcase := NewUserUsecase(userRepo, uploadRepo)

	userTestProfile := &models.SellerData{
		ID:         userTest.ID,
		Name:       "Максим",
		Surname:    "Торжков",
		Telephone:  "+79169230768",
		LinkImages: userTest.LinkImages,
	}

	userRepo.EXPECT().SelectByID(gomock.Eq(userTest.ID)).Return(userTest, nil)

	user, err := userUcase.GetSellerByID(userTest.ID)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, user, userTestProfile)

	//error
	userRepo.EXPECT().SelectByID(gomock.Eq(userTest.ID)).Return(nil, sql.ErrConnDone)

	_, err = userUcase.GetSellerByID(userTest.ID)
	assert.Equal(t, err, errors.UnexpectedInternal(sql.ErrConnDone))
}

func TestUserUsecase_GetSellerByID_UserNotExist(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockUserRepository(ctrl)
	uploadRepo := uMock.NewMockUploadRepository(ctrl)
	userUcase := NewUserUsecase(userRepo, uploadRepo)

	userRepo.EXPECT().SelectByID(gomock.Eq(userTest.ID)).Return(nil, sql.ErrNoRows)

	user, err := userUcase.GetSellerByID(userTest.ID)
	assert.Equal(t, err, errors.Cause(errors.UserNotExist))
	assert.Equal(t, user, (*models.SellerData)(nil))
}
