package Profile

import (
	"2021_1_YSNP/models"
	_tmpDB "2021_1_YSNP/tmp_database"
	"encoding/json"
	"errors"
	"net/http"
)

func GetProfileHandler(w http.ResponseWriter, r *http.Request){
	authorized := false
	session, err := r.Cookie("session_id")
	if err == nil && session != nil {
		authorized = _tmpDB.CheckSession(session.Value)
	}

	if authorized {
		userInfo := _tmpDB.GetUserBySession(session.Value)

		body, marshErr := json.Marshal(userInfo)
		if marshErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(marshErr.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(body)

	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(errors.New("User not authorised or not found").Error()))
		return
	}
}

func ChangeProfileHandler(w http.ResponseWriter, r *http.Request){
	authorized := false
	session, err := r.Cookie("session_id")
	if err == nil && session != nil {
		authorized = _tmpDB.CheckSession(session.Value)
	}

	if authorized {
		signUpData := models.SignUpData{}
		decodeErr := json.NewDecoder(r.Body).Decode(&signUpData)
		if decodeErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(decodeErr.Error()))
			return
		}

		insertErr := _tmpDB.ChangeUserData(session.Value, &signUpData)
		if insertErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(insertErr.Error()))
			return
		}

		body, marshErr := json.Marshal(signUpData)
		if marshErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(marshErr.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(body)

	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(errors.New("User not authorised or not found").Error()))
		return
	}
}
