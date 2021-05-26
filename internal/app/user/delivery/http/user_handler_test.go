package delivery

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/microservices/auth/mocks"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/middleware"
	_ "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/validator"
	uMock "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/user/mocks"
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
	LinkImages: "/static/avatar/profile.webp",
}

//var byteData = bytes.NewReader([]byte(`
//			{
//			"name":"Максим",
//			"surname":"Торжков",
//			"sex":"male",
//			"email":"a@a.ru",
//			"telephone":"+7916923076",
//			"password1":"Qwerty12",
//			"password2":"Qwerty12",
//			"dateBirth":"2021-03-08"
//			}
//	`))

func TestUserHandler_SignUpHandler_OK(t *testing.T) {
	//t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUcase := uMock.NewMockUserUsecase(ctrl)
	sessUcase := mock.NewMockSessionUsecase(ctrl)

	var byteData = bytes.NewReader([]byte(`
			{
			"name":"Максим",
			"surname":"Торжков",
			"sex":"male",
			"email":"a@a.ru",
			"telephone":"+79169230768",
			"password1":"Qwerty12",
			"password2":"Qwerty12",
			"dateBirth":"2021-03-08"
			}
	`))

	r := httptest.NewRequest("POST", "/api/v1/signup", byteData)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	mw := middleware.NewMiddleware(sessUcase, userUcase, nil)
	router.Use(middleware.CorsControlMiddleware)
	router.Use(mw.AccessLogMiddleware)
	userHandler := NewUserHandler(userUcase, sessUcase)
	userHandler.Configure(router, mw)

	userUcase.EXPECT().Create(userTest).Return(nil)
	sessUcase.EXPECT().Create(gomock.Any()).Return(nil)
	userHandler.SignUpHandler(w, r.WithContext(ctx))

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserHandler_SignUpHandler_LoggerError(t *testing.T) {
	//t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUcase := uMock.NewMockUserUsecase(ctrl)
	sessUcase := mock.NewMockSessionUsecase(ctrl)

	var byteData = bytes.NewReader([]byte(`
			{
			"name":"Максим",
			"surname":"Торжков",
			"sex":"male",
			"email":"a@a.ru",
			"telephone":"+79169230768",
			"password1":"Qwerty12",
			"password2":"Qwerty12",
			"dateBirth":"2021-03-08"
			}
	`))

	r := httptest.NewRequest("POST", "/api/v1/signup", byteData)
	ctx := r.Context()
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	userHandler := NewUserHandler(userUcase, sessUcase)
	userHandler.Configure(router, nil)

	userUcase.EXPECT().Create(userTest).Return(nil)
	sessUcase.EXPECT().Create(gomock.Any()).Return(nil)
	userHandler.SignUpHandler(w, r.WithContext(ctx))

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserHandler_SignUpHandler_DecodeError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUcase := uMock.NewMockUserUsecase(ctrl)
	sessUcase := mock.NewMockSessionUsecase(ctrl)

	var byteData = bytes.NewReader([]byte(`
			{
			"name":12,
			"surname":"Торжков",
			"sex":"male",
			"email":"a@a.ru",
			"telephone":"+79169230768",
			"password1":"Qwerty12",
			"password2":"Qwerty12",
			"dateBirth":"2021-03-08"
			}
	`))

	r := httptest.NewRequest("POST", "/api/v1/signup", byteData)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	userHandler := NewUserHandler(userUcase, sessUcase)
	userHandler.Configure(router, nil)

	userHandler.SignUpHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_SignUpHandler_TelephoneAlreadyExists(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUcase := uMock.NewMockUserUsecase(ctrl)
	sessUcase := mock.NewMockSessionUsecase(ctrl)

	var byteData = bytes.NewReader([]byte(`
			{
			"name":"Максим",
			"surname":"Торжков",
			"sex":"male",
			"email":"a@a.ru",
			"telephone":"+79169230768",
			"password1":"Qwerty12",
			"password2":"Qwerty12",
			"dateBirth":"2021-03-08"
			}
	`))

	userTest := &models.UserData{
		ID:        0,
		Name:      "Максим",
		Surname:   "Торжков",
		Sex:       "male",
		Email:     "a@a.ru",
		Telephone: "+79169230768",
		Password:  "Qwerty12",
		DateBirth: "2021-03-08",
		LinkImages: "/static/avatar/profile.webp",
	}

	r := httptest.NewRequest("POST", "/api/v1/signup", byteData)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	userHandler := NewUserHandler(userUcase, sessUcase)
	userHandler.Configure(router, nil)

	userUcase.EXPECT().Create(userTest).Return(errors.Cause(errors.TelephoneAlreadyExists))
	userHandler.SignUpHandler(w, r.WithContext(ctx))

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_SignUpHandler_ValidationError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUcase := uMock.NewMockUserUsecase(ctrl)
	sessUcase := mock.NewMockSessionUsecase(ctrl)

	var byteData = bytes.NewReader([]byte(`
			{
			"name":"Максим",
			"surname":"Торжков",
			"sex":"male",
			"email":"a@a.ru",
			"telephone":"+7916923076",
			"password1":"Qwerty12",
			"password2":"Qwerty12",
			"dateBirth":"2021-03-08"
			}
	`)) // wrong telephone

	r := httptest.NewRequest("POST", "/api/v1/signup", byteData)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	userHandler := NewUserHandler(userUcase, sessUcase)
	userHandler.Configure(router, nil)

	userHandler.SignUpHandler(w, r.WithContext(ctx))

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_GetProfileHandler_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUcase := uMock.NewMockUserUsecase(ctrl)
	sessUcase := mock.NewMockSessionUsecase(ctrl)

	r := httptest.NewRequest("GET", "/api/v1/me", nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, userTest.ID)
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	userHandler := NewUserHandler(userUcase, sessUcase)
	userHandler.Configure(router, nil)

	userTestProfile := &models.ProfileData{
		Name:       "Максим",
		Surname:    "Торжков",
		Sex:        "male",
		Email:      "a@a.ru",
		Telephone:  "+79169230768",
		DateBirth:  "2021-03-08",
		LinkImages: "",
	}

	userUcase.EXPECT().GetByID(gomock.Eq(userTest.ID)).Return(userTestProfile, nil)

	userHandler.GetProfileHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserHandler_GetProfileHandler_LoggerError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUcase := uMock.NewMockUserUsecase(ctrl)
	sessUcase := mock.NewMockSessionUsecase(ctrl)

	r := httptest.NewRequest("GET", "/api/v1/me", nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextUserID, userTest.ID)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	userHandler := NewUserHandler(userUcase, sessUcase)
	userHandler.Configure(router, nil)

	userTestProfile := &models.ProfileData{
		Name:       "Максим",
		Surname:    "Торжков",
		Sex:        "male",
		Email:      "a@a.ru",
		Telephone:  "+79169230768",
		DateBirth:  "2021-03-08",
		LinkImages: "",
	}

	userUcase.EXPECT().GetByID(gomock.Eq(userTest.ID)).Return(userTestProfile, nil)

	userHandler.GetProfileHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserHandler_GetProfileHandler_UserUnauthorized(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUcase := uMock.NewMockUserUsecase(ctrl)
	sessUcase := mock.NewMockSessionUsecase(ctrl)

	r := httptest.NewRequest("GET", "/api/v1/me", nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	userHandler := NewUserHandler(userUcase, sessUcase)
	userHandler.Configure(router, nil)

	userHandler.GetProfileHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUserHandler_GetProfileHandler_UserNotExist(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUcase := uMock.NewMockUserUsecase(ctrl)
	sessUcase := mock.NewMockSessionUsecase(ctrl)

	r := httptest.NewRequest("GET", "/api/v1/me", nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, userTest.ID)
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	userHandler := NewUserHandler(userUcase, sessUcase)
	userHandler.Configure(router, nil)

	userUcase.EXPECT().GetByID(gomock.Eq(userTest.ID)).Return(nil, errors.Cause(errors.UserNotExist))

	userHandler.GetProfileHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_GetSellerHandler_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUcase := uMock.NewMockUserUsecase(ctrl)
	sessUcase := mock.NewMockSessionUsecase(ctrl)

	r := httptest.NewRequest("GET", "/api/v1/user/0", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "0"})
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, userTest.ID)
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	userHandler := NewUserHandler(userUcase, sessUcase)
	userHandler.Configure(router, nil)

	sellerTestProfile := &models.SellerData{
		Name:       "Максим",
		Surname:    "Торжков",
		Telephone:  "+79169230768",
		LinkImages: "",
	}

	userUcase.EXPECT().GetSellerByID(uint64(0)).Return(sellerTestProfile, nil)

	userHandler.GetSellerHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserHandler_GetSellerHandler_LoggerError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUcase := uMock.NewMockUserUsecase(ctrl)
	sessUcase := mock.NewMockSessionUsecase(ctrl)

	r := httptest.NewRequest("GET", "/api/v1/user/0", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "0"})
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextUserID, userTest.ID)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	userHandler := NewUserHandler(userUcase, sessUcase)
	userHandler.Configure(router, nil)

	sellerTestProfile := &models.SellerData{
		Name:       "Максим",
		Surname:    "Торжков",
		Telephone:  "+79169230768",
		LinkImages: "",
	}

	userUcase.EXPECT().GetSellerByID(uint64(0)).Return(sellerTestProfile, nil)

	userHandler.GetSellerHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserHandler_GetSellerHandler_UserNotExist(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUcase := uMock.NewMockUserUsecase(ctrl)
	sessUcase := mock.NewMockSessionUsecase(ctrl)

	r := httptest.NewRequest("GET", "/api/v1/user/0", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "0"})
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, userTest.ID)
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	userHandler := NewUserHandler(userUcase, sessUcase)
	userHandler.Configure(router, nil)

	userUcase.EXPECT().GetSellerByID(uint64(0)).Return(nil, errors.Cause(errors.UserNotExist))

	userHandler.GetSellerHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_ChangeProfileHandler_Success(t *testing.T) {
	//t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUcase := uMock.NewMockUserUsecase(ctrl)
	sessUcase := mock.NewMockSessionUsecase(ctrl)

	var byteData = bytes.NewReader([]byte(`
			{
			"name":"Максим",
			"surname":"Торжков",
			"sex":"male",
			"email":"a@a.ru",
			"telephone":"+79169230768",
			"password1":"Qwerty12",
			"password2":"Qwerty12",
			"dateBirth":"2021-03-08"
			}
	`))

	r := httptest.NewRequest("POST", "/api/v1/user", byteData)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, userTest.ID)
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	userHandler := NewUserHandler(userUcase, sessUcase)
	userHandler.Configure(router, nil)

	userTest.Password = ""
	userTest.LinkImages = ""

	userUcase.EXPECT().UpdateProfile(gomock.Eq(userTest.ID), gomock.Eq(userTest)).Return(userTest, nil)

	userHandler.ChangeProfileHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserHandler_ChangeProfileHandler_LoggerError(t *testing.T) {
	//t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUcase := uMock.NewMockUserUsecase(ctrl)
	sessUcase := mock.NewMockSessionUsecase(ctrl)

	var byteData = bytes.NewReader([]byte(`
			{
			"name":"Максим",
			"surname":"Торжков",
			"sex":"male",
			"email":"a@a.ru",
			"telephone":"+79169230768",
			"password1":"Qwerty12",
			"password2":"Qwerty12",
			"dateBirth":"2021-03-08"
			}
	`))

	r := httptest.NewRequest("POST", "/api/v1/user", byteData)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextUserID, userTest.ID)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	userHandler := NewUserHandler(userUcase, sessUcase)
	userHandler.Configure(router, nil)

	userTest.Password = ""

	userUcase.EXPECT().UpdateProfile(gomock.Eq(userTest.ID), gomock.Eq(userTest)).Return(userTest, nil)

	userHandler.ChangeProfileHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserHandler_ChangeProfileHandler_DecodeError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUcase := uMock.NewMockUserUsecase(ctrl)
	sessUcase := mock.NewMockSessionUsecase(ctrl)

	var byteData = bytes.NewReader([]byte(`
			{
			"name":12,
			"surname":"Торжков",
			"sex":"male",
			"email":"a@a.ru",
			"telephone":"+79169230768",
			"password1":"Qwerty12",
			"password2":"Qwerty12",
			"dateBirth":"2021-03-08"
			}
	`))

	r := httptest.NewRequest("POST", "/api/v1/user", byteData)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, userTest.ID)
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	userHandler := NewUserHandler(userUcase, sessUcase)
	userHandler.Configure(router, nil)

	userHandler.ChangeProfileHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_ChangeProfileHandler_NotAuth(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUcase := uMock.NewMockUserUsecase(ctrl)
	sessUcase := mock.NewMockSessionUsecase(ctrl)

	r := httptest.NewRequest("POST", "/api/v1/user", nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	userHandler := NewUserHandler(userUcase, sessUcase)
	userHandler.Configure(router, nil)

	userHandler.ChangeProfileHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUserHandler_ChangeProfileHandler_ValidateError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUcase := uMock.NewMockUserUsecase(ctrl)
	sessUcase := mock.NewMockSessionUsecase(ctrl)

	var byteData = bytes.NewReader([]byte(`
			{
			"name":"Максим",
			"surname":"Торжков",
			"sex":"male",
			"email":"a@a.ru",
			"telephone":"+7916923076",
			"password1":"Qwerty12",
			"password2":"Qwerty12",
			"dateBirth":"2021-03-08"
			}
	`))

	r := httptest.NewRequest("POST", "/api/v1/user", byteData)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, userTest.ID)
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	userHandler := NewUserHandler(userUcase, sessUcase)
	userHandler.Configure(router, nil)

	userHandler.ChangeProfileHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_ChangeProfileHandler_NoUser(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUcase := uMock.NewMockUserUsecase(ctrl)
	sessUcase := mock.NewMockSessionUsecase(ctrl)

	var byteData = bytes.NewReader([]byte(`
			{
			"name":"Максим",
			"surname":"Торжков",
			"sex":"male",
			"email":"a@a.ru",
			"telephone":"+79169230768",
			"password1":"Qwerty12",
			"password2":"Qwerty12",
			"dateBirth":"2021-03-08"
			}
	`))

	r := httptest.NewRequest("POST", "/api/v1/user", byteData)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, userTest.ID)
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	userHandler := NewUserHandler(userUcase, sessUcase)
	userHandler.Configure(router, nil)

	userUcase.EXPECT().UpdateProfile(gomock.Any(), gomock.Any()).Return(nil, errors.Cause(errors.UserNotExist))

	userHandler.ChangeProfileHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_ChangeProfilePasswordHandler_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUcase := uMock.NewMockUserUsecase(ctrl)
	sessUcase := mock.NewMockSessionUsecase(ctrl)

	var byteData = bytes.NewReader([]byte(`
			{
			"oldPassword":"Qwerty12",
			"newPassword1":"Qwerty123",
			"newPassword2":"Qwerty123"
			}
	`))

	r := httptest.NewRequest("POST", "/api/v1/user/password", byteData)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, userTest.ID)
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	userHandler := NewUserHandler(userUcase, sessUcase)
	userHandler.Configure(router, nil)

	userUcase.EXPECT().UpdatePassword(gomock.Eq(userTest.ID), gomock.Eq("Qwerty123")).Return(userTest, nil)

	userHandler.ChangeProfilePasswordHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserHandler_ChangeProfilePasswordHandler_LoggerError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUcase := uMock.NewMockUserUsecase(ctrl)
	sessUcase := mock.NewMockSessionUsecase(ctrl)

	var byteData = bytes.NewReader([]byte(`
			{
			"oldPassword":"Qwerty12",
			"newPassword1":"Qwerty123",
			"newPassword2":"Qwerty123"
			}
	`))

	r := httptest.NewRequest("POST", "/api/v1/user/password", byteData)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextUserID, userTest.ID)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	userHandler := NewUserHandler(userUcase, sessUcase)
	userHandler.Configure(router, nil)

	userUcase.EXPECT().UpdatePassword(gomock.Eq(userTest.ID), gomock.Eq("Qwerty123")).Return(userTest, nil)

	userHandler.ChangeProfilePasswordHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserHandler_ChangeProfilePasswordHandler_DecodeError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUcase := uMock.NewMockUserUsecase(ctrl)
	sessUcase := mock.NewMockSessionUsecase(ctrl)

	var byteData = bytes.NewReader([]byte(`
			{
			"oldPassword":12,
			"newPassword1":"Qwerty123",
			"newPassword2":"Qwerty123"
			}
	`))

	r := httptest.NewRequest("POST", "/api/v1/user/password", byteData)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, userTest.ID)
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	userHandler := NewUserHandler(userUcase, sessUcase)
	userHandler.Configure(router, nil)

	userHandler.ChangeProfilePasswordHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_ChangeProfilePasswordHandler_NotAuth(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUcase := uMock.NewMockUserUsecase(ctrl)
	sessUcase := mock.NewMockSessionUsecase(ctrl)

	r := httptest.NewRequest("POST", "/api/v1/user/password", nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	userHandler := NewUserHandler(userUcase, sessUcase)
	userHandler.Configure(router, nil)

	userHandler.ChangeProfilePasswordHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUserHandler_ChangeProfilePasswordHandler_ValidateError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUcase := uMock.NewMockUserUsecase(ctrl)
	sessUcase := mock.NewMockSessionUsecase(ctrl)

	var byteData = bytes.NewReader([]byte(`
			{
			"oldPassword":"Qwerty12",
			"newPassword1":"Qwerty123",
			"newPassword2":"Qwerty124"
			}
	`))

	r := httptest.NewRequest("POST", "/api/v1/user/password", byteData)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, userTest.ID)
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	userHandler := NewUserHandler(userUcase, sessUcase)
	userHandler.Configure(router, nil)

	userHandler.ChangeProfilePasswordHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_ChangeProfilePasswordHandler_NoUser(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUcase := uMock.NewMockUserUsecase(ctrl)
	sessUcase := mock.NewMockSessionUsecase(ctrl)

	var byteData = bytes.NewReader([]byte(`
			{
			"oldPassword":"Qwerty12",
			"newPassword1":"Qwerty123",
			"newPassword2":"Qwerty123"
			}
	`))

	r := httptest.NewRequest("POST", "/api/v1/user/password", byteData)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, userTest.ID)
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	userHandler := NewUserHandler(userUcase, sessUcase)
	userHandler.Configure(router, nil)

	userUcase.EXPECT().UpdatePassword(gomock.Any(), gomock.Any()).Return(nil, errors.Cause(errors.UserNotExist))

	userHandler.ChangeProfilePasswordHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_ChangeUSerPositionHandler_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUcase := uMock.NewMockUserUsecase(ctrl)
	sessUcase := mock.NewMockSessionUsecase(ctrl)

	var byteData = bytes.NewReader([]byte(`
			{
			"latitude": 1,
			"longitude": 1,
			"radius": 1,
			"address":"Qwerty123"
			}
	`))

	position := &models.LocationRequest{
		Latitude:  1,
		Longitude: 1,
		Radius:    1,
		Address:   "Qwerty123",
	}

	r := httptest.NewRequest("POST", "/api/v1/user/position", byteData)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, userTest.ID)
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	userHandler := NewUserHandler(userUcase, sessUcase)
	userHandler.Configure(router, nil)

	userUcase.EXPECT().UpdateLocation(gomock.Eq(userTest.ID), gomock.Eq(position)).Return(userTest, nil)

	userHandler.ChangeUserLocationHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserHandler_ChangeUSerPositionHandler_LoggerError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUcase := uMock.NewMockUserUsecase(ctrl)
	sessUcase := mock.NewMockSessionUsecase(ctrl)

	var byteData = bytes.NewReader([]byte(`
			{
			"latitude": 1,
			"longitude": 1,
			"radius": 1,
			"address":"Qwerty123"
			}
	`))

	position := &models.LocationRequest{
		Latitude:  1,
		Longitude: 1,
		Radius:    1,
		Address:   "Qwerty123",
	}

	r := httptest.NewRequest("POST", "/api/v1/user/position", byteData)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextUserID, userTest.ID)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	userHandler := NewUserHandler(userUcase, sessUcase)
	userHandler.Configure(router, nil)

	userUcase.EXPECT().UpdateLocation(gomock.Eq(userTest.ID), gomock.Eq(position)).Return(userTest, nil)

	userHandler.ChangeUserLocationHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserHandler_ChangeUSerPositionHandler_DecodeError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUcase := uMock.NewMockUserUsecase(ctrl)
	sessUcase := mock.NewMockSessionUsecase(ctrl)

	var byteData = bytes.NewReader([]byte(`
			{
			"latitude": "1",
			"longitude": 1,
			"radius": 1,
			"address":"Qwerty123"
			}
	`))

	r := httptest.NewRequest("POST", "/api/v1/user/position", byteData)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, userTest.ID)
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	userHandler := NewUserHandler(userUcase, sessUcase)
	userHandler.Configure(router, nil)

	userHandler.ChangeUserLocationHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_ChangeUSerPositionHandler_ValidationError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUcase := uMock.NewMockUserUsecase(ctrl)
	sessUcase := mock.NewMockSessionUsecase(ctrl)

	var byteData = bytes.NewReader([]byte(`
			{
			"latitude": -1123,
			"longitude": 1,
			"radius": 1,
			"address":"Qwerty123"
			}
	`))

	r := httptest.NewRequest("POST", "/api/v1/user/position", byteData)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, userTest.ID)
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	userHandler := NewUserHandler(userUcase, sessUcase)
	userHandler.Configure(router, nil)

	userHandler.ChangeUserLocationHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_ChangeUSerPositionHandler_NotAuth(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUcase := uMock.NewMockUserUsecase(ctrl)
	sessUcase := mock.NewMockSessionUsecase(ctrl)

	r := httptest.NewRequest("POST", "/api/v1/user/position", nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	userHandler := NewUserHandler(userUcase, sessUcase)
	userHandler.Configure(router, nil)

	userHandler.ChangeUserLocationHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUserHandler_ChangeUSerPositionHandler_NoUser(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUcase := uMock.NewMockUserUsecase(ctrl)
	sessUcase := mock.NewMockSessionUsecase(ctrl)

	var byteData = bytes.NewReader([]byte(`
			{
			"latitude": 1,
			"longitude": 1,
			"radius": 1,
			"address":"Qwerty123"
			}
	`))

	r := httptest.NewRequest("POST", "/api/v1/user/position", byteData)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, userTest.ID)
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	userHandler := NewUserHandler(userUcase, sessUcase)
	userHandler.Configure(router, nil)

	userUcase.EXPECT().UpdateLocation(gomock.Any(), gomock.Any()).Return(nil, errors.Cause(errors.UserNotExist))

	userHandler.ChangeUserLocationHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_UploadAvatarHandler_LoggerError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUcase := uMock.NewMockUserUsecase(ctrl)
	sessUcase := mock.NewMockSessionUsecase(ctrl)

	r := httptest.NewRequest("POST", "/api/v1/upload", nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextUserID, userTest.ID)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	userHandler := NewUserHandler(userUcase, sessUcase)
	userHandler.Configure(router, nil)

	userHandler.UploadAvatarHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_UploadAvatarHandler_ErrorContentType(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUcase := uMock.NewMockUserUsecase(ctrl)
	sessUcase := mock.NewMockSessionUsecase(ctrl)

	r := httptest.NewRequest("POST", "/api/v1/upload", nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, userTest.ID)
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	userHandler := NewUserHandler(userUcase, sessUcase)
	userHandler.Configure(router, nil)

	userHandler.UploadAvatarHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_UploadAvatarHandler_NoAuthError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUcase := uMock.NewMockUserUsecase(ctrl)
	sessUcase := mock.NewMockSessionUsecase(ctrl)

	r := httptest.NewRequest("POST", "/api/v1/upload", nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	userHandler := NewUserHandler(userUcase, sessUcase)
	userHandler.Configure(router, nil)

	userHandler.UploadAvatarHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
