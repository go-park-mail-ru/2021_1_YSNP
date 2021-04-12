package delivery

import (
	"bytes"
	"context"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/middleware"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	sMock "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/session/mocks"
	uMock "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/user/mocks"
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

//func TestUserHandler_SignUpHandler_OK(t *testing.T) {
//	t.Parallel()
//
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	userUcase := uMock.NewMockUserUsecase(ctrl)
//	sessUcase := sMock.NewMockSessionUsecase(ctrl)
//
//	var byteData = bytes.NewReader([]byte(`
//			{
//			"name":"Максим",
//			"surname":"Торжков",
//			"sex":"male",
//			"email":"a@a.ru",
//			"telephone":"+79169230768",
//			"password1":"Qwerty12",
//			"password2":"Qwerty12",
//			"dateBirth":"2021-03-08"
//			}
//	`))
//
//	userTest := &models.UserData{
//		ID:         0,
//		Name:       "Максим",
//		Surname:    "Торжков",
//		Sex:        "male",
//		Email:      "a@a.ru",
//		Telephone:  "+79169230768",
//		Password:   "Qwerty12",
//		DateBirth:  "2021-03-08",
//	}
//
//	session := models.CreateSession(userTest.ID)
//
//	r := httptest.NewRequest("POST", "/api/v1/signup", byteData)
//	ctx := r.Context()
//	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
//		"logger": "LOGRUS",
//	}))
//	logrus.SetOutput(ioutil.Discard)
//	w := httptest.NewRecorder()
//
//	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
//	userHandler := NewUserHandler(userUcase, sessUcase)
//	userHandler.Configure(router, nil)
//
//	userUcase.EXPECT().Create(userTest).Return()
//	sessUcase.EXPECT().Create(session).Return()
//	userHandler.SignUpHandler(w, r.WithContext(ctx))
//
//	assert.Equal(t, http.StatusBadRequest, w.Code)
//}

func TestUserHandler_SignUpHandler_TelephoneAlreadyExists(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUcase := uMock.NewMockUserUsecase(ctrl)
	sessUcase := sMock.NewMockSessionUsecase(ctrl)

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
		ID:         0,
		Name:       "Максим",
		Surname:    "Торжков",
		Sex:        "male",
		Email:      "a@a.ru",
		Telephone:  "+79169230768",
		Password:   "Qwerty12",
		DateBirth:  "2021-03-08",
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