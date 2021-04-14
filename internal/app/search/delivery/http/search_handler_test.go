package delivery

import (
	"context"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/middleware"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	mock "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/search/mocks"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearchHandler_SearchHandler_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	searchUcase := mock.NewMockSearchUsecase(ctrl)

	var userID uint64 = 1

	search := &models.Search{
		Category:   "Шуба",
	}

	r := httptest.NewRequest("POST", "/api/v1/search?category=Шуба", nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, userID)
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := rout.PathPrefix("/api/v1").Subrouter()
	searchHandler := NewSearchHandler(searchUcase)
	searchHandler.Configure(router, nil)

	searchUcase.EXPECT().SelectByFilter(&userID, search).Return([]*models.ProductListData{}, nil)

	searchHandler.SearchHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSearchHandler_SearchHandler_ValidationError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	searchUcase := mock.NewMockSearchUsecase(ctrl)

	var userID uint64 = 1

	r := httptest.NewRequest("POST", "/api/v1/search?category=Шуба&fromAmount=-1", nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, userID)
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := rout.PathPrefix("/api/v1").Subrouter()
	searchHandler := NewSearchHandler(searchUcase)
	searchHandler.Configure(router, nil)

	searchHandler.SearchHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSearchHandler_SearchHandler_NotFoundErr(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	searchUcase := mock.NewMockSearchUsecase(ctrl)

	var userID uint64 = 1

	search := &models.Search{
		Category:   "Шуба",
	}

	r := httptest.NewRequest("POST", "/api/v1/search?category=Шуба", nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, userID)
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := rout.PathPrefix("/api/v1").Subrouter()
	searchHandler := NewSearchHandler(searchUcase)
	searchHandler.Configure(router, nil)

	searchUcase.EXPECT().SelectByFilter(&userID, search).Return(nil, errors.Cause(errors.EmptySearch))

	searchHandler.SearchHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSearchHandler_SearchHandler_QueryParseErr(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	searchUcase := mock.NewMockSearchUsecase(ctrl)

	var userID uint64 = 1

	r := httptest.NewRequest("POST", "/api/v1/search?category=Шуба&fromAmount=sfdfsdfsdf", nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, userID)
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := rout.PathPrefix("/api/v1").Subrouter()
	searchHandler := NewSearchHandler(searchUcase)
	searchHandler.Configure(router, nil)

	searchHandler.SearchHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
