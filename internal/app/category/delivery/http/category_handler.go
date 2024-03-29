package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/category"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
	log "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/logger"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/middleware"
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

// CategoriesHandler godoc
// @Summary      get categories
// @Description  Handler for getting categories
// @Tags         Category
// @Accept       json
// @Produce      json
// @Success      200 {object} []models.Category
// @Failure      400  {object}  errors.Error
// @Failure      404  {object}  errors.Error
// @Failure      500  {object}  errors.Error
// @Router       /categories [get]
func (cat *CategoryHandler) CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	logger, ok := r.Context().Value(middleware.ContextLogger).(*logrus.Entry)
	if !ok {
		logger = log.GetDefaultLogger()
		logger.Warn("no logger")
	}

	categories, errE := cat.categoryUcase.GetAllCategories()
	if errE != nil {
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
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
