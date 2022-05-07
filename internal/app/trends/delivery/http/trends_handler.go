package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
	log "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/logger"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/middleware"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/trends"
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
	r.HandleFunc("/stat", mw.CheckAuthMiddleware(th.CreateTrends)).Methods(http.MethodPost, http.MethodOptions)
}

// CreateTrends godoc
// @Summary      Create trends
// @Description  Handler for creating trends
// @Tags         Trends
// @Accept       json
// @Produce      json
// @Param        body body models.UserInterested true "UserInterested"
// @Success      200 {object} errors.Success
// @Failure      400  {object}  errors.Error
// @Failure      404  {object}  errors.Error
// @Failure      500  {object}  errors.Error
// @Router      /stat [post]
func (th *TrendsHandler) CreateTrends(w http.ResponseWriter, r *http.Request) {
	logger, ok := r.Context().Value(middleware.ContextLogger).(*logrus.Entry)
	if !ok {
		logger = log.GetDefaultLogger()
		logger.Warn("no logger")
	}
	defer r.Body.Close()

	userID, ok := r.Context().Value(middleware.ContextUserID).(uint64)
	if !ok {
		errE := errors.Cause(errors.UserUnauthorized)
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	ui := &models.UserInterested{}
	err := json.NewDecoder(r.Body).Decode(&ui)
	if err != nil {
		logger.Error(err)
		errE := errors.UnexpectedBadRequest(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
	ui.UserID = userID
	th.trendsUsecase.InsertOrUpdate(ui)

	w.WriteHeader(http.StatusOK)
	w.Write(errors.JSONSuccess("Successful stat."))
}
