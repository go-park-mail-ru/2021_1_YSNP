package SignUp

import (
	"2021_1_YSNP/models"
	_tmpDB "2021_1_YSNP/tmp_database"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	signUpData := models.SignUpData{}
	err := json.NewDecoder(r.Body).Decode(&signUpData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	fmt.Println(signUpData)
	insertErr := _tmpDB.NewUser(&signUpData)
	if insertErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
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
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
