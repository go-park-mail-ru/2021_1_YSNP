package product

import (
	"2021_1_YSNP/models"
	_tmpDB "2021_1_YSNP/tmp_database"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func JSONError(message string) []byte {
	jsonError, err := json.Marshal(models.Error{Message: message})
	if err != nil {
		return []byte("")
	}
	return jsonError
}

func ProductIDHandler(w http.ResponseWriter, r *http.Request) {
	productID := strings.TrimPrefix(r.URL.Path, "/api/v1/product/")

	product, err := _tmpDB.GetProduct(productID)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusNotFound)
		w.Write(JSONError(err.Error()))
		return
	}

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
	defer r.Body.Close()
	productData := models.ProductData{}
	err := json.NewDecoder(r.Body).Decode(&productData)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(JSONError(err.Error()))
		return
	}

	err = _tmpDB.NewProduct(&productData)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}

	body, err := json.Marshal(productData)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
