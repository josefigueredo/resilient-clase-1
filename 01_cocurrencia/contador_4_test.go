package contador_test

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestContadorConcurrenteConWaitGroupConMutex incrementa un contador con goroutines, waitgroup y mutex
func TestContadorConcurrenteConWaitGroupConMutex(t *testing.T) {
	var contador int64
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := 0; i < Limite; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(CargaDeTrabajo)
			mu.Lock()
			contador++
			mu.Unlock()
		}()
	}
	wg.Wait()

	if contador != Limite {
		t.Errorf("el contador debía llegara a %d, y llegó a: %d", Limite, contador)
	} else {
		fmt.Println(fmt.Sprintf("%s: el contador alcanzo el límite: %d", t.Name(), contador))
	}
}

// TestContadorConcurrenteConWaitGroupConAtomic incrementa un contador con goroutines, waitgroup y atomic wrappers
func TestContadorConcurrenteConWaitGroupConAtomic(t *testing.T) {
	var contador int64
	var wg sync.WaitGroup

	for i := 0; i < Limite; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(CargaDeTrabajo)
			atomic.AddInt64(&contador, 1)
		}()
	}
	wg.Wait()

	if contador != Limite {
		t.Errorf("el contador debía llegara a %d, y llegó a: %d", Limite, contador)
	} else {
		fmt.Println(fmt.Sprintf("%s: el contador alcanzo el límite: %d", t.Name(), contador))
	}
}
