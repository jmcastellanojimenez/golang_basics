package worker

import (
	"sync"
	"sync/atomic"
)

// Stats se encarga de recolectar métricas del procesamiento. 📊
// Demuestra el uso de contadores atómicos (Módulo 9) para rendimiento máximo
// y RWMutex para proteger estructuras complejas como mapas.
type Stats struct {
	processed atomic.Int64    // Contador global. Rápido y sin contención de cerrojos (lock-free). ⚡
	mu        sync.RWMutex    // Protege el mapa perWorker. Permite múltiples lectores pero solo un escritor. 🔐
	perWorker map[string]int  // Estadísticas detalladas: ¿cuántos procesó cada worker?
}

// NewStats inicializa un recolector de estadísticas nuevo. 🏗️
func NewStats() *Stats {
	return &Stats{
		perWorker: make(map[string]int),
	}
}

// RecordProcessed registra un éxito en el procesamiento. ✅
func (s *Stats) RecordProcessed(workerID string) {
	// 1. Incrementar el contador global de forma atómica.
	// Es mucho más eficiente que un Mutex para un simple contador.
	s.processed.Add(1)

	// 2. Actualizar el mapa de workers.
	// En Go, los mapas NO son thread-safe para escrituras concurrentes. 🚫
	// Necesitamos bloquear para evitar un "fatal error: concurrent map iteration and map write".
	s.mu.Lock()
	s.perWorker[workerID]++
	s.mu.Unlock()
}

// TotalProcessed retorna el total acumulado hasta el momento. 🔢
func (s *Stats) TotalProcessed() int64 {
	// Lectura atómica segura.
	return s.processed.Load()
}

// PerWorker retorna una copia de las estadísticas por worker. 📋
// Usamos RWMutex para permitir lecturas concurrentes.
func (s *Stats) PerWorker() map[string]int {
	// RLock permite que otros hilos lean simultáneamente, pero bloquea si alguien quiere escribir.
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Importante: retornamos una COPIA del mapa. 
	// Si retornamos el original, el llamador podría tener una data race si nosotros seguimos escribiendo.
	copyMap := make(map[string]int, len(s.perWorker))
	for k, v := range s.perWorker {
		copyMap[k] = v
	}
	return copyMap
}
