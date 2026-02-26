package education

import "fmt"

// --- EMBEDDING (COMPOSICIÓN) ---
// Go no tiene herencia (extends). Tiene piezas de Lego. 🧩

type Engine struct {
	Horsepower int
}

func (e Engine) Start() {
	fmt.Println("Vroom! Engine started.")
}

type Car struct {
	Engine // Embedding: Al no ponerle nombre, "promocionamos" sus tripas.
	Brand  string
}

// SHADOWING: El struct de fuera "tapa" al de dentro si tienen el mismo nombre. 👤
func (c Car) Start() {
	fmt.Printf("%s starting... Quietly.\n", c.Brand)
}

func ShowEmbedding() {
	c := Car{
		Brand:  "Tesla",
		Engine: Engine{Horsepower: 300},
	}
	
	// Magia: Accedemos directamente a Horsepower aunque sea de Engine.
	fmt.Printf("HP: %d\n", c.Horsepower) 

	c.Start()        // Llama a Car.Start() (el externo gana).
	c.Engine.Start() // El original no ha muerto, está ahí dentro si lo necesitas.
}

// TRAMPA: "TIENE UN" vs "ES UN". ⚠️
// Un Car TIENE un Engine, pero NO ES un Engine. 
// Go es estricto: la composición no es polimorfismo de tipos.
func Repair(e Engine) {
	fmt.Println("Repairing engine...")
}

func TestRepair() {
	c := Car{}
	// Repair(c) // ❌ ERROR: Car no es de tipo Engine.
	Repair(c.Engine) // ✅ Correcto: pasamos la pieza interna.
}
