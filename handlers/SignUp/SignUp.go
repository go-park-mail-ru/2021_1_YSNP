package SignUp

import (
	"2021_1_YSNP/models"
	"encoding/json"
	"net/http"
)

func SignUpHandler(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
	signUpData := models.SignUpData{}
	err := json.NewDecoder(r.Body).Decode(&signUpData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	body, err := json.Marshal(signUpData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
