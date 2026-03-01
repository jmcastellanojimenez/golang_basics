//go:build race_example

package worker

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

// TestIntentionalRace demuestra lo que ocurre cuando NO protegemos los datos compartidos. 💥
// Si ejecutas este test con la flag '-race', Go detectará el conflicto inmediatamente.
// Ejecución: go test -tags=race_example -race ./internal/worker/
func TestIntentionalRace(t *testing.T) {
	counter := 0 // Variable compartida SIN protección. 🔓
	var wg sync.WaitGroup
	n := 1000

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// ¡DATA RACE! Múltiples goroutines leen, incrementan y escriben al mismo tiempo.
			counter++
		}()
	}

	wg.Wait()

	// Lo más probable es que 'counter' sea menor que 1000 debido a las actualizaciones perdidas.
	fmt.Printf("--- RESULTADO CON DATA RACE: %d (esperado: %d) ---\n", counter, n)
	
	/*
	   NOTA DIDÁCTICA: "Esto es un data race intencional. Ejecuta con 
	   go test -tags=race_example -race ./internal/worker/ para ver el detector 
	   en acción. Luego mira stats.go para ver cómo se soluciona con atomic y mutex." 💡
	*/
}

// TestFixedWithAtomic demuestra la solución usando operaciones atómicas. 🛡️
func TestFixedWithAtomic(t *testing.T) {
	var counter atomic.Int64 // Variable compartida PROTEGIDA. 🔒
	var wg sync.WaitGroup
	n := 1000

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Operación atómica: el hardware garantiza que el incremento sea indivisible. ⚡
			counter.Add(1)
		}()
	}

	wg.Wait()

	if counter.Load() != int64(n) {
		t.Errorf("esperaba %d, obtuvo %d", n, counter.Load())
	}

	fmt.Printf("--- RESULTADO CON ATOMIC: %d (siempre correcto) ---\n", counter.Load())
	
	/*
	   NOTA DIDÁCTICA: "Misma operación, protegida con atomic. 
	   go test -tags=race_example -race pasa limpio." ✅
	*/
}
