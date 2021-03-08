package main

import (
	_login "2021_1_YSNP/handlers/login"
	_mainPage "2021_1_YSNP/handlers/main_page"
	_product "2021_1_YSNP/handlers/product"
	_profile "2021_1_YSNP/handlers/profile"
	_signUp "2021_1_YSNP/handlers/signup"
	"2021_1_YSNP/middleware"
	_tmpDB "2021_1_YSNP/tmp_database"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	_tmpDB.InitDB()
	router := mux.NewRouter()

	router.Use(middleware.CorsControlMiddleware)

	server := http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	router.HandleFunc("/api/v1/", _mainPage.MainPageHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/product/{id}", _product.ProductIDHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/product/create", _product.ProductCreateHandler).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/product/upload", _product.UploadPhotoHandler).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/login", _login.LoginHandler).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/logout", _login.LogoutHandler).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/signup", _signUp.SignUpHandler).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/upload", _signUp.UploadAvatarHandler).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/me", _profile.GetProfileHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/settings", _profile.ChangeProfileHandler).Methods(http.MethodPost)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
