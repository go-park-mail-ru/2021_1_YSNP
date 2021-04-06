package http

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/favorite"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/middleware"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type FavoriteHandler struct {
	favoriteUcase favorite.FavoriteUsecase
}

func NewFavoriteHandler(favoriteUcase favorite.FavoriteUsecase) *FavoriteHandler {
	return &FavoriteHandler{
		favoriteUcase: favoriteUcase,
	}
}

func (ph *FavoriteHandler) Configure(r *mux.Router, mw *middleware.Middleware) {
	r.HandleFunc("/user/favorite/like/{id:[0-9]+}", mw.CheckAuthMiddleware(ph.LikeProductHandler)).Methods(http.MethodPost)
	r.HandleFunc("/user/favorite/dislike/{id:[0-9]+}", mw.CheckAuthMiddleware(ph.DislikeProductHandler)).Methods(http.MethodPost)
	r.HandleFunc("/user/favorite/list", mw.CheckAuthMiddleware(ph.UserFavoriteHandler)).Methods(http.MethodGet)
}

func (ph *FavoriteHandler) LikeProductHandler(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value(middleware.ContextLogger).(*logrus.Entry)

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

	errE := ph.favoriteUcase.LikeProduct(userID, productID)
	if errE != nil {
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(errors.JSONSuccess("Successful like."))
}

func (ph *FavoriteHandler) DislikeProductHandler(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value(middleware.ContextLogger).(*logrus.Entry)

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

	errE := ph.favoriteUcase.DislikeProduct(userID, productID)
	if errE != nil {
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(errors.JSONSuccess("Successful dislike."))
}

func (ph *FavoriteHandler) UserFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value(middleware.ContextLogger).(*logrus.Entry)

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

	_, err = govalidator.ValidateStruct(page)
	if err != nil {
		if allErrs, ok := err.(govalidator.Errors); ok {
			logger.Error(allErrs.Errors())
			errE := errors.UnexpectedBadRequest(allErrs)
			w.WriteHeader(errE.HttpError)
			w.Write(errors.JSONError(errE))
			return
		}
	}

	products, errE := ph.favoriteUcase.GetUserFavorite(userID, page)
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
