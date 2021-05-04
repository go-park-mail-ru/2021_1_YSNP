package http

import (
	"context"
	"fmt"
	mock "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/category/mocks"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/middleware"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCategoryHandler_CategoriesHandler_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	catUcase := mock.NewMockCategoryUsecase(ctrl)

	r := httptest.NewRequest("GET", "/api/v1/categories", nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	catHandler := NewCategoryHandler(catUcase)
	catHandler.Configure(router, nil)

	catUcase.EXPECT().GetAllCategories().Return([]*models.Category{}, nil)

	catHandler.CategoriesHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCategoryHandler_CategoriesHandler_LoggerError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	catUcase := mock.NewMockCategoryUsecase(ctrl)

	r := httptest.NewRequest("GET", "/api/v1/categories", nil)
	ctx := r.Context()
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	catHandler := NewCategoryHandler(catUcase)
	catHandler.Configure(router, nil)

	catUcase.EXPECT().GetAllCategories().Return([]*models.Category{}, nil)

	catHandler.CategoriesHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCategoryHandler_CategoriesHandler_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	catUcase := mock.NewMockCategoryUsecase(ctrl)

	r := httptest.NewRequest("GET", "/api/v1/categories", nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	catHandler := NewCategoryHandler(catUcase)
	catHandler.Configure(router, nil)

	catUcase.EXPECT().GetAllCategories().Return(nil, errors.UnexpectedInternal(fmt.Errorf("")))

	catHandler.CategoriesHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}