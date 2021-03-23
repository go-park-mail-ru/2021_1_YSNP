package delivery

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/session"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/user"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
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

func (sh *SessionHandler) Configure(r *mux.Router) {
	r.HandleFunc("/login", sh.LoginHandler).Methods(http.MethodPost)
	r.HandleFunc("/logout", sh.LogoutHandler).Methods(http.MethodPost)
}

func (sh *SessionHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	login := &models.LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errors.JSONError(err.Error()))
		return
	}

	user, err := sh.userUcase.GetByTelephone(login.Telephone)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errors.JSONError(err.Error()))
		return
	}

	err = sh.userUcase.CheckPassword(user, login.Password)
	if  err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errors.JSONError(err.Error()))
		return
	}

	session := models.CreateSession(user.ID)
	err = sh.sessUcase.Create(session)
	if  err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errors.JSONError(err.Error()))
		return
	}

	cookie := http.Cookie{
		Name:     "session_id",
		Value:    session.Value,
		Expires:  session.ExpiresAt,
		Secure:   false,
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
	w.Write(errors.JSONSuccess("Successful login."))
}

func (sh *SessionHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errors.JSONError(err.Error()))
		return
	}

	err = sh.sessUcase.Delete(session.Value)
	if  err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errors.JSONError(err.Error()))
		return
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)

	w.WriteHeader(http.StatusOK)
	w.Write(errors.JSONSuccess("logout success"))
}