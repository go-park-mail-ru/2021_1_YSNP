package product

import (
	"2021_1_YSNP/models"
	_tmpDB "2021_1_YSNP/tmp_database"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
)

func JSONError(message string) []byte {
	jsonError, err := json.Marshal(models.Error{Message: message})
	if err != nil {
		return []byte("")
	}

	return jsonError
}

func ProductIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]
	product, err := _tmpDB.GetProduct(productID)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusNotFound)
		w.Write(JSONError(err.Error()))
		return
	}

	fmt.Println("ProductIDHandler", productID)

	body, err := json.Marshal(product)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func ProductCreateHandler(w http.ResponseWriter, r *http.Request) {
	authorized := false
	session, err := r.Cookie("session_id")
	if err == nil && session != nil {
		authorized = _tmpDB.CheckSession(session.Value)
	}

	if authorized {
		defer r.Body.Close()
		productData := models.ProductData{}
		err := json.NewDecoder(r.Body).Decode(&productData)
		if err != nil {
			logrus.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(JSONError(err.Error()))
			return
		}

		fmt.Println("ProductCreateHandler", session.Value)

		err = _tmpDB.NewProduct(&productData, session.Value)
		if err != nil {
			logrus.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(JSONError(err.Error()))
			return
		}

		body, err := json.Marshal(map[string]string{"message": "Successful creation."})
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

func UploadPhotoHandler(w http.ResponseWriter, r *http.Request) {
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

		fmt.Println("UploadPhotoHandler")

		files := r.MultipartForm.File["photos"]
		imgs := make(map[string][]string)
		for i, _ := range files {
			file, err := files[i].Open()
			defer file.Close()
			if err != nil {
				logrus.Error(err)
				w.WriteHeader(http.StatusBadRequest)
				w.Write(JSONError(err.Error()))
				return
			}

			str, err := os.Getwd()
			if err != nil {
				logrus.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(JSONError(err.Error()))
				return
			}

			photoPath := "static/product/"
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

			imgs["linkImages"] = append(imgs["linkImages"], _tmpDB.Url+"/static/product/"+photoID.String()+".jpg")
		}

		if len(imgs) == 0 {
			logrus.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(JSONError(errors.New("http: no such file").Error()))
			return
		}
		body, err := json.Marshal(imgs)
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
