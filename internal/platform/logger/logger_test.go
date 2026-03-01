package logger

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"
)

// TestLoggerNormalOperation verifica que en condiciones ideales no se pierda ningún log. ✅
func TestLoggerNormalOperation(t *testing.T) {
	var buf bytes.Buffer
	n := 50
	l := New(&buf, 64)

	for i := 0; i < n; i++ {
		l.Log(fmt.Sprintf("msg %d", i))
	}

	l.Shutdown()

	if l.DroppedCount() != 0 {
		t.Errorf("esperaba 0 mensajes descartados, obtuvo %d", l.DroppedCount())
	}

	// Contamos cuántas líneas hay en el buffer.
	got := strings.Count(buf.String(), "\n")
	if got != n {
		t.Errorf("esperaba %d mensajes en el buffer, obtuvo %d", n, got)
	}
}

// slowWriter simula un disco o red muy lentos. 🐢
type slowWriter struct{}

func (s *slowWriter) Write(p []byte) (n int, err error) {
	time.Sleep(50 * time.Millisecond)
	return len(p), nil
}

// TestLoggerSlowWriter asegura que el Logger NUNCA bloquee al llamador,
// incluso si el destino es lento. Esto es el Drop Pattern en acción. 🛡️
func TestLoggerSlowWriter(t *testing.T) {
	sw := &slowWriter{}
	l := New(sw, 8) // Buffer pequeño para forzar el descarte rápido.
	n := 200

	for i := 0; i < n; i++ {
		start := time.Now()
		l.Log("un mensaje cualquiera")
		duration := time.Since(start)

		// El requisito es que Log() no tarde más de 1ms.
		if duration > time.Millisecond {
			t.Errorf("¡BLOQUEO detectado! Log() tardó %v, máximo permitido 1ms", duration)
		}
	}

	l.Shutdown()

	// Como el writer tarda 50ms por log y el buffer es de 8, 
	// la mayoría de los 200 logs deben haberse descartado.
	if l.DroppedCount() == 0 {
		t.Error("esperaba que se descartaran mensajes por lentitud del writer, pero DroppedCount es 0")
	}
}

// TestLoggerShutdownDrains verifica que al cerrar el logger, se procese
// todo lo que ya estaba en el buffer. No se tira nada que ya haya entrado. 🚰
func TestLoggerShutdownDrains(t *testing.T) {
	var buf bytes.Buffer
	bufferSize := 128
	l := New(&buf, bufferSize)

	// Llenamos el buffer al máximo.
	for i := 0; i < bufferSize; i++ {
		l.Log("drenaje")
	}

	// Shutdown debería esperar a que esos 128 mensajes se escriban.
	l.Shutdown()

	if l.DroppedCount() != 0 {
		t.Errorf("no se deberían haber descartado mensajes, obtuvo %d", l.DroppedCount())
	}

	got := strings.Count(buf.String(), "\n")
	if got != bufferSize {
		t.Errorf("esperaba %d mensajes drenados, obtuvo %d", bufferSize, got)
	}
}
