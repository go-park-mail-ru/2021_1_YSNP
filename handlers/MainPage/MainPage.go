package MainPage

import (
	"2021_1_YSNP/models"
	"encoding/json"
	"log"
	"net/http"
)

func MainPageHandler(w http.ResponseWriter, r *http.Request){
	productData := models.ProductData{
		ID:          0,
		Name:        "iphone",
		Date:        2000,
		Amount:      12000,
		LinkImages:  nil,
		Description: "eto iphone",
		OwnerID:     0,
		Views:       1000,
		Likes:       20,
	}

	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)

	err := encoder.Encode(productData)

	if err != nil {
		log.Printf("error while marshalling JSON: %s", err)
		w.Write([]byte("{}"))
		return
	}


}
