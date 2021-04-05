package delivery

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/middleware"
	"github.com/microcosm-cc/bluemonday"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/session"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/user"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
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
	r.HandleFunc("/signup", uh.SignUpHandler).Methods(http.MethodPost)
	r.HandleFunc("/upload", mw.CheckAuthMiddleware(uh.UploadAvatarHandler)).Methods(http.MethodPost)
	r.HandleFunc("/me", mw.CheckAuthMiddleware(uh.GetProfileHandler)).Methods(http.MethodGet)
	r.HandleFunc("/settings", mw.CheckAuthMiddleware(uh.ChangeProfileHandler)).Methods(http.MethodPost)
	r.HandleFunc("/settings/password", mw.CheckAuthMiddleware(uh.ChangeProfilePasswordHandler)).Methods(http.MethodPost)
}

func (uh *UserHandler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value(middleware.ContextLogger).(*logrus.Entry)

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

	user := &models.UserData{
		Name:       signUp.Name,
		Surname:    signUp.Surname,
		Sex:        signUp.Sex,
		Email:      signUp.Email,
		Telephone:  signUp.Telephone,
		Password:   signUp.Password1,
		DateBirth:  signUp.DateBirth,
		LinkImages: signUp.LinkImages,
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
		Name:     "session_id",
		Value:    session.Value,
		Expires:  session.ExpiresAt,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
	}
	logger.Debug("cookie ", cookie)

	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
	w.Write(errors.JSONSuccess("Successful login."))
}

func (uh *UserHandler) UploadAvatarHandler(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value(middleware.ContextLogger).(*logrus.Entry)

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

	file, handler, err := r.FormFile("file-upload")
	if err != nil {
		logger.Error(err)
		errE := errors.UnexpectedBadRequest(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
	logger.Debug("photo ", handler.Header)
	defer file.Close()
	extension := filepath.Ext(handler.Filename)

	r.FormValue("file-upload")

	str, err := os.Getwd()
	if err != nil {
		logger.Error(err)
		errE := errors.UnexpectedInternal(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	photoPath := "static/avatar/"
	os.Chdir(photoPath)

	photoID, err := uuid.NewRandom()
	if err != nil {
		logger.Error(err)
		errE := errors.UnexpectedInternal(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
	logger.Debug("new photo name ", photoID)

	f, err := os.OpenFile(photoID.String()+extension, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Error(err)
		errE := errors.UnexpectedInternal(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
	defer f.Close()

	os.Chdir(str)

	_, err = io.Copy(f, file)
	if err != nil {
		_ = os.Remove(photoID.String() + extension)
		logger.Error(err)
		errE := errors.UnexpectedInternal(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	avatar := "/static/avatar/" + photoID.String() + extension

	_, errE := uh.userUcase.UpdateAvatar(userID, avatar)
	if errE != nil {
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	body, err := json.Marshal(avatar)
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

func (uh *UserHandler) GetProfileHandler(w http.ResponseWriter, r *http.Request) {
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

	_, err := govalidator.ValidateStruct(userID)
	if err != nil {
		if allErrs, ok := err.(govalidator.Errors); ok {
			logger.Error(allErrs.Errors())
			errE := errors.UnexpectedBadRequest(allErrs)
			w.WriteHeader(errE.HttpError)
			w.Write(errors.JSONError(errE))
			return
		}
	}

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

func (uh *UserHandler) ChangeProfileHandler(w http.ResponseWriter, r *http.Request) {
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
	passwordData.NewPassword = sanitizer.Sanitize(passwordData.NewPassword)
	passwordData.OldPassword = sanitizer.Sanitize(passwordData.OldPassword)
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
