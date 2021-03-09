package signup

import (
	"2021_1_YSNP/models"
	_tmpDB "2021_1_YSNP/tmp_database"
	"encoding/json"
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
	fmt.Println(signUpData)
	err = _tmpDB.NewUser(&signUpData)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
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
	r.Body = http.MaxBytesReader(w, r.Body, 10*1024*1024)
	err := r.ParseMultipartForm(10 * 1024 * 1024)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(JSONError(err.Error()))
		return
	}

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
	fmt.Println(str)
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

	mymap := make(map[string]string)

	mymap["linkImages"] = "http://89.208.199.170:8080/static/avatar/" + photoID.String() + ".jpg"

	body, err := json.Marshal(mymap)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
