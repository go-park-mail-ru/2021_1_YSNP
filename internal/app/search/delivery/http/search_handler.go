package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/sirupsen/logrus"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/search"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
	log "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/logger"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/middleware"
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

// SearchHandler godoc
// @Summary      Search products
// @Description  Handler for searching products
// @Tags         Search
// @Accept       json
// @Produce      json
// @Param        From query integer false "From"
// @Param        Count query integer false "Count"
// @Param        Sorting query string false "Sorting"
// @Param        Search query string false "Search"
// @Param        Longitude query number false "Longitude"
// @Param        Latitude query number false "Latitude"
// @Param        Radius query integer false "Radius"
// @Param        FromAmount query integer false "FromAmount"
// @Param        ToAmount query integer false "ToAmount"
// @Param        Date query string false "Date"
// @Param        Category query string false "Category"
// @Success      200 {object} []models.ProductListData
// @Failure      400  {object}  errors.Error
// @Failure      404  {object}  errors.Error
// @Failure      500  {object}  errors.Error
// @Router      /search [get]
func (sh *SearchHandler) SearchHandler(w http.ResponseWriter, r *http.Request) {
	logger, ok := r.Context().Value(middleware.ContextLogger).(*logrus.Entry)
	if !ok {
		logger = log.GetDefaultLogger()
		logger.Warn("no logger")
	}

	logger.Debug("query ", r.URL.Query())

	search := &models.Search{}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(search, r.URL.Query())
	if err != nil {
		logger.Error(err)
		errE := errors.UnexpectedBadRequest(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
	logger.Info("search ", search)

	_, err = govalidator.ValidateStruct(search)
	if err != nil {
		if allErrs, ok := err.(govalidator.Errors); ok {
			logger.Error(allErrs.Errors())
			errE := errors.UnexpectedBadRequest(allErrs)
			w.WriteHeader(errE.HttpError)
			w.Write(errors.JSONError(errE))
			return
		}
	}

	userID, _ := r.Context().Value(middleware.ContextUserID).(uint64)
	logger.Info("user id ", userID)

	products, errE := sh.searchUsecase.SelectByFilter(&userID, search)
	if errE != nil {
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(products)
	if err != nil {
		logger.Error(err)
		errE := errors.UnexpectedInternal(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
}
