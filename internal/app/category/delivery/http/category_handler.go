package http

import (
	"encoding/json"
	log "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/logger"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/middleware"
	"github.com/sirupsen/logrus"
	"net/http"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/category"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/errors"
	"github.com/gorilla/mux"
)

type CategoryHandler struct {
	categoryUcase category.CategoryUsecase
}

func NewCategoryHandler(productUcase category.CategoryUsecase) *CategoryHandler {
	return &CategoryHandler{
		categoryUcase: productUcase,
	}
}

func (cat *CategoryHandler) Configure(r *mux.Router, mw *middleware.Middleware) {
	r.HandleFunc("/categories", mw.SetCSRFToken(cat.CategoriesHandler)).Methods(http.MethodGet, http.MethodOptions)
}

func (cat *CategoryHandler) CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	logger, ok := r.Context().Value(middleware.ContextLogger).(*logrus.Entry)
	if !ok {
		logger = log.GetDefaultLogger()
		logger.Warn("no logger")
	}

	categories, errCategories := cat.categoryUcase.GetAllCategories()
	if errCategories != nil {
		logger.Error(errCategories.Message)
		w.WriteHeader(errCategories.HttpError)
		w.Write(errors.JSONError(errCategories))
		return
	}
	logger.Debug("categories", categories)

	body, err := json.Marshal(categories)
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
