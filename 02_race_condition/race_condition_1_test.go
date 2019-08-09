package racecondition_test

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// valores que usaremos en todos los tests
const (
	Limite         = 100                 // valor maximo que queremos alcanzar con el contador
	CargaDeTrabajo = 1 * time.Nanosecond // representa una carga de trabajo adicional previo a incrementar el contador
)

// HighLowBlock almacena los limites del bloque de numeros en uso para obtener consecutivos
// una vez que se usa el ultimo numero del bloque, debe renovarse por un nuevo bloque
type HighLowBlock struct {
	Low  uint64 // limite inferior del bloque de numeros en uso
	High uint64 // limite superior del bloque de numeros en uso
}

var (
	blockSize         uint64 = 10 // tamaño del bloque
	renewBlockMutex   sync.Mutex
	highLowBlockInUse HighLowBlock                    // bloque en uso por esta instancia
	lastUsedNumber    uint64                          // ultimo numero usado por esta instancia
	externalStorage   = make(map[string]HighLowBlock) // representa el almacenamiento externo
)

// GetNextNumber devuelve el siguiente numero a ser utilizado
func getNextNumber() (uint64, error) {
	var err error
	var newLow uint64
	var newHigh uint64
	var nextNumber uint64

	// queremos evitar race conditions mientras chequeamos/actualizamos el  HighLowBlock
	renewBlockMutex.Lock()
	if lastUsedNumber == highLowBlockInUse.High {
		// si alcanzamos el final del bloque, obtenemos un nuevo bloque desde el almacenamiento externo
		if newLow, newHigh, err = getNextBlock(); err != nil {
			renewBlockMutex.Unlock()
			return 0, err
		}

		// pasamos los valores recibidos al HighLowBlock
		highLowBlockInUse.Low = newLow
		highLowBlockInUse.High = newHigh

		// re-initializamos con el limite inferior del bloque de numeros
		lastUsedNumber = newLow
	}
	renewBlockMutex.Unlock()

	// usamos sync.atomic para actualizar el contador
	nextNumber = atomic.AddUint64(&lastUsedNumber, 1)

	return nextNumber, nil
}

// getNextBlock obtiene los nuevos limites desde un almacenamiento externo
func getNextBlock() (newLow uint64, newHigh uint64, err error) {
	var externalMutex sync.Mutex
	var highLowBlock HighLowBlock
	var foundBlock bool
	var blockKey = "hlb"

	externalMutex.Lock()         // representa la llamada al servicio externo de Mutex para obtener un Lock
	defer externalMutex.Unlock() // representa la llamada al servicio externo de Mutex para liberar el Lock

	fmt.Println()

	// busco el bloque desde el almacenamiento externo
	// (debería ser el ultimo bloque utilizado por alguna instancia)
	if highLowBlock, foundBlock = externalStorage[blockKey]; !foundBlock {
		// si no lo encontré entonces lo inicializo para poder guardarlo
		highLowBlock = HighLowBlock{Low: 0, High: 0}
	}

	// aquí sucede la magia de la renovación del bloque
	highLowBlock.Low = highLowBlock.High
	highLowBlock.High = highLowBlock.Low + blockSize

	// guardamos el nuevo bloque en almacenamiento externo
	externalStorage[blockKey] = highLowBlock

	fmt.Print(fmt.Sprintf("getNextBlock: %+v", highLowBlock))

	return highLowBlock.Low, highLowBlock.High, nil
}

func TestGetNextNumber(t *testing.T) {
	lastUsedNumber = 0

	for i := 0; i < Limite; i++ {
		time.Sleep(CargaDeTrabajo)
		next, _ := getNextNumber()
		fmt.Printf(" %d", next)
	}

	fmt.Println()
	if lastUsedNumber != Limite {
		t.Errorf("el contador NO alcanzó el límite: %d", lastUsedNumber)
	} else {
		fmt.Println(fmt.Sprintf("%s: el contador alcanzo el límite! (%+v)", t.Name(), highLowBlockInUse))
	}
}
