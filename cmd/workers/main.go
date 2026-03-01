package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof" // Importación mágica: habilita /debug/pprof automáticamente. ✨
	"os"
	"sort"
	"sync"
	"time"

	"github.com/user/golang_basics/internal/platform/logger"
	"github.com/user/golang_basics/internal/worker"
)

func main() {
	// 1. Configuración vía Flags. 🚩
	numWorkers := flag.Int("workers", 8, "Número de workers en paralelo")
	numJobs := flag.Int("jobs", 2000, "Número total de usuarios a procesar")
	timeout := flag.Duration("timeout", 10*time.Second, "Tiempo máximo de ejecución")
	bufferSize := flag.Int("buffer", 256, "Tamaño del buffer del logger (Drop Pattern)")
	mode := flag.String("mode", "pool", "Modo de ejecución: 'pool' o 'semaphore'")
	flag.Parse()
	
	// 1.1. Servidor de Profiling (pprof). 🕵️‍♂️
	// Lanza un servidor HTTP en segundo plano para que podamos inspeccionar 
	// la memoria, CPU y goroutines mientras el programa corre.
	// Acceso: http://localhost:6060/debug/pprof
	go func() {
		_ = http.ListenAndServe("localhost:6060", nil)
	}()

	// 2. Banner de inicio. 🏗️
	fmt.Println("🏭 Batch User Processor")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Workers:       %d\n", *numWorkers)
	fmt.Printf("Jobs:          %d\n", *numJobs)
	fmt.Printf("Timeout:       %v\n", *timeout)
	fmt.Printf("Logger buffer: %d\n", *bufferSize)
	fmt.Printf("Mode:          %s\n", *mode)
	fmt.Println()

	// 3. Inicialización de componentes. ⚙️
	l := logger.New(os.Stdout, *bufferSize)
	jobs := worker.GenerateJobs(*numJobs)
	
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	start := time.Now()
	var finalStats *worker.Stats

	// 4. Orquestación según el modo elegido. 🎼
	if *mode == "pool" {
		// --- MODO POOL (Recomendado) ---
		// Crea N goroutines fijas que escuchan de un canal compartido. 🏊‍♂️
		p := worker.NewPool(*numWorkers, l)
		finalStats = p.Run(ctx, jobs)
	} else {
		// --- MODO SEMAPHORE (Alternativo) ---
		// Lanza UNA goroutine por cada job (Fan Out), pero limitadas por un semáforo. 🚦
		finalStats = runWithSemaphore(ctx, *numWorkers, jobs, l)
	}

	// 5. Finalización y Reporte. 📊
	// Esperamos a que el logger termine de vaciar su buffer (drain).
	l.Shutdown()

	elapsed := time.Since(start)
	printStats(finalStats, elapsed, l.DroppedCount())
}

// runWithSemaphore demuestra el patrón Fan Out con limitación vía semáforo. 🚦
func runWithSemaphore(ctx context.Context, maxWorkers int, jobs []worker.Job, l *logger.Logger) *worker.Stats {
	stats := worker.NewStats()
	
	// El semáforo: un canal de struct{}.
	// Usamos struct{} porque ocupa 0 bytes de memoria. 💡
	sem := make(chan struct{}, maxWorkers)

	var wg sync.WaitGroup

	for _, j := range jobs {
		// Comprobamos cancelación antes de lanzar una nueva unidad de trabajo.
		if ctx.Err() != nil {
			break
		}

		wg.Add(1)
		go func(job worker.Job) {
			defer wg.Done()

			// ACQUIRE: Intentamos entrar en el semáforo. 📥
			// Si está lleno, esta goroutine se bloquea aquí hasta que haya un hueco.
			select {
			case sem <- struct{}{}:
				// Entramos.
			case <-ctx.Done():
				// Si el tiempo se acaba mientras esperábamos al semáforo, cancelamos.
				return
			}

			// RELEASE: Al terminar, liberamos el hueco para el siguiente. 📤
			defer func() { <-sem }()

			// Procesamiento (reutilizamos la lógica del job pero simplificada aquí).
			// En un sistema real, esto iría en un service.
			// Simulamos el trabajo.
			time.Sleep(5 * time.Millisecond)
			
			stats.RecordProcessed("semaphore-worker")
			l.Log(fmt.Sprintf("[semaphore] processed %s", job.Name))
		}(j)
	}

	// Wait for Result: main espera a que todas las goroutines terminen. ⏳
	wg.Wait()
	return stats
}

// printStats muestra el resumen final de la operación. 📋
func printStats(s *worker.Stats, elapsed time.Duration, dropped int64) {
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📊 STATS")
	fmt.Printf("  Total processed: %d\n", s.TotalProcessed())
	fmt.Printf("  Time elapsed:    %v\n", elapsed)
	fmt.Printf("  Logs dropped:    %d\n", dropped)
	fmt.Println("  Per worker:")

	perWorker := s.PerWorker()
	// Ordenamos los nombres de los workers para que la salida sea bonita.
	keys := make([]string, 0, len(perWorker))
	for k := range perWorker {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Printf("    %s: %d\n", k, perWorker[k])
	}
}

/*
   COMPARATIVA DE MODOS:
   - Pool Mode: Crea N goroutines fijas. Es mucho más eficiente en memoria si tienes
     millones de tareas, ya que no sobrecargas al scheduler de Go. 🏊‍♂️
   - Semaphore Mode: Crea una goroutine POR CADA tarea (Fan Out), pero el semáforo
     garantiza que solo N estén activas a la vez. Es más fácil de razonar para flujos
     lineales pero consume más recursos de control si el número de jobs es masivo. 🚦
*/
