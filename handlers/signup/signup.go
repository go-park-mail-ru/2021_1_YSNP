package signup

import (
	"2021_1_YSNP/models"
	_tmpDB "2021_1_YSNP/tmp_database"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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

	err = _tmpDB.NewUser(&signUpData)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(JSONError(err.Error()))
		return
	}

	fmt.Println("SignUpHandler", signUpData)

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    _tmpDB.NewSession(signUpData.Telephone),
		Expires:  time.Now().Add(10 * time.Hour),
		Secure:   false,
		HttpOnly: true,
	}

	body, err := json.Marshal(map[string]string{"message": "Successful registration."})
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

func UploadAvatarHandler(w http.ResponseWriter, r *http.Request) {
	authorized := false
	session, err := r.Cookie("session_id")
	if err == nil && session != nil {
		authorized = _tmpDB.CheckSession(session.Value)
	}

	if authorized {
		r.Body = http.MaxBytesReader(w, r.Body, 10*1024*1024)
		err := r.ParseMultipartForm(10 * 1024 * 1024)
		if err != nil {
			logrus.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(JSONError(err.Error()))
			return
		}

		fmt.Println("UploadAvatarHandler")

		file, _, err := r.FormFile("file-upload")
		if err != nil {
			logrus.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(JSONError(err.Error()))
			return
		}
		defer file.Close()

		r.FormValue("file-upload")

		str, err := os.Getwd()

		if err != nil {
			logrus.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(JSONError(err.Error()))
			return
		}

		photoPath := "static/avatar/"
		os.Chdir(photoPath)

		photoID, err := uuid.NewRandom()
		if err != nil {
			logrus.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(JSONError(err.Error()))
			return
		}

		f, err := os.OpenFile(photoID.String()+".jpg", os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			logrus.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(JSONError(err.Error()))
			return
		}
		defer f.Close()

		os.Chdir(str)

		_, err = io.Copy(f, file)
		if err != nil {
			logrus.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(JSONError(err.Error()))
			return
		}

		var avatar []string
		avatar = append(avatar, _tmpDB.Url+"/static/avatar/"+photoID.String()+".jpg")
		_tmpDB.SetUserAvatar(session.Value, avatar)

		body, err := json.Marshal(avatar)
		if err != nil {
			logrus.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(JSONError(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(body)
	} else {
		err = errors.New("user not authorised or not found")
		logrus.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(JSONError(err.Error()))
		return
	}
}
