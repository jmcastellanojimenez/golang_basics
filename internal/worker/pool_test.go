package worker

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/user/golang_basics/internal/platform/logger"
)

// TestPoolProcessesAllJobs verifica que el pool procese correctamente todos los elementos. ✅
func TestPoolProcessesAllJobs(t *testing.T) {
	var buf bytes.Buffer
	l := logger.New(&buf, 256)
	numWorkers := 4
	numJobs := 100
	
	pool := NewPool(numWorkers, l)
	jobs := GenerateJobs(numJobs)
	
	// Contexto holgado para dar tiempo.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stats := pool.Run(ctx, jobs)

	// Verificaciones de integridad.
	if stats.TotalProcessed() != int64(numJobs) {
		t.Errorf("esperaba %d procesados, obtuvo %d", numJobs, stats.TotalProcessed())
	}

	// Comprobamos el mapa por worker: la suma debe coincidir con el total.
	perWorker := stats.PerWorker()
	totalFromMap := 0
	for _, count := range perWorker {
		totalFromMap += count
	}

	if totalFromMap != numJobs {
		t.Errorf("la suma de stats por worker (%d) no coincide con el total (%d)", totalFromMap, numJobs)
	}

	l.Shutdown()
}

// TestPoolRespectsTimeout asegura que el pool se detenga si el contexto expira. 🛡️
func TestPoolRespectsTimeout(t *testing.T) {
	var buf bytes.Buffer
	l := logger.New(&buf, 10) // Buffer pequeño
	numWorkers := 2
	numJobs := 10000 // Muchos jobs para forzar el timeout.

	pool := NewPool(numWorkers, l)
	jobs := GenerateJobs(numJobs)

	// Timeout extremadamente agresivo.
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	start := time.Now()
	stats := pool.Run(ctx, jobs)
	duration := time.Since(start)

	// El test debe terminar rápido (timeout).
	if duration > 2*time.Second {
		t.Errorf("el test tardó demasiado (%v), el timeout no funcionó correctamente", duration)
	}

	// No debe haber procesado todos los jobs.
	if stats.TotalProcessed() >= int64(numJobs) {
		t.Errorf("procesó todos los jobs (%d) a pesar del timeout", stats.TotalProcessed())
	}

	l.Shutdown()
}
