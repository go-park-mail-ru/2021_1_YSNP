package http

import (
	"bytes"
	"context"
	"database/sql"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	mock "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/product/mocks"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/middleware"
)

var prodTest = &models.ProductData{
	ID:          0,
	Name:        "tovar",
	Date:        "",
	Amount:      10000,
	LinkImages:  nil,
	Description: "Description product aaaaa",
	Category:    "0",
	OwnerID:     1,
}

func TestProductHandler_ProductCreateHandler_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockProductUsecase(ctrl)

	var byteData = bytes.NewReader([]byte(`
			{
				"name":"tovar",
				"amount":10000,
				"description":"Description product aaaaa",
				"category":"0"
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
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	prodUcase.EXPECT().Create(gomock.Eq(prodTest)).Return(nil)

	prodHandler.ProductCreateHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_ProductCreateHandler_LoggerError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockProductUsecase(ctrl)

	var byteData = bytes.NewReader([]byte(`
			{
				"name":"tovar",
				"amount":10000,
				"description":"Description product aaaaa",
				"category":"0"
				}
	`))

	r := httptest.NewRequest("POST", "/product/create", byteData)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextUserID, uint64(1))
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := rout.PathPrefix("/api/v1").Subrouter()
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	prodUcase.EXPECT().Create(gomock.Eq(prodTest)).Return(nil)

	prodHandler.ProductCreateHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_ProductCreateHandler_DecodeError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockProductUsecase(ctrl)

	var byteData = bytes.NewReader([]byte(`
			{
				"name":"tovar",
				"amount":"hello",
				"description":"Description product aaaaa",
				"category":"0"
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
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	prodHandler.ProductCreateHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProductHandler_ProductCreateHandler_UserUnauth(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockProductUsecase(ctrl)

	r := httptest.NewRequest("POST", "/api/v1/product/create", nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := rout.PathPrefix("/api/v1").Subrouter()
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	prodHandler.ProductCreateHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestProductHandler_ProductCreateHandler_ValidateErr(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockProductUsecase(ctrl)

	var byteData = bytes.NewReader([]byte(`
			{
				"name":"tovar",
				"amount":-10000,
				"description":"Description product aaaaa",
				"category":"0"
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
	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	prodHandler.ProductCreateHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProductHandler_ProductIDHandler_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockProductUsecase(ctrl)

	r := httptest.NewRequest("GET", "/api/v1/product/0", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "0"})
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	prodUcase.EXPECT().GetProduct(gomock.Eq(uint64(0))).Return(prodTest, nil)

	prodHandler.ProductIDHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_ProductIDHandler_LoggerError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockProductUsecase(ctrl)

	r := httptest.NewRequest("GET", "/api/v1/product/0", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "0"})
	ctx := r.Context()
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	prodUcase.EXPECT().GetProduct(gomock.Eq(uint64(0))).Return(prodTest, nil)

	prodHandler.ProductIDHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_ProductIDHandler_ProductNotExist(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockProductUsecase(ctrl)

	r := httptest.NewRequest("GET", "/api/v1/product", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "0"})
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	prodUcase.EXPECT().GetProduct(gomock.Eq(uint64(0))).Return(nil, errors.Cause(errors.ProductNotExist))

	prodHandler.ProductIDHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestProductHandler_MainPageHandler_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockProductUsecase(ctrl)

	page := &models.Page{
		From:  0,
		Count: 20,
	}

	var userID uint64 = 1

	r := httptest.NewRequest("POST", "/api/v1/list/?from=0&count=20", nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, userID)
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	prodUcase.EXPECT().ListLatest(&userID, gomock.Eq(page)).Return([]*models.ProductListData{}, nil)

	prodHandler.MainPageHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_MainPageHandler_LoggerError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockProductUsecase(ctrl)

	page := &models.Page{
		From:  0,
		Count: 20,
	}

	var userID uint64 = 1

	r := httptest.NewRequest("POST", "/api/v1/list/?from=0&count=20", nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextUserID, userID)
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	prodUcase.EXPECT().ListLatest(&userID, gomock.Eq(page)).Return([]*models.ProductListData{}, nil)

	prodHandler.MainPageHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_MainPageHandler_DecodeError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockProductUsecase(ctrl)

	r := httptest.NewRequest("POST", "/api/v1/list/?from=fsfsffs&count=20", nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	prodHandler.MainPageHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProductHandler_UserAdHandler_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockProductUsecase(ctrl)

	page := &models.Page{
		From:  0,
		Count: 20,
	}

	var userID uint64 = 1

	r := httptest.NewRequest("POST", "/api/v1/list/?from=0&count=20", nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, userID)
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	prodUcase.EXPECT().UserAdList(userID, gomock.Eq(page)).Return([]*models.ProductListData{}, nil)

	prodHandler.UserAdHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_UserAdHandler_LoggerError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockProductUsecase(ctrl)

	page := &models.Page{
		From:  0,
		Count: 20,
	}

	var userID uint64 = 1

	r := httptest.NewRequest("POST", "/api/v1/list/?from=0&count=20", nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextUserID, userID)
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	prodUcase.EXPECT().UserAdList(userID, gomock.Eq(page)).Return([]*models.ProductListData{}, nil)

	prodHandler.UserAdHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_UserAdHandler_DecodeError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockProductUsecase(ctrl)

	r := httptest.NewRequest("POST", "/api/v1/list/?from=fsfsffs&count=20", nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, uint64(1))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	prodHandler.UserAdHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProductHandler_UserAdHandler_UserUnauth(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockProductUsecase(ctrl)

	r := httptest.NewRequest("POST", "/api/v1/list/?from=fsfsffs&count=20", nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	prodHandler.UserAdHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestProductHandler_UserFavoriteHandler_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockProductUsecase(ctrl)

	page := &models.Page{
		From:  0,
		Count: 20,
	}

	var userID uint64 = 1

	r := httptest.NewRequest("POST", "/api/v1/list/?from=0&count=20", nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, userID)
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	prodUcase.EXPECT().GetUserFavorite(userID, gomock.Eq(page)).Return([]*models.ProductListData{}, nil)

	prodHandler.UserFavoriteHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_UserFavoriteHandler_LoggerError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockProductUsecase(ctrl)

	page := &models.Page{
		From:  0,
		Count: 20,
	}

	var userID uint64 = 1

	r := httptest.NewRequest("POST", "/api/v1/list/?from=0&count=20", nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextUserID, userID)
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	prodUcase.EXPECT().GetUserFavorite(userID, gomock.Eq(page)).Return([]*models.ProductListData{}, nil)

	prodHandler.UserFavoriteHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_UserFavoriteHandler_DecodeError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockProductUsecase(ctrl)

	r := httptest.NewRequest("POST", "/api/v1/list/?from=fsfsffs&count=20", nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, uint64(1))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	prodHandler.UserFavoriteHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProductHandler_UserFavoriteHandler_UserUnauth(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockProductUsecase(ctrl)

	r := httptest.NewRequest("POST", "/api/v1/list/?from=fsfsffs&count=20", nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	prodHandler.UserFavoriteHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestProductHandler_LikeProductHandler_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockProductUsecase(ctrl)

	r := httptest.NewRequest("POST", "/api/v1/list/?from=fsfsffs&count=20", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "0"})
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, uint64(1))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	prodUcase.EXPECT().LikeProduct(gomock.Eq(uint64(1)), gomock.Eq(uint64(0))).Return(nil)

	prodHandler.LikeProductHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_LikeProductHandler_LoggerError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockProductUsecase(ctrl)

	r := httptest.NewRequest("POST", "/api/v1/list/?from=fsfsffs&count=20", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "0"})
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextUserID, uint64(1))
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	prodUcase.EXPECT().LikeProduct(gomock.Eq(uint64(1)), gomock.Eq(uint64(0))).Return(nil)

	prodHandler.LikeProductHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_LikeProductHandler_Unauth(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockProductUsecase(ctrl)

	r := httptest.NewRequest("POST", "/api/v1/list/?from=fsfsffs&count=20", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "0"})
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := rout.PathPrefix("/api/v1").Subrouter()
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	prodHandler.LikeProductHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestProductHandler_LikeProductHandler_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockProductUsecase(ctrl)

	r := httptest.NewRequest("POST", "/api/v1/list/?from=fsfsffs&count=20", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "0"})
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, uint64(1))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	prodUcase.EXPECT().LikeProduct(gomock.Eq(uint64(1)), gomock.Eq(uint64(0))).Return(errors.Cause(errors.ProductAlreadyLiked))

	prodHandler.LikeProductHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProductHandler_DislikeProductHandler_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockProductUsecase(ctrl)

	r := httptest.NewRequest("POST", "/api/v1/list/?from=fsfsffs&count=20", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "0"})
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, uint64(1))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := rout.PathPrefix("/api/v1").Subrouter()
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	prodUcase.EXPECT().DislikeProduct(gomock.Eq(uint64(1)), gomock.Eq(uint64(0))).Return(nil)

	prodHandler.DislikeProductHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_DislikeProductHandler_LoggerError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockProductUsecase(ctrl)

	r := httptest.NewRequest("POST", "/api/v1/list/?from=fsfsffs&count=20", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "0"})
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextUserID, uint64(1))
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := rout.PathPrefix("/api/v1").Subrouter()
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	prodUcase.EXPECT().DislikeProduct(gomock.Eq(uint64(1)), gomock.Eq(uint64(0))).Return(nil)

	prodHandler.DislikeProductHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_DislikeProductHandler_Unauth(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockProductUsecase(ctrl)

	r := httptest.NewRequest("POST", "/api/v1/list/?from=fsfsffs&count=20", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "0"})
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := rout.PathPrefix("/api/v1").Subrouter()
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	prodHandler.DislikeProductHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestProductHandler_UploadPhotoHandler_LoggerError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockProductUsecase(ctrl)

	r := httptest.NewRequest("POST", "/api/v1/product/1/upload", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "0"})
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextUserID, uint64(1))
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := rout.PathPrefix("/api/v1").Subrouter()
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	prodHandler.UploadPhotoHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProductHandler_UploadPhotoHandler_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockProductUsecase(ctrl)

	r := httptest.NewRequest("POST", "/api/v1/product/1/upload", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "0"})
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, uint64(1))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := rout.PathPrefix("/api/v1").Subrouter()
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	prodHandler.UploadPhotoHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProductHandler_UploadPhotoHandler__UserUnauthError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockProductUsecase(ctrl)

	r := httptest.NewRequest("POST", "/api/v1/product/1/upload", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "0"})
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := rout.PathPrefix("/api/v1").Subrouter()
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	prodHandler.UploadPhotoHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestProductHandler_PromoteProductHandler_LoggerError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockProductUsecase(ctrl)

	r := httptest.NewRequest("POST", "/api/v1/product/1/upload", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "0"})
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextUserID, uint64(1))
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := rout.PathPrefix("/api/v1").Subrouter()
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	prodHandler.PromoteProductHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProductHandler_PromoteProductHandler_PartFormError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockProductUsecase(ctrl)

	r := httptest.NewRequest("POST", "/api/v1/product/1/upload", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "0"})
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, uint64(1))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := rout.PathPrefix("/api/v1").Subrouter()
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	prodHandler.PromoteProductHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProductHandler_TrendsPageHandler_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockProductUsecase(ctrl)

	r := httptest.NewRequest("POST", "/api/v1/list/?from=fsfsffs&count=20", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "0"})
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, uint64(1))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := rout.PathPrefix("/api/v1").Subrouter()
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	var userID uint64 = 1;
	prodUcase.EXPECT().TrendList(&userID).Return([]*models.ProductListData{}, nil)

	prodHandler.TrendsPageHandler(w, r.WithContext(ctx))

	assert.Equal(t, http.StatusOK, w.Code)

	//error
	w = httptest.NewRecorder()
	prodUcase.EXPECT().TrendList(&userID).Return(nil, errors.UnexpectedInternal(sql.ErrConnDone))

	prodHandler.TrendsPageHandler(w, r.WithContext(ctx))

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestProductHandler_ProductCloseHandler_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockProductUsecase(ctrl)

	r := httptest.NewRequest("GET", "/api/v1/product/0", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "0"})
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextUserID, uint64(1))
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	prodUcase.EXPECT().Close(prodTest.ID, uint64(1)).Return(nil)

	prodHandler.ProductCloseHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusOK, w.Code)

	//error
	w = httptest.NewRecorder()
	prodUcase.EXPECT().Close(prodTest.ID, uint64(1)).Return(errors.UnexpectedInternal(sql.ErrConnDone))

	prodHandler.ProductCloseHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusInternalServerError, w.Code)

}

func TestProductHandler_ProductCloseHandler_Noauth(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockProductUsecase(ctrl)

	r := httptest.NewRequest("GET", "/api/v1/product/0", nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	prodHandler.ProductCloseHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestProductHandler_ProductEditHandler_Successt (t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockProductUsecase(ctrl)

	var byteData = bytes.NewReader([]byte(`
			{
				"name":"tovar",
				"amount":10000,
				"description":"Description product aaaaa",
				"category":"0",
				"ownerId": 1
				}
	`))

	r := httptest.NewRequest("POST", "/product/create", byteData)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, uint64(1))
	r = mux.SetURLVars(r, map[string]string{"id": "0"})
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := rout.PathPrefix("/api/v1").Subrouter()
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	prodUcase.EXPECT().Edit(gomock.Any()).Return(nil)

	prodHandler.ProductEditHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_ProductEditHandler_Unauth(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockProductUsecase(ctrl)


	r := httptest.NewRequest("POST", "/product/create", nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := rout.PathPrefix("/api/v1").Subrouter()
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	prodHandler.ProductEditHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusUnauthorized, w.Code)

}

func TestProductHandler_ProductEditHandler_Unmarsherr(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockProductUsecase(ctrl)

	var byteData = bytes.NewReader([]byte(`
			{
				"name":"tovar",
				"amount":10000,
				"description":"Description product aaaaa",
				"category":"0",
				"ownerId": 1
				
	`))

	r := httptest.NewRequest("POST", "/product/create", byteData)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, uint64(1))
	r = mux.SetURLVars(r, map[string]string{"id": "0"})
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := rout.PathPrefix("/api/v1").Subrouter()
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	prodHandler.ProductEditHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProductHandler_ProductEditHandler_WronOwner(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockProductUsecase(ctrl)

	var byteData = bytes.NewReader([]byte(`
			{
				"name":"tovar",
				"amount":10000,
				"description":"Description product aaaaa",
				"category":"0",
				"ownerId": 2
				}
	`))

	r := httptest.NewRequest("POST", "/product/create", byteData)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, uint64(1))
	r = mux.SetURLVars(r, map[string]string{"id": "0"})
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := rout.PathPrefix("/api/v1").Subrouter()
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	prodHandler.ProductEditHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusForbidden, w.Code)

}

func TestProductHandler_ProductEditHandler_ValidError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockProductUsecase(ctrl)

	var byteData = bytes.NewReader([]byte(`
			{
				"name":"tovar",
				"amount":-10000,
				"description":"Description product aaaaa",
				"category":"0",
				"ownerId": 1
				}
	`))

	r := httptest.NewRequest("POST", "/product/create", byteData)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, uint64(1))
	r = mux.SetURLVars(r, map[string]string{"id": "0"})
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := rout.PathPrefix("/api/v1").Subrouter()
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	prodHandler.ProductEditHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusBadRequest, w.Code)

}

func TestProductHandler_ProductEditHandler_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodUcase := mock.NewMockProductUsecase(ctrl)

	var byteData = bytes.NewReader([]byte(`
			{
				"name":"tovar",
				"amount":10000,
				"description":"Description product aaaaa",
				"category":"0",
				"ownerId": 1
				}
	`))

	r := httptest.NewRequest("POST", "/product/create", byteData)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, uint64(1))
	r = mux.SetURLVars(r, map[string]string{"id": "0"})
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := rout.PathPrefix("/api/v1").Subrouter()
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	prodUcase.EXPECT().Edit(gomock.Any()).Return(errors.UnexpectedInternal(sql.ErrConnDone))

	prodHandler.ProductEditHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusInternalServerError, w.Code)

}