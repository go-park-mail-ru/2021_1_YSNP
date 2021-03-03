package main

import (
	"2021_1_YSNP/middleware"
	_signIn "2021_1_YSNP/handlers/SignIn"
	_signUp "2021_1_YSNP/handlers/SignUp"
	_mainPage "2021_1_YSNP/handlers/MainPage"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func testFunc (resp http.ResponseWriter, req *http.Request){
	resp.Write([]byte("Hr"))
}


func main(){
		router := mux.NewRouter()

		router.Use(middleware.CorsControlMiddleware)

		server := http.Server{
			Addr:              ":8080",
			Handler: 		   router,
			ReadTimeout:       60 * time.Second,
			WriteTimeout:      60 * time.Second,
		}

		router.HandleFunc("/", _mainPage.MainPageHandler).Methods(http.MethodGet)
		router.HandleFunc("/signup", _signUp.SignUpHandler).Methods(http.MethodPost)
		router.HandleFunc("/signin", _signIn.SignInHandler).Methods(http.MethodGet, http.MethodPost)
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}


}