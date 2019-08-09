package main

import (
	"fmt"
	"net/http"

	"golang.org/x/time/rate" // obtener la dependencia con go get -u golang.org/x/time/rate
)

// Configuramos nuestro Rate Limit con Token Buckets
var limiter = rate.NewLimiter(2, 1)

// Middleware de los requests
func limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validamos si hay tokens disponibles,
		// si no devolvemos un HTTP Status 429
		if limiter.Allow() == false {
			http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", okHandler)

	// Levantamos un server con el middleware
	fmt.Println("Starting server...")
	http.ListenAndServe(":4000", limit(mux))
}

func okHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
