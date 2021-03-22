package http

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/product"
	"github.com/gorilla/mux"
	"net/http"
)

type ProductHandler struct {
	productUcase product.ProductUsecase
}

func NewProductHandler(productUcase product.ProductUsecase) *ProductHandler{
	return &ProductHandler{
		productUcase: productUcase,
	}
}

func (ph *ProductHandler) Configure(r *mux.Router) {
	r.HandleFunc("/product/list", ph.MainPageHandler).Methods(http.MethodGet)
	r.HandleFunc("/product/{id:[0-9]+}", ph.ProductIDHandler).Methods(http.MethodGet)
	r.HandleFunc("/product/create", ph.ProductCreateHandler).Methods(http.MethodPost)
	r.HandleFunc("/product/upload", ph.UploadPhotoHandler).Methods(http.MethodPost)
}

func (ph *ProductHandler) MainPageHandler(w http.ResponseWriter, r *http.Request) {

}

func (ph *ProductHandler) ProductIDHandler(w http.ResponseWriter, r *http.Request) {

}

func (ph *ProductHandler) ProductCreateHandler(w http.ResponseWriter, r *http.Request) {

}

func (ph *ProductHandler) UploadPhotoHandler(w http.ResponseWriter, r *http.Request) {

}

