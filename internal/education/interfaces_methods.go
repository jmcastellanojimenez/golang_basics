package education

import (
	"fmt"
)

// --- METHODS: VALUE vs POINTER RECEIVER ---

type Counter struct {
	Count int
}

// Pointer receiver: ¡Modifica el original! 💉
// Go se encarga de coger la dirección si llamas desde un valor (c.Increment()).
func (c *Counter) Increment() {
	c.Count++
}

// Value receiver: Recibe una COPIA. 🔒
// Úsalo cuando no quieras mutar nada o para tipos inmutables.
func (c Counter) GetValue() int {
	return c.Count
}

// --- INTERFACES: EL CONTRATO (DUCK TYPING) ---
// En Go no dices "implements". Si tienes los métodos, YA ERES la interfaz. 🦆

type Speaker interface {
	Speak() string
}

type Dog struct {
	Name string
}

// Dog implementa Speaker implícitamente. ¡Magia!
func (d Dog) Speak() string {
	return "Guau!"
}

// --- METHOD SETS: LA REGLA DE ORO ---

type Updater interface {
	Update()
}

type Item struct {
	Val int
}

// Si el método usa un PUNTERO...
func (i *Item) Update() {
	i.Val++
}

func InterfaceRules() {
	it := Item{Val: 0}
	
	// ... la interfaz REQUIERE un puntero. ❌ var u Updater = it (FALLARÍA)
	var u Updater = &it // ✅ CORRECTO
	u.Update()
}

// --- ANY (interface{}) Y TYPE ASSERTION ---

// any es la caja donde cabe todo. Pero Go pierde el rastro del tipo real.
func InspectAny(v any) {
	// Type Switch: La forma elegante de abrir la caja. 📦
	switch val := v.(type) {
	case int:
		fmt.Printf("Es int: %d\n", val)
	case string:
		fmt.Printf("Es string: %s\n", val)
	default:
		fmt.Println("Tipo desconocido")
	}
	
	// Comma-OK idiom: La forma segura de forzar el tipo. 🛡️
	// Si no usas el 'ok', y no es del tipo, Go lanzará un PANIC.
	s, ok := v.(string)
	if ok {
		fmt.Printf("Aserción segura: %s\n", s)
	}
}
