package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/sony/gobreaker" // obtener la dependencia con: go get github.com/sony/gobreaker
)

const (
	requestLimitNewLine = 60
	useCircuitBreaker   = true
	targetAddress       = "localhost"
	targetPort          = 4000
	targetPath          = "ping"
)

var (
	ticker       *time.Ticker
	targetURL    string
	requestCount uint32
	st           gobreaker.Settings
	cb           *gobreaker.CircuitBreaker
)

func init() {
	targetURL = fmt.Sprintf("http://%s:%d/%s", targetAddress, targetPort, targetPath)
	requestCount = 0

	if useCircuitBreaker {
		st = gobreaker.Settings{
			Name:          "Health Check",  // The name of the CircuitBreaker.
			MaxRequests:   1,               // Maximum number of requests allowed to pass through when state is half-open.
			Interval:      0,               // Period of the closed state after clearing the internal Counts.
			Timeout:       5 * time.Second, // Time to stay on the open state, after which the state becomes half-open.
			OnStateChange: myOnStateChange, // is called whenever the state of the CircuitBreaker changes.
		}

		cb = gobreaker.NewCircuitBreaker(st)
	}
}

func myOnStateChange(name string, from gobreaker.State, to gobreaker.State) {
	if to == gobreaker.StateOpen {
		fmt.Print("|")
	} else if to == gobreaker.StateHalfOpen {
		fmt.Print("/")
	} else if to == gobreaker.StateClosed {
		fmt.Print("-")
	}

	incrementRequestCount()
}

// getWithCB wraps http.Get in CircuitBreaker.
func getWithCB(url string) (*http.Response, error) {
	resp, err := cb.Execute(func() (interface{}, error) {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode == http.StatusInternalServerError {
			return nil, fmt.Errorf("Internal Server Error")
		}

		return resp, nil
	})
	if err != nil {
		return nil, err
	}

	return resp.(*http.Response), nil
}

func main() {
	fmt.Println("Cliente iniciado e chequeando la salud del server")

	// inicializo un ticker para poder invocar ritmicamente al endpoint de healthCheck del server
	for ticker = time.NewTicker(200 * time.Millisecond); true; <-ticker.C {
		var err error
		if cb != nil {
			_, err = getWithCB(targetURL)
		} else {
			_, err = http.Get(targetURL)
		}
		if err != nil {
			fmt.Print("x")
		} else {
			fmt.Print(".")
		}

		incrementRequestCount()
	}
}

func incrementRequestCount() {
	atomic.AddUint32(&requestCount, 1)
	if atomic.CompareAndSwapUint32(&requestCount, requestLimitNewLine, 0) {
		fmt.Println()
	}
}
