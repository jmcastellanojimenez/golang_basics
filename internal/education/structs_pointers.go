package education

import (
	"fmt"
)

// --- STRUCTS Y PADDING ---
// El orden de los factores SÍ altera el producto (en memoria). 📏

// BadOrder: Ocupa 24 bytes. 
// Go alinea los datos. El bool (1 byte) deja 7 vacíos para que el int64 empiece en un múltiplo de 8.
type BadOrder struct {
	A bool  // 1 byte + 7 padding
	B int64 // 8 bytes
	C bool  // 1 byte + 7 padding
}

// GoodOrder: Ocupa solo 16 bytes.
// Al agrupar los pequeños al final, ahorramos espacio. ¡Eficiencia pura!
type GoodOrder struct {
	B int64 // 8 bytes
	A bool  // 1 byte
	C bool  // 1 byte (+ 6 padding al final)
}

// --- POINTERS: VALUE vs POINTER SEMANTICS ---

type Person struct {
	Name string
}

// Value semantics: Recibe una COPIA. 
// Lo que pase dentro de la función, se queda en la función. 🔒
func IncrementAgeValue(age int) {
	age++
}

// Pointer semantics: Recibe la DIRECCIÓN (puntero).
// Modifica el dato original. ¡Cuidado con los efectos secundarios! 💉
func IncrementAgePointer(age *int) {
	*age++ // Dereferencing: accedemos al valor que hay en la dirección
}

// --- FACTORY FUNCTIONS & ESCAPE ANALYSIS ---

// NewPersonValue: Se queda en el STACK (rápido). 🏔️
// Crea una persona, la copia y la devuelve. La original muere con la función.
func NewPersonValue(name string) Person {
	return Person{Name: name}
}

// NewPersonPointer: Escapa al HEAP (más lento). 🏨
// Retornas la dirección de algo local. Go se da cuenta y dice: 
// "Oye, esto tiene que sobrevivir fuera de aquí, lo muevo al Heap".
func NewPersonPointer(name string) *Person {
	return &Person{Name: name} 
}

// --- SHARING UP vs DOWN ---

// sharing DOWN: Pasas un puntero hacia abajo. ⬇️
// Es seguro. Main (el dueño) sigue vivo mientras esta función trabaja. Se queda en Stack.
func shareDown(p *Person) {
	p.Name = "Updated Down"
}

// sharing UP: Retornas un puntero hacia arriba. ⬆️
// Peligro: El dato debe sobrevivir a su creador. Escapa al Heap.
func shareUp() *Person {
	p := Person{Name: "Created Up"}
	return &p 
}

// --- PASS-BY-VALUE CON POINTERS (LA TRAMPA) ---
// En Go, TODO se pasa por valor, ¡incluso los punteros!

func ChangePointer(p *Person) {
	// ❌ ERROR: Estás reasignando tu COPIA local del puntero.
	// El 'p' de main sigue apuntando al sitio original.
	p = &Person{Name: "New Person"} 
}

func ChangeValue(p *Person) {
	// ✅ OK: Estás usando tu copia del puntero para ir a la dirección compartida
	// y cambiar lo que hay allí.
	p.Name = "Changed Value" 
}

func ShowPointers() {
	x := 10
	IncrementAgeValue(x)
	fmt.Printf("Value (no cambia): %d\n", x)
	
	IncrementAgePointer(&x)
	fmt.Printf("Pointer (sí cambia): %d\n", x)
	
	p := &Person{Name: "Original"}
	ChangePointer(p)
	fmt.Printf("After ChangePointer: %s\n", p.Name)
	
	ChangeValue(p)
	fmt.Printf("After ChangeValue: %s\n", p.Name)
}
