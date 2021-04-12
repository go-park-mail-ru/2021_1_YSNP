package delivery

import (
	"bytes"
	"context"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/middleware"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	sessMock "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/session/mocks"
	userMock "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/user/mocks"
	_ "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/validator"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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
	}

func TestSessionHandler_LoginHandler_OK(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessUcase := sessMock.NewMockSessionUsecase(ctrl)
	userUcase := userMock.NewMockUserUsecase(ctrl)

	var byteData = bytes.NewReader([]byte(`
	{
		"telephone": "+79169230768",
		"password": "Qwerty12"
	}
	`))

	r := httptest.NewRequest("POST", "/api/v1/login", byteData)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	sessHandler := NewSessionHandler(sessUcase, userUcase)
	sessHandler.Configure(router, nil)

	userUcase.EXPECT().GetByTelephone(userTest.Telephone).Return(userTest, nil)
	userUcase.EXPECT().CheckPassword(userTest, userTest.Password).Return(nil)
	sessUcase.EXPECT().Create(gomock.Any()).Return(nil)

	sessHandler.LoginHandler(w, r.WithContext(ctx))

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSessionHandler_LoginHandler_InternalError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessUcase := sessMock.NewMockSessionUsecase(ctrl)
	userUcase := userMock.NewMockUserUsecase(ctrl)

	var byteData = bytes.NewReader([]byte(`
	{
		"telephone": "+79169230768",
		"password": "Qwerty12"
	}
	`))

	r := httptest.NewRequest("POST", "/api/v1/login", byteData)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	sessHandler := NewSessionHandler(sessUcase, userUcase)
	sessHandler.Configure(router, nil)

	userUcase.EXPECT().GetByTelephone(userTest.Telephone).Return(userTest, nil)
	userUcase.EXPECT().CheckPassword(userTest, userTest.Password).Return(nil)
	sessUcase.EXPECT().Create(gomock.Any()).Return(errors.Cause(errors.InternalError))

	sessHandler.LoginHandler(w, r.WithContext(ctx))

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestSessionHandler_LoginHandler_ValidationError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessUcase := sessMock.NewMockSessionUsecase(ctrl)
	userUcase := userMock.NewMockUserUsecase(ctrl)

	var byteData = bytes.NewReader([]byte(`
	{
		"telephone": "+7916923076",
		"password": "Qwerty12"
	}
	`)) // short telephone number

	r := httptest.NewRequest("POST", "/api/v1/login", byteData)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	sessHandler := NewSessionHandler(sessUcase, userUcase)
	sessHandler.Configure(router, nil)

	sessHandler.LoginHandler(w, r.WithContext(ctx))

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSessionHandler_LoginHandler_UserNotExist(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessUcase := sessMock.NewMockSessionUsecase(ctrl)
	userUcase := userMock.NewMockUserUsecase(ctrl)

	var byteData = bytes.NewReader([]byte(`
	{
		"telephone": "+79169230768",
		"password": "Qwerty12"
	}
	`))

	r := httptest.NewRequest("POST", "/api/v1/login", byteData)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	sessHandler := NewSessionHandler(sessUcase, userUcase)
	sessHandler.Configure(router, nil)

	userUcase.EXPECT().GetByTelephone(userTest.Telephone).Return(nil, errors.Cause(errors.UserNotExist))

	sessHandler.LoginHandler(w, r.WithContext(ctx))

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSessionHandler_LoginHandler_WrongPassword(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessUcase := sessMock.NewMockSessionUsecase(ctrl)
	userUcase := userMock.NewMockUserUsecase(ctrl)

	var byteData = bytes.NewReader([]byte(`
	{
		"telephone": "+79169230768",
		"password": "Qwerty12"
	}
	`))

	r := httptest.NewRequest("POST", "/api/v1/login", byteData)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	sessHandler := NewSessionHandler(sessUcase, userUcase)
	sessHandler.Configure(router, nil)

	userUcase.EXPECT().GetByTelephone(userTest.Telephone).Return(userTest, nil)
	userUcase.EXPECT().CheckPassword(userTest, userTest.Password).Return(errors.Cause(errors.WrongPassword))

	sessHandler.LoginHandler(w, r.WithContext(ctx))

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSessionHandler_LogoutHandler_OK(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessUcase := sessMock.NewMockSessionUsecase(ctrl)
	userUcase := userMock.NewMockUserUsecase(ctrl)

	session := models.CreateSession(0)
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    session.Value,
		Expires:  session.ExpiresAt,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
	}

	r := httptest.NewRequest("POST", "/api/v1/logout", nil)
	r.AddCookie(&cookie)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	sessHandler := NewSessionHandler(sessUcase, userUcase)
	sessHandler.Configure(router, nil)

	sessUcase.EXPECT().Delete(session.Value).Return(nil)

	sessHandler.LogoutHandler(w, r.WithContext(ctx))

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSessionHandler_LogoutHandler_NoCookie(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessUcase := sessMock.NewMockSessionUsecase(ctrl)
	userUcase := userMock.NewMockUserUsecase(ctrl)

	r := httptest.NewRequest("POST", "/api/v1/logout", nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	sessHandler := NewSessionHandler(sessUcase, userUcase)
	sessHandler.Configure(router, nil)

	sessHandler.LogoutHandler(w, r.WithContext(ctx))

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestSessionHandler_LogoutHandler_SessionNotExist(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessUcase := sessMock.NewMockSessionUsecase(ctrl)
	userUcase := userMock.NewMockUserUsecase(ctrl)

	session := models.CreateSession(0)
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    session.Value,
		Expires:  session.ExpiresAt,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
	}

	r := httptest.NewRequest("POST", "/api/v1/logout", nil)
	r.AddCookie(&cookie)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	sessHandler := NewSessionHandler(sessUcase, userUcase)
	sessHandler.Configure(router, nil)

	sessUcase.EXPECT().Delete(session.Value).Return(errors.Cause(errors.SessionNotExist))

	sessHandler.LogoutHandler(w, r.WithContext(ctx))

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}