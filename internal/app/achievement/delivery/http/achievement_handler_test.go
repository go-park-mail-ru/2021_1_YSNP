package http

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	mock "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/achievement/mocks"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/middleware"
)

func TestAchievementHandler_AchievementsSellerHandler_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	achUsecase := mock.NewMockAchievementUsecase(ctrl)

	r := httptest.NewRequest("GET", "/api/v1/achievements", nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	achHandler := NewAchievementHandler(achUsecase)
	achHandler.Configure(router, nil)

	achUsecase.EXPECT().GetUserAchievements(gomock.Eq(0)).Return([]*models.Achievement{}, nil)

	achHandler.achievementsSellerHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAchievementHandler_AchievementsSellerHandler_NoLogger(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	achUsecase := mock.NewMockAchievementUsecase(ctrl)

	r := httptest.NewRequest("GET", "/api/v1/achievements", nil)
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	achHandler := NewAchievementHandler(achUsecase)
	achHandler.Configure(router, nil)

	achUsecase.EXPECT().GetUserAchievements(gomock.Eq(0)).Return([]*models.Achievement{}, nil)

	achHandler.achievementsSellerHandler(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAchievementHandler_AchievementsSellerHandler_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	achUsecase := mock.NewMockAchievementUsecase(ctrl)

	r := httptest.NewRequest("GET", "/api/v1/achievements", nil)
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	achHandler := NewAchievementHandler(achUsecase)
	achHandler.Configure(router, nil)

	achUsecase.EXPECT().GetUserAchievements(gomock.Eq(0)).Return(nil, errors.UnexpectedInternal(fmt.Errorf("")))

	achHandler.achievementsSellerHandler(w, r)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAchievementHandler_AchievementsHandler_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	achUsecase := mock.NewMockAchievementUsecase(ctrl)

	r := httptest.NewRequest("GET", "/api/v1/achievements/{id:[0-9]+}", nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.ContextLogger, logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	}))
	logrus.SetOutput(ioutil.Discard)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	achHandler := NewAchievementHandler(achUsecase)
	achHandler.Configure(router, nil)

	achUsecase.EXPECT().GetUserAchievements(gomock.Eq(0)).Return([]*models.Achievement{}, nil)

	achHandler.achievementsHandler(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAchievementHandler_AchievementsHandler_NoLogger(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	achUsecase := mock.NewMockAchievementUsecase(ctrl)

	r := httptest.NewRequest("GET", "/api/v1/achievements/{id:[0-9]+}", nil)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	achHandler := NewAchievementHandler(achUsecase)
	achHandler.Configure(router, nil)

	achUsecase.EXPECT().GetUserAchievements(gomock.Eq(0)).Return([]*models.Achievement{}, nil)

	achHandler.achievementsHandler(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAchievementHandler_AchievementsHandler_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	achUsecase := mock.NewMockAchievementUsecase(ctrl)

	r := httptest.NewRequest("GET", "/api/v1/achievements/{id:[0-9]+}", nil)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	achHandler := NewAchievementHandler(achUsecase)
	achHandler.Configure(router, nil)

	achUsecase.EXPECT().GetUserAchievements(gomock.Eq(0)).Return(nil, errors.UnexpectedInternal(fmt.Errorf("")))

	achHandler.achievementsHandler(w, r)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}