package racecondition_test

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

const (
	Workers = 3 // representa la cantidad de solicitudes concurrentes que obtienen números
)

func TestGetNextNumberConcurrenteConWaitGroup(t *testing.T) {
	var wg sync.WaitGroup

	lastUsedNumber = 0

	for j := 0; j < Workers; j++ {
		wg.Add(1)
		go func(workerName int) {
			defer wg.Done()
			for i := 0; i < Limite; i++ {
				time.Sleep(CargaDeTrabajo)
				next, _ := getNextNumber()
				fmt.Printf(" %d:%d", workerName, next)
			}
		}(j)
	}
	wg.Wait()

	fmt.Println()
	if lastUsedNumber != Limite*Workers {
		t.Errorf("el contador NO alcanzó el límite: %d", lastUsedNumber)
	} else {
		fmt.Println(fmt.Sprintf("%s: el contador alcanzo el límite!", t.Name()))
	}
	if lastUsedNumber != highLowBlockInUse.High {
		t.Errorf("el contador difiere del bloque: %+v <> %d", highLowBlockInUse, lastUsedNumber)
	}
}
