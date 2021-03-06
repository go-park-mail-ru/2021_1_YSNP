package Login

import (
	"2021_1_YSNP/models"
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
	
	cookie := &http.Cookie{
		Name:       "session_id",
		Value:      "",
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
