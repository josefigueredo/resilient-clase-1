package contador_test

import (
	"fmt"
	"testing"
	"time"
)

// valores que usaremos en todos los tests
const (
	Limite         = 1000000             // valor maximo que queremos alcanzar con el contador
	CargaDeTrabajo = 1 * time.Nanosecond // representa una carga de trabajo adicional previo a incrementar el contador
)

// TestContadorSecuencial incrementa un contador dentro de un loop
func TestContadorSecuencial(t *testing.T) {
	var contador int64

	for i := 0; i < Limite; i++ {
		time.Sleep(CargaDeTrabajo)
		contador++
	}

	if contador != Limite {
		t.Errorf("el contador NO alcanzó el límite: %d", contador)
	} else {
		fmt.Println(fmt.Sprintf("%s: el contador alcanzo el límite!", t.Name()))
	}
}
