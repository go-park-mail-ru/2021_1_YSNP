package http

import (
	"encoding/json"
	"net/http"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/middleware"
	"github.com/microcosm-cc/bluemonday"
	"github.com/sirupsen/logrus"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/category"
	"github.com/gorilla/mux"
)

type CategoryHandler struct {
	categoryUcase category.CategoryUsecase
}

func NewProductHandler(productUcase category.CategoryUsecase) *CategoryHandler {
	return &CategoryUsecase {
		categoryUcase: productUcase,
	}
}

func (cat *CategoryHandler) Configure(r *mux.Router, mw *middleware.Middleware) {
	r.HandleFunc("/categories", cat.CategoriesHandler).Methods(http.MethodGet)
}

func (cat *CategoryHandler) CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	
	var categories []*models.Category
	categories = append(categories, &models.Category {Title: "Транспорт"}, &models.Category {Title: "Недвижмость"}, &models.Category {Title: "Хобби и отдых"}, &models.Category {Title: "Работа"}, &models.Category {Title: "Для дома и дачи"}, &models.Category {Title: "Бытовая электрика"}, &models.Category {Title: "Личные вещи"}, &models.Category {Title: "Животные"})
	body, err := json.Marshal(categories)
	if err != nil {
		logrus.Error(err)
		errE := errors.UnexpectedInternal(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
