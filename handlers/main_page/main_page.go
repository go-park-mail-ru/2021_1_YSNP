package main_page

import (
	_tmpDB "2021_1_YSNP/tmp_database"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

func MainPageHandler(w http.ResponseWriter, r *http.Request){
	products := _tmpDB.GetProducts()

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	err := encoder.Encode(products)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}
}
