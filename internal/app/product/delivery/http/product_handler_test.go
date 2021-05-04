package http

import (
	"bytes"
	"context"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	mock "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/product/mocks"
	errors2 "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
	middleware2 "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/middleware"
	errors "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
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
	ctx = context.WithValue(ctx, middleware2.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware2.ContextUserID, uint64(1))
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
	ctx = context.WithValue(ctx, middleware2.ContextUserID, uint64(1))
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
	ctx = context.WithValue(ctx, middleware2.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware2.ContextUserID, uint64(1))
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
	ctx = context.WithValue(ctx, middleware2.ContextLogger, logrus.WithFields(logrus.Fields{
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
	ctx = context.WithValue(ctx, middleware2.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware2.ContextUserID, uint64(1))
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
	ctx = context.WithValue(ctx, middleware2.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	prodUcase.EXPECT().GetByID(gomock.Eq(uint64(0))).Return(prodTest, nil)

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

	prodUcase.EXPECT().GetByID(gomock.Eq(uint64(0))).Return(prodTest, nil)

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
	ctx = context.WithValue(ctx, middleware2.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	prodUcase.EXPECT().GetByID(gomock.Eq(uint64(0))).Return(nil, errors.Cause(errors.ProductNotExist))

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
	ctx = context.WithValue(ctx, middleware2.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware2.ContextUserID, userID)
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
	ctx = context.WithValue(ctx, middleware2.ContextUserID, userID)
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
	ctx = context.WithValue(ctx, middleware2.ContextLogger, logrus.WithFields(logrus.Fields{
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
	ctx = context.WithValue(ctx, middleware2.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware2.ContextUserID, userID)
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
	ctx = context.WithValue(ctx, middleware2.ContextUserID, userID)
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
	ctx = context.WithValue(ctx, middleware2.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware2.ContextUserID, uint64(1))
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
	ctx = context.WithValue(ctx, middleware2.ContextLogger, logrus.WithFields(logrus.Fields{
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
	ctx = context.WithValue(ctx, middleware2.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware2.ContextUserID, userID)
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
	ctx = context.WithValue(ctx, middleware2.ContextUserID, userID)
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
	ctx = context.WithValue(ctx, middleware2.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware2.ContextUserID, uint64(1))
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
	ctx = context.WithValue(ctx, middleware2.ContextLogger, logrus.WithFields(logrus.Fields{
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
	ctx = context.WithValue(ctx, middleware2.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware2.ContextUserID, uint64(1))
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
	ctx = context.WithValue(ctx, middleware2.ContextUserID, uint64(1))
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
	ctx = context.WithValue(ctx, middleware2.ContextLogger, logrus.WithFields(logrus.Fields{
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
	ctx = context.WithValue(ctx, middleware2.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware2.ContextUserID, uint64(1))
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
	ctx = context.WithValue(ctx, middleware2.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware2.ContextUserID, uint64(1))
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
	ctx = context.WithValue(ctx, middleware2.ContextUserID, uint64(1))
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
	ctx = context.WithValue(ctx, middleware2.ContextLogger, logrus.WithFields(logrus.Fields{
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
	ctx = context.WithValue(ctx, middleware2.ContextUserID, uint64(1))
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
	ctx = context.WithValue(ctx, middleware2.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware2.ContextUserID, uint64(1))
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
	ctx = context.WithValue(ctx, middleware2.ContextLogger, logrus.WithFields(logrus.Fields{
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
	ctx = context.WithValue(ctx, middleware2.ContextUserID, uint64(1))
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
	ctx = context.WithValue(ctx, middleware2.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware2.ContextUserID, uint64(1))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	rout := mux.NewRouter()
	router := rout.PathPrefix("/api/v1").Subrouter()
	prodHandler := NewProductHandler(prodUcase)
	prodHandler.Configure(router, rout, nil)

	prodHandler.PromoteProductHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
