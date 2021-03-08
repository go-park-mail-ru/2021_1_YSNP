package login

import (
	"2021_1_YSNP/models"
	_tmpDB "2021_1_YSNP/tmp_database"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func JSONError(message string) []byte {
	jsonError, err := json.Marshal(models.Error{Message: message})
	if err != nil {
		return []byte("")
	}
	return jsonError
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	signInData := models.LoginData{}
	err := json.NewDecoder(r.Body).Decode(&signInData)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(JSONError(err.Error()))
		return
	}
	fmt.Println(signInData)
	user, err := _tmpDB.GetUserByLogin(signInData.Telephone)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusNotFound)
		w.Write(JSONError(err.Error()))
		return
	}

	if user.Password != signInData.Password {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(JSONError(errors.New("Wrong password").Error()))
		return
	}

	cookie := http.Cookie{
		Name:     "session_id",
		Value:    _tmpDB.NewSession(user.Telephone),
		Expires:  time.Now().Add(10000 * time.Hour),
		Secure:   false,
		HttpOnly: false,
	}
	body, err := json.Marshal(map[string]string{"message": "Successful login."})
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(JSONError(err.Error()))
		return
	}

	_tmpDB.DeleteSession(session.Value)

	body, err := json.Marshal("logout success")
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
