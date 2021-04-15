package delivery

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	"github.com/sirupsen/logrus"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/errors"
	log "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/logger"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/middleware"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/session"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/user"
)

type SessionHandler struct {
	sessUcase session.SessionUsecase
	userUcase user.UserUsecase
}

func NewSessionHandler(sessUcase session.SessionUsecase, userUcase user.UserUsecase) *SessionHandler {
	return &SessionHandler{
		sessUcase: sessUcase,
		userUcase: userUcase,
	}
}

func (sh *SessionHandler) Configure(r *mux.Router, mw *middleware.Middleware) {
	r.HandleFunc("/login", sh.LoginHandler).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/logout", mw.CheckAuthMiddleware(sh.LogoutHandler)).Methods(http.MethodPost, http.MethodOptions)
}

func (sh *SessionHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	logger, ok := r.Context().Value(middleware.ContextLogger).(*logrus.Entry)
	if !ok {
		logger = log.GetDefaultLogger()
		logger.Warn("no logger")
	}
	defer r.Body.Close()

	login := &models.LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		logger.Error(err)
		errE := errors.UnexpectedBadRequest(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
	logger.Info("user data ", login)

	sanitizer := bluemonday.UGCPolicy()
	login.Telephone = sanitizer.Sanitize(login.Telephone)
	login.Password = sanitizer.Sanitize(login.Password)
	logger.Debug("sanitize user data ", login)

	_, err = govalidator.ValidateStruct(login)
	if err != nil {
		if allErrs, ok := err.(govalidator.Errors); ok {
			logger.Error(allErrs.Errors())
			errE := errors.UnexpectedBadRequest(allErrs)
			w.WriteHeader(errE.HttpError)
			w.Write(errors.JSONError(errE))
			return
		}
	}

	user, errE := sh.userUcase.GetByTelephone(login.Telephone)
	if errE != nil {
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
	logger.Debug("user ", user)

	errE = sh.userUcase.CheckPassword(user, login.Password)
	if errE != nil {
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	session := models.CreateSession(user.ID)
	errE = sh.sessUcase.Create(session)
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

func (sh *SessionHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	logger, ok := r.Context().Value(middleware.ContextLogger).(*logrus.Entry)
	if !ok {
		logger = log.GetDefaultLogger()
		logger.Warn("no logger")
	}

	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		errE := errors.Cause(errors.UserUnauthorized)
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
	logger.Info("session ", session)

	errE := sh.sessUcase.Delete(session.Value)
	if errE != nil {
		logger.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)

	w.WriteHeader(http.StatusOK)
	w.Write(errors.JSONSuccess("Successful logout."))
}
