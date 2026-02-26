package education

import "fmt"

// --- CONSTANTS ---
// Las constantes viven solo en tiempo de compilación. 🚀

const (
	// iota: el "contador mágico" de Go para Enums.
	_     = iota            // 0: Saltamos el cero (frecuente para evitar el "zero-value")
	Lunes                   // 1: ¡No hace falta repetir iota! Go lo hace por ti.
	Martes                  // 2
	Miercoles               // 3
)

const (
	// Flags usando bitwise: Ideal para combinar permisos. 🚩
	Read  = 1 << iota // 1 (binario: 001)
	Write = 1 << iota // 2 (binario: 010)
	Exec  = 1 << iota // 4 (binario: 100)
)

// ShowConstants demuestra el uso de constantes typed vs untyped
func ShowConstants() {
	// Untyped Constants: Son como "comodines". No tienen tipo fijo hasta que se usan.
	// Esto permite usarlas con float32 o float64 sin quejarse.
	const Pi = 3.14 
	
	var f float64 = Pi  // ✅ Ok
	var f32 float32 = Pi // ✅ Ok, el compilador adapta Pi mágicamente
	
	fmt.Printf("Pi (f64): %v, Pi (f32): %v\n", f, f32)
	fmt.Printf("Lunes: %v, Martes: %v\n", Lunes, Martes)
	
	// Combinamos flags usando el operador OR (|)
	permisos := Read | Write
	// Verificamos si tiene el flag usando AND (&)
	fmt.Printf("Tiene Read: %v\n", permisos&Read != 0)
}
