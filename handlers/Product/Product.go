package Product

import (
	"2021_1_YSNP/models"
	_tmpDB "2021_1_YSNP/tmp_database"
	"encoding/json"
	"net/http"
	"strings"
)

func ProductIDHandler(w http.ResponseWriter, r *http.Request){
	productID := strings.TrimPrefix(r.URL.Path,"/product/")

	product, err := _tmpDB.GetProduct(productID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	body, err := json.Marshal(product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func ProductCreateHandler (w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
	productData := models.ProductData{}
	parseErr := json.NewDecoder(r.Body).Decode(&productData)
	if parseErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(parseErr.Error()))
		return
	}

	insertErr := _tmpDB.NewProduct(&productData)
	if insertErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(insertErr.Error()))
		return
	}

	body, err := json.Marshal(productData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

