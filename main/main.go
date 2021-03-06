package main

import (
	_login "2021_1_YSNP/handlers/Login"
	_mainPage "2021_1_YSNP/handlers/MainPage"
	_product "2021_1_YSNP/handlers/Product"
	_profile "2021_1_YSNP/handlers/Profile"
	_signUp "2021_1_YSNP/handlers/SignUp"
	"2021_1_YSNP/middleware"
	_tmpDB "2021_1_YSNP/tmp_database"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)


func main(){
		_tmpDB.InitDB()
		router := mux.NewRouter()

		router.Use(middleware.CorsControlMiddleware)

		server := http.Server{
			Addr:              ":8080",
			Handler: 		   router,
			ReadTimeout:       60 * time.Second,
			WriteTimeout:      60 * time.Second,
		}

		router.HandleFunc("/", _mainPage.MainPageHandler).Methods(http.MethodGet)
		router.HandleFunc("/product/{id}", _product.ProductIDHandler).Methods(http.MethodGet)
		router.HandleFunc("/product/create", _product.ProductCreateHandler).Methods(http.MethodPost)
		router.HandleFunc("/login", _login.LoginHandler).Methods(http.MethodPost)
		router.HandleFunc("/signup", _signUp.SignUpHandler).Methods(http.MethodPost)
		router.HandleFunc("/me", _profile.GetProfileHandler).Methods(http.MethodGet)
		router.HandleFunc("/settings", _profile.ChangeProfileHandler).Methods(http.MethodPost)
	err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}


}