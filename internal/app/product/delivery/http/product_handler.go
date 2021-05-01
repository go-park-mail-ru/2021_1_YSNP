package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/microcosm-cc/bluemonday"
	"github.com/sirupsen/logrus"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/middleware"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/product"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
	log "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/logger"
)

type ProductHandler struct {
	productUcase product.ProductUsecase
}

func NewProductHandler(productUcase product.ProductUsecase) *ProductHandler {
	return &ProductHandler{
		productUcase: productUcase,
	}
}

func (ph *ProductHandler) Configure(r *mux.Router, rNoCSRF *mux.Router, mw *middleware.Middleware) {
	r.HandleFunc("/product/close", mw.CheckAuthMiddleware(ph.ProductCloseHandler)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/product/create", mw.CheckAuthMiddleware(ph.ProductCreateHandler)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/product/edit", mw.CheckAuthMiddleware(ph.ProductEditHandler)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/product/upload/{pid:[0-9]+}", mw.CheckAuthMiddleware(ph.UploadPhotoHandler)).Methods(http.MethodPost, http.MethodOptions)
	rNoCSRF.HandleFunc("/product/promote", ph.PromoteProductHandler).Methods(http.MethodPost, http.MethodOptions)

	r.HandleFunc("/product/{id:[0-9]+}", mw.SetCSRFToken(ph.ProductIDHandler)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/product/list", mw.SetCSRFToken(mw.CheckAuthMiddleware(ph.MainPageHandler))).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/user/ad/list", mw.SetCSRFToken(mw.CheckAuthMiddleware(ph.UserAdHandler))).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/user/favorite/list", mw.SetCSRFToken(mw.CheckAuthMiddleware(ph.UserFavoriteHandler))).Methods(http.MethodGet, http.MethodOptions)

	r.HandleFunc("/user/favorite/like/{id:[0-9]+}", mw.CheckAuthMiddleware(ph.LikeProductHandler)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/user/favorite/dislike/{id:[0-9]+}", mw.CheckAuthMiddleware(ph.DislikeProductHandler)).Methods(http.MethodPost, http.MethodOptions)
}


func (ph *ProductHandler) ProductEditHandler(w http.ResponseWriter, r *http.Request) {
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
	logger.Info("user id ", userID)

	productData := &models.ProductData{}
	err := json.NewDecoder(r.Body).Decode(&productData)
	if err != nil {
		logger.Error(err)
		errE := errors.UnexpectedBadRequest(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
	logger.Info("product data ", productData)

	if productData.OwnerID != userID {
		errE := errors.Cause(errors.WrongOwner)
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	sanitizer := bluemonday.UGCPolicy()
	productData.Name = sanitizer.Sanitize(productData.Name)
	productData.Description = sanitizer.Sanitize(productData.Description)
	productData.Category = sanitizer.Sanitize(productData.Category)
	logger.Debug("sanitize product data ", productData)

	_, err = govalidator.ValidateStruct(productData)
	if err != nil {
		if allErrs, ok := err.(govalidator.Errors); ok {
			logger.Error(allErrs.Errors())
			errE := errors.UnexpectedBadRequest(allErrs)
			w.WriteHeader(errE.HttpError)
			w.Write(errors.JSONError(errE))
			return
		}
	}

	errE := ph.productUcase.Edit(productData)
	if errE != nil {
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
	logger.Debug("product id ", productData.ID)

	w.WriteHeader(http.StatusOK)
	w.Write(errors.JSONSuccess("Successful creation.", productData.ID))
}

func (ph *ProductHandler) ProductCloseHandler(w http.ResponseWriter, r *http.Request) {
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
	logger.Info("user id ", userID)
	productData := &models.ProductData{}
	err := json.NewDecoder(r.Body).Decode(&productData)
	if err != nil {
		logger.Error(err)
		errE := errors.UnexpectedBadRequest(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
	logger.Info("product id ", productData.ID)

	errE := ph.productUcase.Close(productData, int(userID))
	if errE != nil {
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(errors.JSONSuccess("Successful close.", productData.ID))
}

func (ph *ProductHandler) ProductCreateHandler(w http.ResponseWriter, r *http.Request) {
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
	logger.Info("user id ", userID)

	productData := &models.ProductData{}
	err := json.NewDecoder(r.Body).Decode(&productData)
	if err != nil {
		logger.Error(err)
		errE := errors.UnexpectedBadRequest(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
	logger.Info("product data ", productData)

	sanitizer := bluemonday.UGCPolicy()
	productData.Name = sanitizer.Sanitize(productData.Name)
	productData.Description = sanitizer.Sanitize(productData.Description)
	productData.Category = sanitizer.Sanitize(productData.Category)
	logger.Debug("sanitize product data ", productData)

	_, err = govalidator.ValidateStruct(productData)
	if err != nil {
		if allErrs, ok := err.(govalidator.Errors); ok {
			logger.Error(allErrs.Errors())
			errE := errors.UnexpectedBadRequest(allErrs)
			w.WriteHeader(errE.HttpError)
			w.Write(errors.JSONError(errE))
			return
		}
	}

	errE := ph.productUcase.Create(productData)
	if errE != nil {
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
	logger.Debug("product id ", productData.ID)

	w.WriteHeader(http.StatusOK)
	w.Write(errors.JSONSuccess("Successful creation.", productData.ID))
}

func (ph *ProductHandler) UploadPhotoHandler(w http.ResponseWriter, r *http.Request) {
	logger, ok := r.Context().Value(middleware.ContextLogger).(*logrus.Entry)
	if !ok {
		logger = log.GetDefaultLogger()
		logger.Warn("no logger")
	}

	vars := mux.Vars(r)
	productID, _ := strconv.ParseUint(vars["pid"], 10, 64)
	logger.Info("product id ", productID)

	userId, ok := r.Context().Value(middleware.ContextUserID).(uint64)
	if !ok {
		errE := errors.Cause(errors.UserUnauthorized)
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
	logger.Info("user id ", userId)

	r.Body = http.MaxBytesReader(w, r.Body, 10*1024*1024)
	err := r.ParseMultipartForm(10 * 1024 * 1024)
	if err != nil {
		logger.Error(err)
		errE := errors.UnexpectedBadRequest(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	files := r.MultipartForm.File["photos"]
	_, errE := ph.productUcase.UpdatePhoto(productID, userId, files)
	if errE != nil {
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}


	w.WriteHeader(http.StatusOK)
	w.Write(errors.JSONSuccess("Successful upload."))
}

func (ph *ProductHandler) PromoteProductHandler(w http.ResponseWriter, r *http.Request) {
	logger, ok := r.Context().Value(middleware.ContextLogger).(*logrus.Entry)
	if !ok {
		logger = log.GetDefaultLogger()
		logger.Warn("no logger")
	}

	err := r.ParseForm()
	if err != nil {
		logger.Error(err)
		errE := errors.UnexpectedBadRequest(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	label := r.PostFormValue("label")
	if label == "" {
		errE := errors.Cause(errors.PromoteEmptyLabel)
		logger.Error(errE)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
	logger.Debug("label ", label)

	data := strings.Split(label, ",")
	productID, err := strconv.ParseUint(data[0], 10, 64)
	if err != nil {
		logger.Error(err)
		errE := errors.UnexpectedBadRequest(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
	logger.Info("product id ", productID)

	tariff, err := strconv.Atoi(data[1])
	if err != nil {
		logger.Error(err)
		errE := errors.UnexpectedBadRequest(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
	logger.Info("tariff ", tariff)

	errE := ph.productUcase.SetTariff(productID, tariff)
	if errE != nil {
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(errors.JSONSuccess("Successful promotion."))
}

func (ph *ProductHandler) ProductIDHandler(w http.ResponseWriter, r *http.Request) {
	logger, ok := r.Context().Value(middleware.ContextLogger).(*logrus.Entry)
	if !ok {
		logger = log.GetDefaultLogger()
		logger.Warn("no logger")
	}

	vars := mux.Vars(r)
	productID, _ := strconv.ParseUint(vars["id"], 10, 64)
	logger.Info("product id ", productID)

	product, errE := ph.productUcase.GetByID(productID)
	if errE != nil {
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
	logger.Debug("product by id ", product)

	body, err := json.Marshal(product)
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

func (ph *ProductHandler) MainPageHandler(w http.ResponseWriter, r *http.Request) {
	logger, ok := r.Context().Value(middleware.ContextLogger).(*logrus.Entry)
	if !ok {
		logger = log.GetDefaultLogger()
		logger.Warn("no logger")
	}

	page := &models.Page{}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(page, r.URL.Query())
	if err != nil {
		logger.Error(err)
		errE := errors.UnexpectedBadRequest(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
	logger.Info("page ", page)

	userID, _ := r.Context().Value(middleware.ContextUserID).(uint64)
	logger.Info("user id ", userID)

	products, errE := ph.productUcase.ListLatest(&userID, page)
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

func (ph *ProductHandler) UserAdHandler(w http.ResponseWriter, r *http.Request) {
	logger, ok := r.Context().Value(middleware.ContextLogger).(*logrus.Entry)
	if !ok {
		logger = log.GetDefaultLogger()
		logger.Warn("no logger")
	}

	userID, ok := r.Context().Value(middleware.ContextUserID).(uint64)
	if !ok {
		errE := errors.Cause(errors.UserUnauthorized)
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
	logger.Info("user id ", userID)

	page := &models.Page{}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(page, r.URL.Query())
	if err != nil {
		logger.Error(err)
		errE := errors.UnexpectedBadRequest(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
	logger.Info("page ", page)

	products, errE := ph.productUcase.UserAdList(userID, page)
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

func (ph *ProductHandler) UserFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	logger, ok := r.Context().Value(middleware.ContextLogger).(*logrus.Entry)
	if !ok {
		logger = log.GetDefaultLogger()
		logger.Warn("no logger")
	}

	userID, ok := r.Context().Value(middleware.ContextUserID).(uint64)
	if !ok {
		errE := errors.Cause(errors.UserUnauthorized)
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
	logger.Info("user id ", userID)

	page := &models.Page{}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(page, r.URL.Query())
	if err != nil {
		logger.Error(err)
		errE := errors.UnexpectedBadRequest(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
	logger.Info("page ", page)

	products, errE := ph.productUcase.GetUserFavorite(userID, page)
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

func (ph *ProductHandler) LikeProductHandler(w http.ResponseWriter, r *http.Request) {
	logger, ok := r.Context().Value(middleware.ContextLogger).(*logrus.Entry)
	if !ok {
		logger = log.GetDefaultLogger()
		logger.Warn("no logger")
	}

	userID, ok := r.Context().Value(middleware.ContextUserID).(uint64)
	if !ok {
		errE := errors.Cause(errors.UserUnauthorized)
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
	logger.Info("user id ", userID)

	vars := mux.Vars(r)
	productID, _ := strconv.ParseUint(vars["id"], 10, 64)
	logger.Info("product id ", productID)

	errE := ph.productUcase.LikeProduct(userID, productID)
	if errE != nil {
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(errors.JSONSuccess("Successful like."))
}

func (ph *ProductHandler) DislikeProductHandler(w http.ResponseWriter, r *http.Request) {
	logger, ok := r.Context().Value(middleware.ContextLogger).(*logrus.Entry)
	if !ok {
		logger = log.GetDefaultLogger()
		logger.Warn("no logger")
	}

	userID, ok := r.Context().Value(middleware.ContextUserID).(uint64)
	if !ok {
		errE := errors.Cause(errors.UserUnauthorized)
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
	logger.Info("user id ", userID)

	vars := mux.Vars(r)
	productID, _ := strconv.ParseUint(vars["id"], 10, 64)
	logger.Info("product id ", productID)

	errE := ph.productUcase.DislikeProduct(userID, productID)
	if errE != nil {
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(errors.JSONSuccess("Successful dislike."))
}
