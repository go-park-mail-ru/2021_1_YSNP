package main

import (
	"database/sql"
	_ "github.com/jackc/pgx/stdlib"

	userHandler "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/user/delivery/http"
	userRepo "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/user/repository/postgres"
	userUsecase "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/user/usecase"

	sessionHandler "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/session/delivery/http"
	sessionRepo "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/session/repository/postgres"
	sessionUsecase "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/session/usecase"

	productHandler "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/product/delivery/http"
	productRepo "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/product/repository/postgres"
	productUsecase "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/product/usecase"

	"github.com/go-park-mail-ru/2021_1_YSNP/configs"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/middleware"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	configs, err := configs.LoadConfig("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	dbConn, err := sql.Open("pgx", configs.GetDBConfig())
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	if err := dbConn.Ping(); err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	//_tmpDB.InitDB()

	userRepo := userRepo.NewUserRepository(dbConn)
	sessRepo := sessionRepo.NewSessionRepository(dbConn)
	prodRepo := productRepo.NewProductRepository(dbConn)

	userUcase := userUsecase.NewUserUsecase(userRepo)
	sessUcase := sessionUsecase.NewSessionUsecase(sessRepo)
	prodUcase := productUsecase.NewProductUsecase(prodRepo)

	mw := middleware.NewMiddleware(sessUcase, userUcase)

	logrus.SetFormatter(&logrus.TextFormatter{DisableColors: true})
	logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
		"host":   "89.208.199.170",
		"port":   "8080",
	}).Info("Starting server")

	router.Use(mw.AccessLogMiddleware)

	router.Use(middleware.CorsControlMiddleware)

	server := http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	api := router.PathPrefix("/api/v1").Subrouter()

	userHandler := userHandler.NewUserHandler(userUcase, sessUcase)
	sessHandler := sessionHandler.NewSessionHandler(sessUcase, userUcase)
	prodHandler := productHandler.NewProductHandler(prodUcase)

	userHandler.Configure(api)
	sessHandler.Configure(api)
	prodHandler.Configure(api)
	//api.HandleFunc("/product/list", _mainPage.MainPageHandler).Methods(http.MethodGet)
	//api.HandleFunc("/product/{id:[0-9]+}", _product.ProductIDHandler).Methods(http.MethodGet)
	//api.HandleFunc("/product/create", _product.ProductCreateHandler).Methods(http.MethodPost)
	//api.HandleFunc("/product/upload", _product.UploadPhotoHandler).Methods(http.MethodPost)
	//api.HandleFunc("/login", _login.LoginHandler).Methods(http.MethodPost)
	//api.HandleFunc("/logout", _login.LogoutHandler).Methods(http.MethodPost)
	//api.HandleFunc("/signup", _signUp.SignUpHandler).Methods(http.MethodPost)
	//api.HandleFunc("/upload", _signUp.UploadAvatarHandler).Methods(http.MethodPost)
	//api.HandleFunc("/me", _profile.GetProfileHandler).Methods(http.MethodGet)
	//api.HandleFunc("/settings", _profile.ChangeProfileHandler).Methods(http.MethodPost)
	//api.HandleFunc("/settings/password", _profile.ChangeProfilePasswordHandler).Methods(http.MethodPost)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
