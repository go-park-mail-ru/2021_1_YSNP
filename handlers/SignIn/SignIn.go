package SignIn

import (
	"2021_1_YSNP/models"
	"encoding/json"
	"net/http"
)

func SignInHandler(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
	signInData := models.SignInData{}
	err := json.NewDecoder(r.Body).Decode(&signInData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	body, err := json.Marshal(signInData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
