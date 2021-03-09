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

	apiv1 := router.PathPrefix("/api/v1").Subrouter()

	apiv1.HandleFunc("/", _mainPage.MainPageHandler).Methods(http.MethodGet)
	apiv1.HandleFunc("/product/{id}", _product.ProductIDHandler).Methods(http.MethodGet)
	apiv1.HandleFunc("/product/create", _product.ProductCreateHandler).Methods(http.MethodPost)
	apiv1.HandleFunc("/product/upload", _product.UploadPhotoHandler).Methods(http.MethodPost)
	apiv1.HandleFunc("/login", _login.LoginHandler).Methods(http.MethodPost)
	apiv1.HandleFunc("/logout", _login.LogoutHandler).Methods(http.MethodPost)
	apiv1.HandleFunc("/signup", _signUp.SignUpHandler).Methods(http.MethodPost)
	apiv1.HandleFunc("/upload", _signUp.UploadAvatarHandler).Methods(http.MethodPost)
	apiv1.HandleFunc("/me", _profile.GetProfileHandler).Methods(http.MethodGet)
	apiv1.HandleFunc("/settings", _profile.ChangeProfileHandler).Methods(http.MethodPost)
	apiv1.HandleFunc("/settings/password", _profile.ChangeProfilePasswordHandler).Methods(http.MethodPost)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
