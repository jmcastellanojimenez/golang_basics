package domain

import "errors"

// --- LOS 3 TIPOS DE ERROR ---
// En Go, los errores son valores. No hay excepciones que rompan el flujo. 🛑

// 1. Sentinel Error (Estado/Valor) 🚩
// Es una variable global. Se usa para comparaciones rápidas: "Is this error EXACTLY this one?".
var ErrNotFound = errors.New("user not found")

// 2. Type Error (Estructura/Contexto) 🏷️
// Un struct que implementa Error(). Útil para dar más info (qué campo falló).
// Se captura con errors.As().
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// 3. Behavior Error (Interfaz/Capacidad) 🔄
// No nos importa el tipo, solo si se comporta de cierta forma. 
// Por ejemplo: "¿Es este error temporal? ¿Puedo reintentar?".
type Temporary interface {
	IsTemporary() bool
}
