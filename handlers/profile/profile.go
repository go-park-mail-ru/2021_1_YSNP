package profile

import (
	"2021_1_YSNP/models"
	_tmpDB "2021_1_YSNP/tmp_database"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

func JSONError(message string) []byte {
	jsonError, err := json.Marshal(models.Error{Message: message})
	if err != nil {
		return []byte("")
	}

	return jsonError
}

func GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	authorized := false
	session, err := r.Cookie("session_id")
	if err == nil && session != nil {
		authorized = _tmpDB.CheckSession(session.Value)
	}

	if !authorized {
		err = errors.New("user not authorised or not found")
		logrus.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(JSONError(err.Error()))
		return
	}

	userInfo := _tmpDB.GetUserBySession(session.Value)
	userInfo.Password = ""

	fmt.Println("GetProfileHandler", session.Value)

	body, err := json.Marshal(userInfo)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func ChangeProfileHandler(w http.ResponseWriter, r *http.Request) {
	authorized := false
	session, err := r.Cookie("session_id")
	if err == nil && session != nil {
		authorized = _tmpDB.CheckSession(session.Value)
	}

	if !authorized {
		err = errors.New("user not authorised or not found")
		logrus.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(JSONError(err.Error()))
		return
	}

	signUpData := models.SignUpData{}
	err = json.NewDecoder(r.Body).Decode(&signUpData)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(JSONError(err.Error()))
		return
	}

	fmt.Println("ChangeProfileHandler", session.Value)

	err = _tmpDB.ChangeUserData(session.Value, &signUpData)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}

	body, err := json.Marshal(map[string]string{"message": "Successful change."})
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func ChangeProfilePasswordHandler(w http.ResponseWriter, r *http.Request) {
	authorized := false
	session, err := r.Cookie("session_id")
	if err == nil && session != nil {
		authorized = _tmpDB.CheckSession(session.Value)
	}

	if !authorized {
		err = errors.New("user not authorised or not found")
		logrus.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(JSONError(err.Error()))
		return
	}

	passwordData := models.PasswordChange{}
	err = json.NewDecoder(r.Body).Decode(&passwordData)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(JSONError(err.Error()))
		return
	}

	fmt.Println("ChangeProfilePasswordHandler", session.Value)

	err = _tmpDB.ChangeUserPassword(session.Value, &passwordData)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(JSONError(err.Error()))
		return
	}

	body, err := json.Marshal(map[string]string{"message": "Successful change."})
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
