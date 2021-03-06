package MainPage

import (
	_tmpDB "2021_1_YSNP/tmp_database"
	"encoding/json"
	"log"
	"net/http"
)

func MainPageHandler(w http.ResponseWriter, r *http.Request){
	products := _tmpDB.GetProducts()

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	err := encoder.Encode(products)
	if err != nil {
		log.Printf("error while marshalling JSON: %s", err)
		w.Write([]byte("{}"))
		return
	}
}
