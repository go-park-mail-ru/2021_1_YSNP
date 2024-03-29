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

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/product"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
	log "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/logger"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/middleware"
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
	r.HandleFunc("/product/create", mw.CheckAuthMiddleware(ph.ProductCreateHandler)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/product/close/{id:[0-9]+}", mw.CheckAuthMiddleware(ph.ProductCloseHandler)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/product/edit", mw.CheckAuthMiddleware(ph.ProductEditHandler)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/product/upload/{pid:[0-9]+}", mw.CheckAuthMiddleware(ph.UploadPhotoHandler)).Methods(http.MethodPost, http.MethodOptions)
	rNoCSRF.HandleFunc("/product/promote", ph.PromoteProductHandler).Methods(http.MethodPost, http.MethodOptions)

	r.HandleFunc("/product/{id:[0-9]+}", mw.SetCSRFToken(ph.ProductIDHandler)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/product/list", mw.SetCSRFToken(mw.CheckAuthMiddleware(ph.MainPageHandler))).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/user/{id:[0-9]+}/ad/list", mw.SetCSRFToken(ph.SellerAdHandler)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/user/ad/list", mw.SetCSRFToken(mw.CheckAuthMiddleware(ph.UserAdHandler))).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/user/favorite/list", mw.SetCSRFToken(mw.CheckAuthMiddleware(ph.UserFavoriteHandler))).Methods(http.MethodGet, http.MethodOptions)

	r.HandleFunc("/user/favorite/like/{id:[0-9]+}", mw.CheckAuthMiddleware(ph.LikeProductHandler)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/user/favorite/dislike/{id:[0-9]+}", mw.CheckAuthMiddleware(ph.DislikeProductHandler)).Methods(http.MethodPost, http.MethodOptions)

	r.HandleFunc("/product/trend/list", mw.SetCSRFToken(mw.CheckAuthMiddleware(ph.TrendHandler))).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/product/{id:[0-9]+}/trend/list", mw.SetCSRFToken(ph.ProductTrendsHandler)).Methods(http.MethodGet, http.MethodOptions)

	r.HandleFunc("/product/buyer/{id:[0-9]+}", mw.CheckAuthMiddleware(ph.SetProductBuyer)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/product/review/{id:[0-9]+}", mw.CheckAuthMiddleware(ph.CreateProductReview)).Methods(http.MethodPost, http.MethodOptions)

	r.HandleFunc("/user/{id:[0-9]+}/reviews/{type:seller|buyer}", mw.SetCSRFToken(ph.GetUserReviews)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/user/reviews/await/{type:seller|buyer}", mw.SetCSRFToken(mw.CheckAuthMiddleware(ph.GetWaitingReviews))).Methods(http.MethodGet, http.MethodOptions)
}

// ProductCreateHandler godoc
// @Summary      create product
// @Description  Handler for creating new product
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        body body models.ProductData true "Product"
// @Success      200 {object} errors.Success
// @Failure      400  {object}  errors.Error
// @Failure      404  {object}  errors.Error
// @Failure      500  {object}  errors.Error
// @Router       /product/create [post]
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
	productData.OwnerID = userID

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

// ProductCloseHandler godoc
// @Summary      close product
// @Description  Handler for closing new product
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        id path int64 true "Product ID"
// @Success      200 {object} []models.UserData
// @Failure      400  {object}  errors.Error
// @Failure      404  {object}  errors.Error
// @Failure      500  {object}  errors.Error
// @Router       /product/close/{id} [post]
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

	vars := mux.Vars(r)
	productID, _ := strconv.ParseUint(vars["id"], 10, 64)
	logger.Info("product id ", productID)

	errE := ph.productUcase.Close(productID, userID)
	if errE != nil {
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	users, errE := ph.productUcase.GetProductReviewers(productID, userID)
	if errE != nil {
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	body, err := json.Marshal(users)
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

// ProductEditHandler godoc
// @Summary      edit product
// @Description  Handler for editing new product
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        body body models.ProductData true "Product"
// @Success      200 {object} errors.Success
// @Failure      400  {object}  errors.Error
// @Failure      404  {object}  errors.Error
// @Failure      500  {object}  errors.Error
// @Router       /product/edit [post]
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
	w.Write(errors.JSONSuccess("Successful update.", productData.ID))
}

// UploadPhotoHandler godoc
// @Summary      upload photos for products
// @Description  Handler for uploading a product's photos
// @Tags         Product
// @Accept       mpfd
// @Produce      json
// @Param        pid path int64 true "Product ID"
// @Param 		 photos formData file true "Image to upload"
// @Success      200  {object}  errors.Success
// @Failure      400  {object}  errors.Error
// @Failure      404  {object}  errors.Error
// @Failure      500  {object}  errors.Error
// @Router       /product/upload/{pid} [post]
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

// PromoteProductHandler godoc
// @Summary      promote product
// @Description  Handler for promoting new product
// @Tags         Product
// @Accept       mpfd
// @Produce      json
// @Param 		 photos formData file true "Image to upload"
// @Success      200 {object} errors.Success
// @Failure      400  {object}  errors.Error
// @Failure      404  {object}  errors.Error
// @Failure      500  {object}  errors.Error
// @Router       /product/promote [post]
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

// ProductIDHandler godoc
// @Summary      get product by id
// @Description  Handler for getting product by id
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        id path int64 true "Product ID"
// @Success      200 {object} models.ProductData
// @Failure      400  {object}  errors.Error
// @Failure      404  {object}  errors.Error
// @Failure      500  {object}  errors.Error
// @Router       /product/{id} [get]
func (ph *ProductHandler) ProductIDHandler(w http.ResponseWriter, r *http.Request) {
	logger, ok := r.Context().Value(middleware.ContextLogger).(*logrus.Entry)
	if !ok {
		logger = log.GetDefaultLogger()
		logger.Warn("no logger")
	}

	vars := mux.Vars(r)
	productID, _ := strconv.ParseUint(vars["id"], 10, 64)
	logger.Info("product id ", productID)

	product, errE := ph.productUcase.GetProduct(productID)
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

// MainPageHandler godoc
// @Summary      list of products
// @Description  Handler for getting list of products
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        From query integer false "From"
// @Param        Count query integer false "Count"
// @Success      200 {object} []models.ProductListData
// @Failure      400  {object}  errors.Error
// @Failure      404  {object}  errors.Error
// @Failure      500  {object}  errors.Error
// @Router       /product/list [get]
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

// UserAdHandler godoc
// @Summary      list of user ad products
// @Description  Handler for getting list of user ad products
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        From query integer false "From"
// @Param        Count query integer false "Count"
// @Success      200 {object} []models.ProductListData
// @Failure      400  {object}  errors.Error
// @Failure      404  {object}  errors.Error
// @Failure      500  {object}  errors.Error
// @Router       /user/ad/list [get]
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

// SellerAdHandler godoc
// @Summary      list of seller ad products
// @Description  Handler for getting list of seller ad products
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        From query integer false "From"
// @Param        Count query integer false "Count"
// @Param        id path int64 true "User ID"
// @Success      200 {object} []models.ProductListData
// @Failure      400  {object}  errors.Error
// @Failure      404  {object}  errors.Error
// @Failure      500  {object}  errors.Error
// @Router       /user/{id}/ad/list [get]
func (ph *ProductHandler) SellerAdHandler(w http.ResponseWriter, r *http.Request) {
	logger, ok := r.Context().Value(middleware.ContextLogger).(*logrus.Entry)
	if !ok {
		logger = log.GetDefaultLogger()
		logger.Warn("no logger")
	}

	vars := mux.Vars(r)
	userID, _ := strconv.ParseUint(vars["id"], 10, 64)
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

// UserFavoriteHandler godoc
// @Summary      list of user favorite products
// @Description  Handler for getting list of user favorite products
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        From query integer false "From"
// @Param        Count query integer false "Count"
// @Success      200 {object} []models.ProductListData
// @Failure      400  {object}  errors.Error
// @Failure      404  {object}  errors.Error
// @Failure      500  {object}  errors.Error
// @Router       /user/favorite/list [get]
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

// LikeProductHandler godoc
// @Summary      add product to favorites
// @Description  Handler for adding product to favorites
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        id path int64 true "Product ID"
// @Success      200 {object} errors.Success
// @Failure      400  {object}  errors.Error
// @Failure      404  {object}  errors.Error
// @Failure      500  {object}  errors.Error
// @Router       /user/favorite/like/{id} [post]
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

// DislikeProductHandler godoc
// @Summary      remove product to favorites
// @Description  Handler for removing product from favorites
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        id path int64 true "Product ID"
// @Success      200 {object} errors.Success
// @Failure      400  {object}  errors.Error
// @Failure      404  {object}  errors.Error
// @Failure      500  {object}  errors.Error
// @Router       /user/favorite/dislike/{id} [post]
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

// TrendHandler godoc
// @Summary      get recommendation to user
// @Description  Handler for getting recommendation to user
// @Tags         Product
// @Accept       json
// @Produce      json
// @Success      200 {object} []models.ProductListData
// @Failure      400  {object}  errors.Error
// @Failure      404  {object}  errors.Error
// @Failure      500  {object}  errors.Error
// @Router       /product/trend/list [get]
func (ph *ProductHandler) TrendHandler(w http.ResponseWriter, r *http.Request) {
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

	products, errE := ph.productUcase.TrendList(&userID)
	if errE != nil {
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(products)
	if err != nil {
		logger.Error(err)
		errE := errors.UnexpectedInternal(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
}

// ProductTrendsHandler godoc
// @Summary      get recommendation to product
// @Description  Handler for getting recommendation to product
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        id path int64 true "Product ID"
// @Success      200 {object} []models.ProductListData
// @Failure      400  {object}  errors.Error
// @Failure      404  {object}  errors.Error
// @Failure      500  {object}  errors.Error
// @Router       /product/{id}/trend/list [get]
func (ph *ProductHandler) ProductTrendsHandler(w http.ResponseWriter, r *http.Request) {
	logger, ok := r.Context().Value(middleware.ContextLogger).(*logrus.Entry)
	if !ok {
		logger = log.GetDefaultLogger()
		logger.Warn("no logger")
	}

	userID, _ := r.Context().Value(middleware.ContextUserID).(uint64)
	logger.Info("user id ", userID)

	vars := mux.Vars(r)
	productID, _ := strconv.ParseUint(vars["id"], 10, 64)
	logger.Info("product id ", productID)

	products, errE := ph.productUcase.RecommendationList(productID, userID)
	if errE != nil {
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(products)
	if err != nil {
		logger.Error(err)
		errE := errors.UnexpectedInternal(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
}

// SetProductBuyer godoc
// @Summary      set product buyer id
// @Description  Handler for setting product buyer id
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        id path int64 true "Product ID"
// @Param        body body models.SetProductBuyerRequest true "Body"
// @Success      200 {object} errors.Success
// @Failure      400  {object}  errors.Error
// @Failure      404  {object}  errors.Error
// @Failure      500  {object}  errors.Error
// @Router       /product/buyer/{id} [post]
func (ph *ProductHandler) SetProductBuyer(w http.ResponseWriter, r *http.Request) {
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

	req := &models.SetProductBuyerRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Error(err)
		errE := errors.UnexpectedBadRequest(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	errE := ph.productUcase.SetProductBuyer(productID, req.Buyer_id)
	if errE != nil {
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(errors.JSONSuccess("Successful set."))
}

// CreateProductReview godoc
// @Summary      Create product review
// @Description  Handler for creating product review
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        id path int64 true "Product ID"
// @Param        body body models.Review true "Review"
// @Success      200 {object} errors.Success
// @Failure      400  {object}  errors.Error
// @Failure      404  {object}  errors.Error
// @Failure      500  {object}  errors.Error
// @Router       /product/review/{id} [post]
func (ph *ProductHandler) CreateProductReview(w http.ResponseWriter, r *http.Request) {
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

	review := &models.Review{}
	err := json.NewDecoder(r.Body).Decode(&review)
	if err != nil {
		logger.Error(err)
		errE := errors.UnexpectedBadRequest(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
	review.ReviewerID = userID

	errE := ph.productUcase.CreateProductReview(review)
	if errE != nil {
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(errors.JSONSuccess("Successful creation."))
}

// GetUserReviews godoc
// @Summary      get user reviews
// @Description  Handler for getting user reviews
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        id path int64 true "User ID"
// @Param        type path string true "Type" Enums(seller, buyer)
// @Success      200 {object} []models.Review
// @Failure      400  {object}  errors.Error
// @Failure      404  {object}  errors.Error
// @Failure      500  {object}  errors.Error
// @Router       /user/{id}/reviews/{type} [get]
func (ph *ProductHandler) GetUserReviews(w http.ResponseWriter, r *http.Request) {
	logger, ok := r.Context().Value(middleware.ContextLogger).(*logrus.Entry)
	if !ok {
		logger = log.GetDefaultLogger()
		logger.Warn("no logger")
	}

	page := &models.PageWithSort{}
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

	vars := mux.Vars(r)
	userID, _ := strconv.ParseUint(vars["id"], 10, 64)
	reviewType := vars["type"]
	logger.Info("user id ", userID)

	reviews, errE := ph.productUcase.GetUserReviews(userID, reviewType, page)
	if errE != nil {
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(reviews)
	if err != nil {
		logger.Error(err)
		errE := errors.UnexpectedInternal(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
}

// GetWaitingReviews godoc
// @Summary      get waiting reviews
// @Description  Handler for getting waiting reviews
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        type path string true "Type" Enums(seller, buyer)
// @Param        From query integer false "From"
// @Param        Count query integer false "Count"
// @Success      200 {object} []models.WaitingReview
// @Failure      400  {object}  errors.Error
// @Failure      404  {object}  errors.Error
// @Failure      500  {object}  errors.Error
// @Router       /user/reviews/await/{type} [get]
func (ph *ProductHandler) GetWaitingReviews(w http.ResponseWriter, r *http.Request) {
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
	reviewType := vars["type"]
	
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

	reviews, errE := ph.productUcase.GetWaitingReviews(userID, reviewType, page)
	if errE != nil {
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(reviews)
	if err != nil {
		logger.Error(err)
		errE := errors.UnexpectedInternal(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
}
