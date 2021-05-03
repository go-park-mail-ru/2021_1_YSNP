package http

import (
	"bytes"
	"context"
	"database/sql"
	mock "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/chat/mocks"
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

func TestChatHandler_CreateChat_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatUcase := mock.NewMockChatUsecase(ctrl)

	var byteData = bytes.NewReader([]byte(`
	{
		"productID": 2,
  		"partnerID": 2
	}
	`))

	req := &models.ChatCreateReq{ProductID: uint64(2), PartnerID: uint64(2)}

	r := httptest.NewRequest("POST", "/api/v1/login", byteData)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, uint64(1))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	chatHandler := NewChatHandler(chatUcase)
	chatHandler.Configure(router, nil, nil)

	chatUcase.EXPECT().CreateChat(req, uint64(1)).Return(&models.ChatResponse{}, nil)

	chatHandler.CreateChat(w, r.WithContext(ctx))

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestChatHandler_CreateChat_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatUcase := mock.NewMockChatUsecase(ctrl)

	var byteData = bytes.NewReader([]byte(`
	{
		"productID": 2,
  		"partnerID": 2
	}
	`))

	req := &models.ChatCreateReq{ProductID: uint64(2), PartnerID: uint64(2)}

	r := httptest.NewRequest("POST", "/api/v1/login", byteData)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, uint64(1))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	chatHandler := NewChatHandler(chatUcase)
	chatHandler.Configure(router, nil, nil)

	chatUcase.EXPECT().CreateChat(req, uint64(1)).Return(nil, errors.UnexpectedInternal(sql.ErrNoRows))

	chatHandler.CreateChat(w, r.WithContext(ctx))

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestChatHandler_CreateChat_NoAuth(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatUcase := mock.NewMockChatUsecase(ctrl)

	r := httptest.NewRequest("POST", "/api/v1/chat", nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	chatHandler := NewChatHandler(chatUcase)
	chatHandler.Configure(router, nil, nil)

	chatHandler.CreateChat(w, r.WithContext(ctx))

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestChatHandler_GetChatByID_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatUcase := mock.NewMockChatUsecase(ctrl)

	r := httptest.NewRequest("POST", "/api/v1/login", nil)
	r = mux.SetURLVars(r, map[string]string{"cid": "0"})
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, uint64(1))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	chatHandler := NewChatHandler(chatUcase)
	chatHandler.Configure(router, nil, nil)

	chatUcase.EXPECT().GetChatById(uint64(0), uint64(1)).Return(&models.ChatResponse{}, nil)

	chatHandler.GetChatByID(w, r.WithContext(ctx))

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestChatHandler_GetChatByID_NoAuth(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatUcase := mock.NewMockChatUsecase(ctrl)

	r := httptest.NewRequest("POST", "/api/v1/login", nil)
	r = mux.SetURLVars(r, map[string]string{"cid": "0"})
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	chatHandler := NewChatHandler(chatUcase)
	chatHandler.Configure(router, nil, nil)

	chatHandler.GetChatByID(w, r.WithContext(ctx))

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestChatHandler_GetChatByID_Nochat(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatUcase := mock.NewMockChatUsecase(ctrl)

	r := httptest.NewRequest("POST", "/api/v1/login", nil)
	r = mux.SetURLVars(r, map[string]string{"cid": "0"})
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, uint64(1))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	chatHandler := NewChatHandler(chatUcase)
	chatHandler.Configure(router, nil, nil)

	chatUcase.EXPECT().GetChatById(uint64(0), uint64(1)).Return(nil, errors.Cause(errors.ChatNotExist))

	chatHandler.GetChatByID(w, r.WithContext(ctx))

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestChatHandler_GetUserChats_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatUcase := mock.NewMockChatUsecase(ctrl)

	r := httptest.NewRequest("POST", "/api/v1/login", nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, uint64(1))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	chatHandler := NewChatHandler(chatUcase)
	chatHandler.Configure(router, nil, nil)

	chatUcase.EXPECT().GetUserChats(uint64(1)).Return([]*models.ChatResponse{}, nil)

	chatHandler.GetUserChats(w, r.WithContext(ctx))

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestChatHandler_GetUserChats_NoAuth(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatUcase := mock.NewMockChatUsecase(ctrl)

	r := httptest.NewRequest("POST", "/api/v1/login", nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	chatHandler := NewChatHandler(chatUcase)
	chatHandler.Configure(router, nil, nil)

	chatHandler.GetUserChats(w, r.WithContext(ctx))

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestChatHandler_GetUserChats_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatUcase := mock.NewMockChatUsecase(ctrl)

	r := httptest.NewRequest("POST", "/api/v1/login", nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, uint64(1))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	chatHandler := NewChatHandler(chatUcase)
	chatHandler.Configure(router, nil, nil)

	chatUcase.EXPECT().GetUserChats(uint64(1)).Return(nil, errors.UnexpectedInternal(sql.ErrNoRows))

	chatHandler.GetUserChats(w, r.WithContext(ctx))

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestChatHandler_ServeWs_NoAuth(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatUcase := mock.NewMockChatUsecase(ctrl)
	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	chatHandler := NewChatHandler(chatUcase)
	chatHandler.Configure(router, nil, nil)

	r := httptest.NewRequest("POST", "/api/v1/login", nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	chatHandler.ServeWs(nil)(w, r.WithContext(ctx))

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestChatHandler_ServeWs_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatUcase := mock.NewMockChatUsecase(ctrl)
	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	chatHandler := NewChatHandler(chatUcase)
	chatHandler.Configure(router, nil, nil)

	r := httptest.NewRequest("POST", "/api/v1/login", nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	ctx = context.WithValue(ctx, middleware.ContextUserID, uint64(1))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	chatHandler.ServeWs(nil)(w, r.WithContext(ctx))

	assert.Equal(t, http.StatusBadRequest, w.Code)
}