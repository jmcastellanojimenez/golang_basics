package worker

import (
	"context"
	"fmt"
	"io"
	"sync"
	"testing"

	"github.com/user/golang_basics/internal/platform/logger"
)

// 🏎️ BenchmarkAtomic vs BenchmarkMutex
// Compara el coste de incrementar un contador desde múltiples goroutines.
// Spoiler: Atomic suele ganar porque no suspende goroutines (es lock-free).
func BenchmarkAtomic(b *testing.B) {
	stats := NewStats()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			stats.processed.Add(1)
		}
	})
}

func BenchmarkMutex(b *testing.B) {
	var mu sync.Mutex
	var count int64
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.Lock()
			count++
			mu.Unlock()
		}
	})
}

// 🎯 BenchmarkOptimalWorkers
// ¿Cuántos workers son demasiados?
// Este benchmark ejecuta el Pool con distintas cantidades de trabajadores
// para encontrar el punto óptimo donde la latencia baja pero el throughput sube.
func BenchmarkOptimalWorkers(b *testing.B) {
	// 1. Silenciamos el logger para no benchmarkear la consola.
	l := logger.New(io.Discard, 1000)
	defer l.Shutdown()

	// 2. Generamos una carga de trabajo fija (1000 jobs).
	jobs := GenerateJobs(1000)
	ctx := context.Background()

	// 3. Probamos distintas configuraciones.
	workerCounts := []int{1, 2, 4, 8, 16, 32, 64}

	for _, count := range workerCounts {
		b.Run(fmt.Sprintf("Workers-%d", count), func(b *testing.B) {
			p := NewPool(count, l)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				p.Run(ctx, jobs)
			}
		})
	}
}
