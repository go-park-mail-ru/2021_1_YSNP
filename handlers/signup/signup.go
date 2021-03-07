package SignUp

import (
	"2021_1_YSNP/models"
	_tmpDB "2021_1_YSNP/tmp_database"
	"encoding/json"
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

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	signUpData := models.SignUpData{}
	err := json.NewDecoder(r.Body).Decode(&signUpData)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(JSONError(err.Error()))
		return
	}
	fmt.Println(signUpData)
	err = _tmpDB.NewUser(&signUpData)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    _tmpDB.NewSession(signUpData.Telephone),
		Expires:  time.Now().Add(10 * time.Hour),
		Secure:   false,
		HttpOnly: false,
	}

	body, err := json.Marshal(signUpData)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
