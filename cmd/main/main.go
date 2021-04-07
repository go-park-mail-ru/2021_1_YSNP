package main

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/tarantool/go-tarantool"

	userHandler "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/user/delivery/http"
	userRepo "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/user/repository/postgres"
	userUsecase "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/user/usecase"

	sessionHandler "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/session/delivery/http"
	sessionRepo "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/session/repository/tarantool"
	sessionUsecase "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/session/usecase"

	productHandler "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/product/delivery/http"
	productRepo "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/product/repository/postgres"
	productUsecase "github.com/go-park-mail-ru/2021_1_YSNP/internal/app/product/usecase"

	"github.com/go-park-mail-ru/2021_1_YSNP/configs"

	"log"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/middleware"
	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

func main() {
	configs, err := configs.LoadConfig("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	sqlConn, err := sql.Open("pgx", configs.GetDBConfig())
	if err != nil {
		log.Fatal(err)
	}
	defer sqlConn.Close()

	if err := sqlConn.Ping(); err != nil {
		log.Fatal(err)
	}

	opts := tarantool.Opts{
		User: "admin",
		Pass: "pass",
	}
	tarConn, err := tarantool.Connect("127.0.0.1:3301", opts)

	if err != nil {
		fmt.Println("baa: Connection refused:", err)
		return
	}

	router := mux.NewRouter()

	userRepo := userRepo.NewUserRepository(sqlConn)
	sessRepo := sessionRepo.NewSessionRepository(tarConn)
	prodRepo := productRepo.NewProductRepository(sqlConn)

	userUcase := userUsecase.NewUserUsecase(userRepo)
	sessUcase := sessionUsecase.NewSessionUsecase(sessRepo)
	prodUcase := productUsecase.NewProductUsecase(prodRepo)

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "02-01-2006 15:04:05",
	})
	logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
		"host":   "89.208.199.170",
		"port":   "8080",
	}).Info("Starting server")

	mw := middleware.NewMiddleware(sessUcase, userUcase)
	contextLogger := logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	})
	mw.LogrusLogger = contextLogger
	logrus.SetLevel(logrus.DebugLevel)

	router.Use(mw.AccessLogMiddleware)

	router.Use(middleware.CorsControlMiddleware)

	server := http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	//router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	api := router.PathPrefix("/api/v1").Subrouter()

	userHandler := userHandler.NewUserHandler(userUcase, sessUcase)
	sessHandler := sessionHandler.NewSessionHandler(sessUcase, userUcase)
	prodHandler := productHandler.NewProductHandler(prodUcase)

	userHandler.Configure(api, mw)
	sessHandler.Configure(api, mw)
	prodHandler.Configure(api, mw)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
