package http

import (
	"context"
	mock "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/chat/mocks"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/middleware"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func TestChatHandler_CreateChat(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatUcase := mock.NewMockChatUsecase(ctrl)

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


	chatHandler.CreateChat(w, r.WithContext(ctx))
}