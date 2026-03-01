package logger

import (
	"io"
	"sync"
	"sync/atomic"
)

// Logger representa un sistema de registro de eventos que implementa el Drop Pattern. 🛡️
// Se asegura de que la aplicación nunca se bloquee por culpa del sistema de logs,
// prefiriendo descartar (drop) mensajes si el sistema de salida está saturado.
type Logger struct {
	ch      chan string      // Canal con buffer donde se encolan los mensajes. 📥
	writer  io.Writer        // Destino final de los logs (ej: os.Stdout o un archivo). 🖋️
	dropped atomic.Int64    // Contador atómico de mensajes descartados por saturación. 📉
	wg      sync.WaitGroup   // Sincronizador para asegurar que el drenaje termine antes del cierre total. ⏳
}

// New crea una instancia de Logger con un buffer de tamaño específico. 🏗️
// Lanza una goroutine en segundo plano que se encarga de procesar la cola.
func New(w io.Writer, bufferSize int) *Logger {
	l := &Logger{
		ch:     make(chan string, bufferSize),
		writer: w,
	}

	// Añadimos 1 al WaitGroup antes de lanzar la goroutine trabajadora.
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		// Iterar con range sobre el channel es el patrón recomendado:
		// la goroutine morirá automáticamente cuando el canal se cierre y se vacíe. 💀
		for msg := range l.ch {
			_, _ = io.WriteString(l.writer, msg+"\n")
		}
	}()

	return l
}

// Log intenta enviar un mensaje a la cola de logs. 🚀
// Implementa el Drop Pattern: si el buffer está lleno, incrementa el contador de descartes
// y retorna inmediatamente sin bloquear la ejecución del llamador.
func (l *Logger) Log(msg string) {
	select {
	case l.ch <- msg:
		// El mensaje se encoló correctamente. ✅
	default:
		// ¡Pánico controlado! El buffer está lleno, así que descartamos el log. 📤
		// Usamos atomic para evitar condiciones de carrera sin usar Mutex.
		l.dropped.Add(1)
	}
}

// Shutdown cierra el sistema de logs de forma limpia. 🛑
// Primero cierra el canal para que la goroutine deje de recibir nuevos datos
// y luego espera a que termine de procesar los mensajes pendientes en el buffer.
func (l *Logger) Shutdown() {
	// Cerramos el canal: señal para que el 'range' termine.
	close(l.ch)

	// Esperamos a que el "drain" (drenaje) se complete.
	l.wg.Wait()

	/*
	   NOTA DIDÁCTICA SOBRE CANALES Y SHUTDOWN:
	   Si en lugar de un WaitGroup usáramos un canal de resultado para señalizar el fin,
	   ese canal DEBERÍA ser buffered con capacidad 1 (ej: chan struct{1}).
	   ¿Por qué? Para prevenir un "goroutine leak": si la goroutine principal que espera
	   se marcha antes (ej: por un timeout), la goroutine de log se quedaría bloqueada
	   intentando enviar al canal de resultado para siempre. El buffer de 1 permite
	   enviar y terminar aunque nadie esté escuchando al otro lado. 💡
	*/
}

// DroppedCount retorna cuántos mensajes se han perdido por el camino. 📊
func (l *Logger) DroppedCount() int64 {
	return l.dropped.Load()
}
