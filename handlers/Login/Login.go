package Login

import (
	"2021_1_YSNP/models"
	_tmpDB "2021_1_YSNP/tmp_database"
	"encoding/json"
	"net/http"
	"time"
)



func LoginHandler(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
	signInData := models.LoginData{}
	err := json.NewDecoder(r.Body).Decode(&signInData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	user, err := _tmpDB.GetUserByLogin(signInData.Telephone)
	if err != nil {
		http.Error(w, `no user`, 404)
		return
	}

	if user.Password != signInData.Password {
		http.Error(w, `bad pass`, 400)
		return
	}

	cookie := &http.Cookie{
		Name:       "session_id",
		Value:      _tmpDB.NewSession(user.Telephone),
		Expires:    time.Now().Add(10 * time.Hour),
		Secure:     false,
		HttpOnly:   false,
	}
	body, err := json.Marshal(signInData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
