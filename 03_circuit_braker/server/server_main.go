package main

import (
	"fmt"
	"net/http"
	"sync"
)

const (
	requestLimitNewLine = 60
	serverPort          = 4000
)

var (
	healthy                            bool
	requestCount                       int
	countMutex                         sync.Mutex
	statusByHealth, visualClueByHealth map[bool]string
)

func init() {
	healthy = true
	requestCount = 0
	statusByHealth = make(map[bool]string)
	statusByHealth[true] = ":-)"
	statusByHealth[false] = ":-("
	visualClueByHealth = make(map[bool]string)
	visualClueByHealth[true] = "."
	visualClueByHealth[false] = "x"
}

func main() {
	router := http.NewServeMux()

	router.HandleFunc("/ping", healthCheckHandler)
	router.HandleFunc("/change", healthChangeHandler)

	fmt.Println(fmt.Sprintf("Server iniciado y escuchando en el puerto:%d", serverPort))
	http.ListenAndServe(fmt.Sprintf(":%d", serverPort), router)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	defer incrementRequestCount()

	fmt.Print(visualClueByHealth[healthy])
	if !healthy {
		http.Error(w, statusByHealth[healthy], http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w)
}

func healthChangeHandler(w http.ResponseWriter, r *http.Request) {
	defer incrementRequestCount()

	healthy = !healthy
	fmt.Print("|")
	fmt.Fprint(w, fmt.Sprintf("%s\n", statusByHealth[healthy]))
}

func incrementRequestCount() {
	countMutex.Lock()
	requestCount++
	if requestCount >= requestLimitNewLine {
		requestCount = 0
		fmt.Println()
	}
	countMutex.Unlock()
}
