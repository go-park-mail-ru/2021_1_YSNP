package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/vk"

	"github.com/go-park-mail-ru/2021_1_YSNP/configs"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/session"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/errors"
	log "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/logger"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/tools/middleware"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/user"
)

type UserHandler struct {
	userUcase user.UserUsecase
	sessUcase session.SessionUsecase
}

func NewUserHandler(userUcase user.UserUsecase, sessUcase session.SessionUsecase) *UserHandler {
	return &UserHandler{
		userUcase: userUcase,
		sessUcase: sessUcase,
	}
}

func (uh *UserHandler) Configure(r *mux.Router, mw *middleware.Middleware) {
	r.HandleFunc("/signup", uh.SignUpHandler).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/upload", mw.CheckAuthMiddleware(uh.UploadAvatarHandler)).Methods(http.MethodPost, http.MethodOptions)

	r.HandleFunc("/me", mw.SetCSRFToken(mw.CheckAuthMiddleware(uh.GetProfileHandler))).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/user/{id:[0-9]+}", mw.SetCSRFToken(mw.CheckAuthMiddleware(uh.GetSellerHandler))).Methods(http.MethodGet, http.MethodOptions)

	r.HandleFunc("/user", mw.CheckAuthMiddleware(uh.ChangeProfileHandler)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/user/password", mw.CheckAuthMiddleware(uh.ChangeProfilePasswordHandler)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/user/position", mw.CheckAuthMiddleware(uh.ChangeUserLocationHandler)).Methods(http.MethodPost, http.MethodOptions)

	r.HandleFunc("/oauth/vk", uh.VKOauth).Methods(http.MethodOptions, http.MethodGet)
}

func (uh *UserHandler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	logger, ok := r.Context().Value(middleware.ContextLogger).(*logrus.Entry)
	if !ok {
		logger = log.GetDefaultLogger()
		logger.Warn("no logger")
	}
	defer r.Body.Close()

	signUp := models.SignUpRequest{}
	err := json.NewDecoder(r.Body).Decode(&signUp)
	if err != nil {
		logger.Error(err)
		errE := errors.UnexpectedBadRequest(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
	logger.Info("user data ", signUp)

	sanitizer := bluemonday.UGCPolicy()
	signUp.Name = sanitizer.Sanitize(signUp.Name)
	signUp.Surname = sanitizer.Sanitize(signUp.Surname)
	signUp.Sex = sanitizer.Sanitize(signUp.Sex)
	signUp.Email = sanitizer.Sanitize(signUp.Email)
	signUp.Telephone = sanitizer.Sanitize(signUp.Telephone)
	signUp.Password1 = sanitizer.Sanitize(signUp.Password1)
	signUp.Password2 = sanitizer.Sanitize(signUp.Password2)
	signUp.DateBirth = sanitizer.Sanitize(signUp.DateBirth)
	logger.Debug("sanitize user data ", signUp)

	_, err = govalidator.ValidateStruct(signUp)
	if err != nil {
		if allErrs, ok := err.(govalidator.Errors); ok {
			logger.Error(allErrs.Errors())
			errE := errors.UnexpectedBadRequest(allErrs)
			w.WriteHeader(errE.HttpError)
			w.Write(errors.JSONError(errE))
			return
		}
	}

	//TODO(Maxim) мне кажется это нужно делать в usecase
	user := &models.UserData{
		Name:       signUp.Name,
		Surname:    signUp.Surname,
		Sex:        signUp.Sex,
		Email:      signUp.Email,
		Telephone:  signUp.Telephone,
		Password:   signUp.Password1,
		DateBirth:  signUp.DateBirth,
		LinkImages: "/static/avatar/profile.webp",
	}

	errE := uh.userUcase.Create(user)
	if errE != nil {
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
	logger.Debug("user ", user)

	session := models.CreateSession(user.ID)
	errE = uh.sessUcase.Create(session)
	if errE != nil {
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
	logger.Debug("session ", session)

	cookie := http.Cookie{
		Name:    "session_id",
		Value:   session.Value,
		Expires: session.ExpiresAt,
		Path:    "/",
		//Secure:   true,
		//SameSite: http.SameSiteLaxMode,
		//HttpOnly: true,
	}
	logger.Debug("cookie ", cookie)

	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
	w.Write(errors.JSONSuccess("Successful login."))
}

func (uh *UserHandler) UploadAvatarHandler(w http.ResponseWriter, r *http.Request) {
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
	logger.Debug("user id ", userID)

	r.Body = http.MaxBytesReader(w, r.Body, 3*1024*1024)
	err := r.ParseMultipartForm(3 * 1024 * 1024)
	if err != nil {
		logger.Error(err)
		errE := errors.UnexpectedBadRequest(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	if len(r.MultipartForm.File["file-upload"]) != 0 {
		file := r.MultipartForm.File["file-upload"][0]
		_, errE := uh.userUcase.UpdateAvatar(userID, file)
		if errE != nil {
			logger.Error(errE.Message)
			w.WriteHeader(errE.HttpError)
			w.Write(errors.JSONError(errE))
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write(errors.JSONSuccess("Successful upload."))
}

func (uh *UserHandler) GetProfileHandler(w http.ResponseWriter, r *http.Request) {
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

	user, errE := uh.userUcase.GetByID(userID)
	if errE != nil {
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
	logger.Debug("user ", user)

	body, err := json.Marshal(user)
	if err != nil {
		logger.Error(err)
		errE := errors.UnexpectedInternal(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (uh *UserHandler) GetSellerHandler(w http.ResponseWriter, r *http.Request) {
	logger, ok := r.Context().Value(middleware.ContextLogger).(*logrus.Entry)
	if !ok {
		logger = log.GetDefaultLogger()
		logger.Warn("no logger")
	}

	vars := mux.Vars(r)
	sellerID, _ := strconv.ParseUint(vars["id"], 10, 64)
	logger.Info("seller id ", sellerID)

	userID, ok := r.Context().Value(middleware.ContextUserID).(uint64)
	if !ok {
		errE := errors.Cause(errors.UserUnauthorized)
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
	logger.Info("user id ", userID)

	user, errE := uh.userUcase.GetSellerByID(sellerID)
	if errE != nil {
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
	logger.Debug("seller ", user)

	body, err := json.Marshal(user)
	if err != nil {
		logger.Error(err)
		errE := errors.UnexpectedInternal(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (uh *UserHandler) ChangeProfileHandler(w http.ResponseWriter, r *http.Request) {
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

	changeData := models.SignUpRequest{}
	err := json.NewDecoder(r.Body).Decode(&changeData)
	if err != nil {
		logger.Error(err)
		errE := errors.UnexpectedBadRequest(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
	logger.Info("user data ", changeData)

	sanitizer := bluemonday.UGCPolicy()
	changeData.Name = sanitizer.Sanitize(changeData.Name)
	changeData.Surname = sanitizer.Sanitize(changeData.Surname)
	changeData.Sex = sanitizer.Sanitize(changeData.Sex)
	changeData.Email = sanitizer.Sanitize(changeData.Email)
	changeData.Telephone = sanitizer.Sanitize(changeData.Telephone)
	changeData.DateBirth = sanitizer.Sanitize(changeData.DateBirth)
	logger.Debug("sanitize user data ", changeData)

	_, err = govalidator.ValidateStruct(changeData)
	if err != nil {
		if allErrs, ok := err.(govalidator.Errors); ok {
			logger.Error(allErrs.Errors())
			errE := errors.UnexpectedBadRequest(allErrs)
			w.WriteHeader(errE.HttpError)
			w.Write(errors.JSONError(errE))
			return
		}
	}

	//TODO(Maxim) мне кажется это нужно делать в usecase
	user := &models.UserData{
		Name:      changeData.Name,
		Surname:   changeData.Surname,
		Sex:       changeData.Sex,
		Email:     changeData.Email,
		Telephone: changeData.Telephone,
		DateBirth: changeData.DateBirth,
	}
	logger.Debug("user ", user)

	_, errE := uh.userUcase.UpdateProfile(userID, user)
	if errE != nil {
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(errors.JSONSuccess("Successful change."))
}

func (uh *UserHandler) ChangeProfilePasswordHandler(w http.ResponseWriter, r *http.Request) {
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

	passwordData := models.PasswordChangeRequest{}
	err := json.NewDecoder(r.Body).Decode(&passwordData)
	if err != nil {
		logger.Error(err)
		errE := errors.UnexpectedBadRequest(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
	logger.Info("user password ", passwordData)

	sanitizer := bluemonday.UGCPolicy()
	passwordData.OldPassword = sanitizer.Sanitize(passwordData.OldPassword)
	passwordData.NewPassword1 = sanitizer.Sanitize(passwordData.NewPassword1)
	passwordData.NewPassword2 = sanitizer.Sanitize(passwordData.NewPassword2)
	logger.Debug("sanitize user data ", passwordData)

	_, err = govalidator.ValidateStruct(passwordData)
	if err != nil {
		if allErrs, ok := err.(govalidator.Errors); ok {
			logger.Error(allErrs.Errors())
			errE := errors.UnexpectedBadRequest(allErrs)
			w.WriteHeader(errE.HttpError)
			w.Write(errors.JSONError(errE))
			return
		}
	}

	_, errE := uh.userUcase.UpdatePassword(userID, passwordData.NewPassword1)
	if errE != nil {
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(errors.JSONSuccess("Successful change."))
}

func (uh *UserHandler) ChangeUserLocationHandler(w http.ResponseWriter, r *http.Request) {
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

	positionData := models.LocationRequest{}
	err := json.NewDecoder(r.Body).Decode(&positionData)
	if err != nil {
		logger.Error(err)
		errE := errors.UnexpectedBadRequest(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
	logger.Info("user position ", positionData)

	sanitizer := bluemonday.UGCPolicy()
	positionData.Address = sanitizer.Sanitize(positionData.Address)
	logger.Debug("sanitize user position ", positionData)

	_, err = govalidator.ValidateStruct(positionData)
	if err != nil {
		if allErrs, ok := err.(govalidator.Errors); ok {
			logger.Error(allErrs.Errors())
			errE := errors.UnexpectedBadRequest(allErrs)
			w.WriteHeader(errE.HttpError)
			w.Write(errors.JSONError(errE))
			return
		}
	}

	_, errE := uh.userUcase.UpdateLocation(userID, &positionData)
	if errE != nil {
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(errors.JSONSuccess("Successful change."))
}

func (uh *UserHandler) VKOauth(w http.ResponseWriter, r *http.Request) {
	logger, ok := r.Context().Value(middleware.ContextLogger).(*logrus.Entry)
	if !ok {
		logger = log.GetDefaultLogger()
		logger.Warn("no logger")
	}
	defer r.Body.Close()

	ctx := r.Context()
	code := r.FormValue("code")
	conf := &oauth2.Config{
		ClientID:     configs.Configs.GetVKAppID(),
		ClientSecret: configs.Configs.GetVKAppKey(),
		RedirectURL:  configs.Configs.GetVKRedirectUrl(),
		Endpoint:     vk.Endpoint,
	}
	logger.Debug("code ", code)

	token, err := conf.Exchange(ctx, code)
	if err != nil {
		logger.Error(err)
		errE := errors.UnexpectedBadRequest(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		http.Redirect(w, r, configs.Configs.GetFrontendUrl(), http.StatusTemporaryRedirect)
		return
	}
	logger.Debug("token ", token)

	userID := token.Extra("user_id")
	logger.Info("userID ", userID)

	client := conf.Client(ctx, token)
	resp, err := client.Get(fmt.Sprintf(configs.Configs.GetVKAppUrl(), token.AccessToken))
	if err != nil {
		logger.Error(err)
		errE := errors.UnexpectedBadRequest(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		http.Redirect(w, r, configs.Configs.GetFrontendUrl(), http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	data := &models.Response{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		logger.Error(err)
		errE := errors.UnexpectedBadRequest(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		http.Redirect(w, r, configs.Configs.GetFrontendUrl(), http.StatusTemporaryRedirect)
		return
	}
	logger.Info("oauth data ", data)

	userOAuth := &models.UserOAuthRequest{
		LastName:      data.Response[0].LastName,
		FirstName:     data.Response[0].FirstName,
		Photo:         data.Response[0].Photo,
		UserOAuthID:   userID.(float64),
		UserOAuthType: "vk",
	}
	logger.Debug("user oauth data ", userOAuth)

	errE := uh.userUcase.CreateOrLogin(userOAuth)
	if errE != nil {
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		http.Redirect(w, r, configs.Configs.GetFrontendUrl(), http.StatusTemporaryRedirect)
		return
	}
	logger.Debug("userID ", userOAuth.ID)

	session := models.CreateSession(userOAuth.ID)
	errE = uh.sessUcase.Create(session)
	if errE != nil {
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		http.Redirect(w, r, configs.Configs.GetFrontendUrl(), http.StatusTemporaryRedirect)
		return
	}
	logger.Debug("session ", session)

	cookie := http.Cookie{
		Name:    "session_id",
		Value:   session.Value,
		Expires: session.ExpiresAt,
		Path:    "/",
		//Secure:   true,
		//SameSite: http.SameSiteLaxMode,
		//HttpOnly: true,
	}
	logger.Debug("cookie ", cookie)

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, configs.Configs.GetFrontendUrl(), http.StatusTemporaryRedirect)
}
