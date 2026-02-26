package education

import (
	"fmt"
)

// ShowTypes muestra la esencia de los tipos básicos y sus alias.
// 💡 Tip: Go es de tipado fuerte, ¡no le gusta mezclar peras con manzanas!
func ShowTypes() {
	// --- Booleanos ---
	var b bool = true
	fmt.Printf("Bool: %v\n", b)

	// --- Numéricos ---
	// int y uint dependen de la arquitectura (32 o 64 bits)
	var i int = -42
	var u uint = 42
	var f float64 = 3.1415
	var c complex128 = complex(1, 2) // Para los matemáticos 🤓
	fmt.Printf("Int: %d, Uint: %d, Float: %f, Complex: %v\n", i, u, f, c)

	// --- Strings ---
	// ⚠️ Inmutables y secuencia de bytes UTF-8. 
	// No puedes cambiar un carácter de un string directamente.
	var s string = "Hola Go 🚀"
	fmt.Printf("String: %s\n", s)

	// --- Otros (Alias con nombre propio) ---
	var by byte = 255  // 🎭 Alias de uint8 (ideal para datos binarios)
	var r rune = 'A'   // 🎭 Alias de int32 (representa un punto de código Unicode)
	fmt.Printf("Byte: %d, Rune: %U\n", by, r)
}
