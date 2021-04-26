package delivery

import (
	"encoding/json"
	errors2 "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
	logger2 "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/logger"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/sirupsen/logrus"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/middleware"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/search"
)

type SearchHandler struct {
	searchUsecase search.SearchUsecase
}

func NewSearchHandler(searchUsecase search.SearchUsecase) *SearchHandler {
	return &SearchHandler{
		searchUsecase: searchUsecase,
	}
}

func (sh *SearchHandler) Configure(r *mux.Router, mw *middleware.Middleware) {
	r.HandleFunc("/search", mw.SetCSRFToken(mw.CheckAuthMiddleware(sh.SearchHandler))).Methods(http.MethodGet, http.MethodOptions)
}

func (sh *SearchHandler) SearchHandler(w http.ResponseWriter, r *http.Request) {
	logger, ok := r.Context().Value(middleware.ContextLogger).(*logrus.Entry)
	if !ok {
		logger = logger2.GetDefaultLogger()
		logger.Warn("no logger")
	}

	logger.Debug("query ", r.URL.Query())

	search := &models.Search{}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(search, r.URL.Query())
	if err != nil {
		logger.Error(err)
		errE := errors2.UnexpectedBadRequest(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors2.JSONError(errE))
		return
	}
	logger.Info("search ", search)

	_, err = govalidator.ValidateStruct(search)
	if err != nil {
		if allErrs, ok := err.(govalidator.Errors); ok {
			logger.Error(allErrs.Errors())
			errE := errors2.UnexpectedBadRequest(allErrs)
			w.WriteHeader(errE.HttpError)
			w.Write(errors2.JSONError(errE))
			return
		}
	}

	userID, _ := r.Context().Value(middleware.ContextUserID).(uint64)
	logger.Info("user id ", userID)

	products, errE := sh.searchUsecase.SelectByFilter(&userID, search)
	if errE != nil {
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors2.JSONError(errE))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(products)
	if err != nil {
		logger.Error(err)
		errE := errors2.UnexpectedInternal(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors2.JSONError(errE))
		return
	}
}
