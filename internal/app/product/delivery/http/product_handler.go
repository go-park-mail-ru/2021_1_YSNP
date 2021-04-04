package http

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/middleware"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/errors"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/product"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type ProductHandler struct {
	productUcase product.ProductUsecase
}

func NewProductHandler(productUcase product.ProductUsecase) *ProductHandler {
	return &ProductHandler{
		productUcase: productUcase,
	}
}

func (ph *ProductHandler) Configure(r *mux.Router, mw *middleware.Middleware) {
	r.HandleFunc("/product/list", ph.MainPageHandler).Methods(http.MethodPost)
	r.HandleFunc("/product/{id:[0-9]+}", ph.ProductIDHandler).Methods(http.MethodGet)
	r.HandleFunc("/product/create", mw.CheckAuthMiddleware(ph.ProductCreateHandler)).Methods(http.MethodPost)
	r.HandleFunc("/product/upload/{pid:[0-9]+}", mw.CheckAuthMiddleware(ph.UploadPhotoHandler)).Methods(http.MethodPost)
}

func (ph *ProductHandler) MainPageHandler(w http.ResponseWriter, r *http.Request) {
	page := models.Page{}
	err := json.NewDecoder(r.Body).Decode(&page)
	if err != nil {
		logrus.Error(err)
		errE := errors.UnexpectedBadRequest(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	products, errE := ph.productUcase.ListLatest(&page.Content)
	if errE != nil {
		logrus.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(products)
	if err != nil {
		logrus.Error(err)
		errE := errors.UnexpectedInternal(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}
}

func (ph *ProductHandler) ProductIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID, _ := strconv.ParseUint(vars["id"], 10, 64)

	product, errE := ph.productUcase.GetByID(productID)
	if errE != nil {
		logrus.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	body, err := json.Marshal(product)
	if err != nil {
		logrus.Error(err)
		errE := errors.UnexpectedInternal(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (ph *ProductHandler) ProductCreateHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userID, ok := r.Context().Value(middleware.ContextUserID).(uint64)
	if !ok {
		errE := errors.Cause(errors.UserUnauthorized)
		logrus.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	productData := &models.ProductData{}
	err := json.NewDecoder(r.Body).Decode(&productData)
	if err != nil {
		logrus.Error(err)
		errE := errors.UnexpectedBadRequest(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	productData.OwnerID = userID

	errE := ph.productUcase.Create(productData)
	if errE != nil {
		logrus.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(errors.JSONSuccess("Successful creation.", productData.ID))
}

func (ph *ProductHandler) UploadPhotoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID, _ := strconv.ParseUint(vars["pid"], 10, 64)

	_, ok := r.Context().Value(middleware.ContextUserID).(uint64)
	if !ok {
		errE := errors.Cause(errors.UserUnauthorized)
		logrus.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 10*1024*1024)
	err := r.ParseMultipartForm(10 * 1024 * 1024)
	if err != nil {
		logrus.Error(err)
		errE := errors.UnexpectedBadRequest(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	files := r.MultipartForm.File["photos"]
	imgs := make(map[string][]string)
	for i := range files {
		file, err := files[i].Open()
		if err != nil {
			logrus.Error(err)
			errE := errors.UnexpectedBadRequest(err)
			w.WriteHeader(errE.HttpError)
			w.Write(errors.JSONError(errE))
			return
		}
		defer file.Close()
		extension := filepath.Ext(files[i].Filename)

		str, err := os.Getwd()
		if err != nil {
			logrus.Error(err)
			errE := errors.UnexpectedInternal(err)
			w.WriteHeader(errE.HttpError)
			w.Write(errors.JSONError(errE))
			return
		}

		photoPath := "static/product/"
		os.Chdir(photoPath)

		photoID, err := uuid.NewRandom()
		if err != nil {
			logrus.Error(err)
			errE := errors.UnexpectedInternal(err)
			w.WriteHeader(errE.HttpError)
			w.Write(errors.JSONError(errE))
			return
		}

		f, err := os.OpenFile(photoID.String()+extension, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			logrus.Error(err)
			errE := errors.UnexpectedInternal(err)
			w.WriteHeader(errE.HttpError)
			w.Write(errors.JSONError(errE))
			return
		}
		defer f.Close()

		os.Chdir(str)

		_, err = io.Copy(f, file)
		if err != nil {
			_ = os.Remove(photoID.String() + extension)
			logrus.Error(err)
			errE := errors.UnexpectedInternal(err)
			w.WriteHeader(errE.HttpError)
			w.Write(errors.JSONError(errE))
			return
		}

		imgs["linkImages"] = append(imgs["linkImages"], "/static/product/"+photoID.String()+extension)
	}

	//if len(imgs) == 0 {
	//	logrus.Error(err)
	//	w.WriteHeader(http.StatusBadRequest)
	//	w.Write(errors.JSONError(errors.Error{Message: "http: no such file"}.Error()))
	//	return
	//}

	_, errE := ph.productUcase.UpdatePhoto(productID, imgs["linkImages"])
	if errE != nil {
		logrus.Error(errE.Message)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	body, err := json.Marshal(imgs)
	if err != nil {
		logrus.Error(err)
		errE := errors.UnexpectedInternal(err)
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
