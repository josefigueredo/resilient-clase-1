package contador_test

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// TestContadorConcurrenteConWaitGroup incrementa un contador con goroutines y waitgroup
func TestContadorConcurrenteConWaitGroup(t *testing.T) {
	var contador int64
	var wg sync.WaitGroup

	for i := 0; i < Limite; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(CargaDeTrabajo)
			contador++
		}()
	}
	wg.Wait()

	if contador != Limite {
		t.Errorf("el contador NO alcanzó el límite: %d", contador)
	} else {
		fmt.Println(fmt.Sprintf("%s: el contador alcanzo el límite!", t.Name()))
	}
}
