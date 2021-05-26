package delivery

import (
	"bytes"
	"context"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/middleware"
	mock "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/trends/mocks"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTrendsHandler_CreateTrends(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockTrendsUsecase(ctrl)

	ui := &models.UserInterested{
		UserID: 1,
		Text:   "aaaaaa",
	}

	var byteData = bytes.NewReader([]byte(`
			{
				"userID":0,
				"text":"aaaaaa"
				}
	`))

	r := httptest.NewRequest("POST", "/product/create", byteData)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, uint64(1))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := rout.PathPrefix("/api/v1").Subrouter()
	prodHandler := NewTrendsHandler(prodUcase)
	prodHandler.Configure(router, nil)

	prodUcase.EXPECT().InsertOrUpdate(ui).Return(nil)

	prodHandler.CreateTrends(w, r.WithContext(ctx))

	assert.Equal(t, http.StatusOK, w.Code)

	//error
	byteData = bytes.NewReader([]byte(`
			{
				"userID":0,
				"text":"aaaaaa"
	`))

	r = httptest.NewRequest("POST", "/product/create", byteData)
	w = httptest.NewRecorder()

	prodHandler.CreateTrends(w, r.WithContext(ctx))

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTrendsHandler_CreateTrends__NoLogger(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockTrendsUsecase(ctrl)

	ui := &models.UserInterested{
		UserID: 1,
		Text:   "aaaaaa",
	}

	var byteData = bytes.NewReader([]byte(`
			{
				"userID":0,
				"text":"aaaaaa"
				}
	`))

	r := httptest.NewRequest("POST", "/product/create", byteData)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextUserID, uint64(1))
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := rout.PathPrefix("/api/v1").Subrouter()
	prodHandler := NewTrendsHandler(prodUcase)
	prodHandler.Configure(router, nil)

	prodUcase.EXPECT().InsertOrUpdate(ui).Return(nil)

	prodHandler.CreateTrends(w, r.WithContext(ctx))

	assert.Equal(t, http.StatusOK, w.Code)

	//error
	byteData = bytes.NewReader([]byte(`
			{
				"userID":0,
				"text":"aaaaaa"
	`))

	r = httptest.NewRequest("POST", "/product/create", byteData)
	w = httptest.NewRecorder()

	prodHandler.CreateTrends(w, r.WithContext(ctx))

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTrendsHandler_CreateTrends_Unauth(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockTrendsUsecase(ctrl)

	r := httptest.NewRequest("POST", "/product/create", nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := rout.PathPrefix("/api/v1").Subrouter()
	prodHandler := NewTrendsHandler(prodUcase)
	prodHandler.Configure(router, nil)

	prodHandler.CreateTrends(w, r.WithContext(ctx))

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}