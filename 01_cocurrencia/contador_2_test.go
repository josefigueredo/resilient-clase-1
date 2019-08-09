package contador_test

import (
	"fmt"
	"testing"
	"time"
)

// TestContadorConcurrente incrementa un contador con goroutines
func TestContadorConcurrente(t *testing.T) {
	var contador int64

	for i := 0; i < Limite; i++ {
		go func() {
			time.Sleep(CargaDeTrabajo)
			contador++
		}()
	}

	if contador != Limite {
		t.Errorf("el contador NO alcanzó el límite: %d", contador)
	} else {
		fmt.Println(fmt.Sprintf("%s: el contador alcanzo el límite!", t.Name()))
	}
}
