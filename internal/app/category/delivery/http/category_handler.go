package http

import (
	"encoding/json"
	"net/http"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/middleware"
	"github.com/sirupsen/logrus"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/category"
	"github.com/gorilla/mux"
)

type CategoryHandler struct {
	categoryUcase category.CategoryUsecase
}

func NewCategoryHandler(productUcase category.CategoryUsecase) *CategoryHandler {
	return &CategoryHandler {
		categoryUcase: productUcase,
	}
}

func (cat *CategoryHandler) Configure(r *mux.Router, mw *middleware.Middleware) {
	r.HandleFunc("/categories",  mw.SetCSRFToken(cat.CategoriesHandler)).Methods(http.MethodGet, http.MethodOptions)
}

func (cat *CategoryHandler) CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value(middleware.ContextLogger).(*logrus.Entry)

	var categories []*models.Category
	categories, errCategories := cat.categoryUcase.GetAllCategories()
	if errCategories != nil {
		logger.Error(errCategories.Message)
		w.WriteHeader(errCategories.HttpError)
		w.Write(errors.JSONError(errCategories))
		return
	}

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
