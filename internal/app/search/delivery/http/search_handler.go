package delivery

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/middleware"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/search"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
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
	r.HandleFunc("/search", sh.MainPageHandler).Methods(http.MethodPost)
}

func (sh *SearchHandler) MainPageHandler(w http.ResponseWriter, r *http.Request) {
	page := models.Search{}
	err := json.NewDecoder(r.Body).Decode(&page)
	if err != nil {
		logrus.Error(err)
		errE := errors.UnexpectedBadRequest(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	products, errE := sh.searchUsecase.SelectByFilter(&page)
	if errE != nil {
		logrus.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(products)
	if err != nil {
		logrus.Error(err)
		errE := errors.UnexpectedInternal(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
}