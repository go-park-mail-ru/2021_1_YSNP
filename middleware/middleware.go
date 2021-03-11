package middleware

import (
	"net/http"
)

func CorsControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		// TODO(Sergey) localhost
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if req.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, req)
	})
}
