package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/achievement"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
	log "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/logger"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/middleware"
)

type AchievementHandler struct {
	achUcase achievement.AchievementUsecase
}

func NewAchievementHandler(achievementUcase achievement.AchievementUsecase) *AchievementHandler {
	return &AchievementHandler{
		achUcase: achievementUcase,
	}
}

func (ah *AchievementHandler) Configure(r *mux.Router, mw *middleware.Middleware) {
	r.HandleFunc("/achievements",  mw.SetCSRFToken(mw.CheckAuthMiddleware(ah.achievementsHandler))).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/achievements/{id:[0-9]+}",  mw.SetCSRFToken(ah.achievementsSellerHandler)).Methods(http.MethodGet, http.MethodOptions)
}

func (ah *AchievementHandler) achievementsSellerHandler(w http.ResponseWriter, r *http.Request) {
	logger, ok := r.Context().Value(middleware.ContextLogger).(*logrus.Entry)
	if !ok {
		logger = log.GetDefaultLogger()
		logger.Warn("no logger")
	}

	userID, _ := strconv.Atoi(mux.Vars(r)["id"])
	logger.Info("user id ", userID)

	achievement, errE := ah.achUcase.GetUserAchievements(userID)

	if errE != nil {
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	body, err := json.Marshal(achievement)

	if err != nil {
		logger.Error(err)
		errE := errors.UnexpectedInternal(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (ah *AchievementHandler) achievementsHandler(w http.ResponseWriter, r *http.Request) {
	logger, ok := r.Context().Value(middleware.ContextLogger).(*logrus.Entry)
	if !ok {
		logger = log.GetDefaultLogger()
		logger.Warn("no logger")
	}

	userID, _ := r.Context().Value(middleware.ContextUserID).(uint64)
	logger.Info("user id ", userID)

	achievement, errE := ah.achUcase.GetUserAchievements(int(userID))

	if errE != nil {
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	body, err := json.Marshal(achievement)

	if err != nil {
		logger.Error(err)
		errE := errors.UnexpectedInternal(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
