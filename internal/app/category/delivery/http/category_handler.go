package http

import (
	"encoding/json"
	errors2 "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
	logger2 "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/logger"
	middleware2 "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/middleware"
	"github.com/sirupsen/logrus"
	"net/http"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/category"
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

func (cat *CategoryHandler) Configure(r *mux.Router, mw *middleware2.Middleware) {
	r.HandleFunc("/categories", mw.SetCSRFToken(cat.CategoriesHandler)).Methods(http.MethodGet, http.MethodOptions)
}

func (cat *CategoryHandler) CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	logger, ok := r.Context().Value(middleware2.ContextLogger).(*logrus.Entry)
	if !ok {
		logger = logger2.GetDefaultLogger()
		logger.Warn("no logger")
	}

	categories, errE := cat.categoryUcase.GetAllCategories()
	if errE != nil {
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors2.JSONError(errE))
		return
	}

	body, err := json.Marshal(categories)
	if err != nil {
		logger.Error(err)
		errE := errors2.UnexpectedInternal(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors2.JSONError(errE))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
