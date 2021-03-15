package main

import (
	_login "2021_1_YSNP/handlers/login"
	_mainPage "2021_1_YSNP/handlers/main_page"
	_product "2021_1_YSNP/handlers/product"
	_profile "2021_1_YSNP/handlers/profile"
	_signUp "2021_1_YSNP/handlers/signup"
	"2021_1_YSNP/middleware"
	_tmpDB "2021_1_YSNP/tmp_database"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	_tmpDB.InitDB()

	router := mux.NewRouter()

	logrus.SetFormatter(&logrus.TextFormatter{DisableColors: true})
	logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
		"host":   "89.208.199.170",
		"port":   "8080",
	}).Info("Starting server")

	AccessLogOut := new(middleware.AccessLogger)

	contextLogger := logrus.WithFields(logrus.Fields{
		"mode":   "[access_log]",
		"logger": "LOGRUS",
	})
	logrus.SetFormatter(&logrus.JSONFormatter{})
	AccessLogOut.LogrusLogger = contextLogger

	router.Use(AccessLogOut.AccessLogMiddleware)


	router.Use(middleware.CorsControlMiddleware)

	server := http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/product/list", _mainPage.MainPageHandler).Methods(http.MethodGet)
	api.HandleFunc("/product/{id:[0-9]+}", _product.ProductIDHandler).Methods(http.MethodGet)
	api.HandleFunc("/product/create", _product.ProductCreateHandler).Methods(http.MethodPost)
	api.HandleFunc("/product/upload", _product.UploadPhotoHandler).Methods(http.MethodPost)
	api.HandleFunc("/login", _login.LoginHandler).Methods(http.MethodPost)
	api.HandleFunc("/logout", _login.LogoutHandler).Methods(http.MethodPost)
	api.HandleFunc("/signup", _signUp.SignUpHandler).Methods(http.MethodPost)
	api.HandleFunc("/upload", _signUp.UploadAvatarHandler).Methods(http.MethodPost)
	api.HandleFunc("/me", _profile.GetProfileHandler).Methods(http.MethodGet)
	api.HandleFunc("/settings", _profile.ChangeProfileHandler).Methods(http.MethodPost)
	api.HandleFunc("/settings/password", _profile.ChangeProfilePasswordHandler).Methods(http.MethodPost)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
