package worker

import (
	"context"
	"crypto/sha256"
	"fmt"
	"sync"
	"time"

	"github.com/user/golang_basics/internal/platform/logger"
)

// Pool orquesta un conjunto de trabajadores (workers) que procesan jobs en paralelo. 🏊‍♂️
// Demuestra los patrones: Goroutines (Mod 8), Pooling (Mod 10) y Fan Out Bounded (Mod 10).
type Pool struct {
	numWorkers int
	stats      *Stats
	logger     *logger.Logger
}

// NewPool crea una instancia del pool de trabajadores. 🏗️
func NewPool(numWorkers int, l *logger.Logger) *Pool {
	return &Pool{
		numWorkers: numWorkers,
		stats:      NewStats(),
		logger:     l,
	}
}

// Run lanza el orquestador y los trabajadores. 🚀
func (p *Pool) Run(ctx context.Context, jobs []Job) *Stats {
	// a) Creamos el canal de comunicación. 📦
	// Usamos un buffer pequeño para permitir que el productor siempre lleve un poco de ventaja.
	// Si lo pusiéramos unbuffered (make(chan Job)), el productor se bloquearía
	// hasta que un worker esté libre para recibir exactamente ese job.
	ch := make(chan Job, p.numWorkers)

	// b) El Productor: envía los jobs al canal. 🏗️
	go func() {
		// Al cerrar el canal, enviamos la señal de "Pooling" a los trabajadores. 🛑
		// Cuando el canal se vacíe, los workers saldrán de su loop 'range'.
		defer close(ch)

		for _, job := range jobs {
			select {
			case ch <- job:
				// Job enviado a la cola.
			case <-ctx.Done():
				// Si el contexto expira, dejamos de enviar inmediatamente. 🛑
				return
			}
		}
	}()

	// c) Los Workers: procesan lo que llega por el canal (Fan Out Bounded). 👷‍♂️
	var wg sync.WaitGroup
	for i := 0; i < p.numWorkers; i++ {
		workerID := fmt.Sprintf("worker-%02d", i)
		wg.Add(1)

		go func(id string) {
			defer wg.Done()

			// Patrón Pooling: iteramos sobre el canal hasta que se cierre.
			for job := range ch {
				// Comprobamos si el contexto ha sido cancelado antes de empezar una tarea larga. 🛑
				if ctx.Err() != nil {
					return
				}
				p.process(ctx, id, job)
			}
		}(workerID)
	}

	// d) Esperamos a que todos los workers terminen su jornada laboral. ⏳
	wg.Wait()

	return p.stats
}

// process simula el trabajo pesado de cada worker. ⚙️
func (p *Pool) process(ctx context.Context, workerID string, job Job) {
	// --- Fase CPU-bound: Demuestra el consumo de ciclos de procesador. 🌋 ---
	// Calculamos el hash SHA256 muchas veces para simular carga real de CPU.
	startCPU := time.Now()
	data := []byte(job.Name + job.Email)
	for i := 0; i < 5000; i++ {
		h := sha256.New()
		h.Write(data)
		_ = h.Sum(nil)
	}
	cpuDur := time.Since(startCPU)

	// --- Fase IO-bound: Demuestra la espera de servicios externos. 🌊 ---
	// Simulamos el guardado en una base de datos con un Sleep (pero respetando cancelación).
	startIO := time.Now()
	select {
	case <-time.After(3 * time.Millisecond):
		// Simulamos que el "guardado" tardó 3ms. ✅
	case <-ctx.Done():
		// Si se cancela mientras esperamos "al disco", salimos. 🛑
		return
	}
	ioDur := time.Since(startIO)

	// Registramos estadísticas y loggeamos lo sucedido.
	p.stats.RecordProcessed(workerID)
	p.logger.Log(fmt.Sprintf("[%s] processed %s (cpu: %v, io: %v)", workerID, job.Name, cpuDur, ioDur))
}

// Stats retorna las estadísticas actuales del pool. 📊
func (p *Pool) Stats() *Stats {
	return p.stats
}
