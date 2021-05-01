package delivery

import (
	"encoding/json"
	"net/http"

	errors2 "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
	logger2 "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/logger"
	"github.com/sirupsen/logrus"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/middleware"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/trends"
	"github.com/gorilla/mux"
)

type TrendsHandler struct {
	trendsUsecase trends.TrendsUsecase

}

func NewTrendsHandler(trendsUsecase trends.TrendsUsecase) *TrendsHandler {
	return &TrendsHandler{
		trendsUsecase: trendsUsecase,
	}
}

func (th *TrendsHandler) Configure(r *mux.Router, mw *middleware.Middleware) {
	r.HandleFunc("/stat", mw.SetCSRFToken(mw.CheckAuthMiddleware(th.LogoutHandler))).Methods(http.MethodPost, http.MethodOptions)
}


func (th *TrendsHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	logger, ok := r.Context().Value(middleware.ContextLogger).(*logrus.Entry)
	if !ok {
		logger = logger2.GetDefaultLogger()
		logger.Warn("no logger")
	}
	defer r.Body.Close()

	userID, ok := r.Context().Value(middleware.ContextUserID).(uint64)
	if !ok {
		errE := errors2.Cause(errors2.UserUnauthorized)
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors2.JSONError(errE))
		return
	}

	ui := &models.UserInterested{}
	err := json.NewDecoder(r.Body).Decode(&ui)
	if err != nil {
		return
	}
	ui.UserID = userID
	th.trendsUsecase.InsertOrUpdate(ui)
	w.WriteHeader(http.StatusNoContent)
}
