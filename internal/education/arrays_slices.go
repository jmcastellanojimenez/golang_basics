package education

import (
	"fmt"
)

// --- ARRAYS vs SLICES ---

func ShowArraysSlices() {
	// Array: Tamaño fijo. ¡Ojo! Al pasar un array a una función, Go COPIA todo.
	// Si el array tiene 1 millón de elementos, copiarás 1 millón de elementos. 🐌
	arr := [3]int{10, 20, 30}
	fmt.Printf("Array: %v, Len: %d\n", arr, len(arr))
	
	// Slice: Es un "mando a distancia" (ptr, len, cap). Pesa solo 24 bytes.
	// No importa si apunta a un array de 10 elementos o de 10 millones. ⚡
	sli := []int{10, 20, 30}
	
	// PILAR 2: Slicing comparte memoria. 🔪
	// 'vista' es una ventana sobre el array de 'sli'. ¡Si tocas uno, tocas el otro!
	vista := sli[1:3]
	vista[0] = 999
	fmt.Printf("Slice original modificado por vista: %v\n", sli)

	// PILAR 3: append y el truco de la Capacidad. 🪄
	// Si hay hueco (Len < Cap), append escribe en el mismo array.
	// Si NO hay hueco, Go crea un array nuevo y se muda. El vínculo se rompe.
	s1 := make([]int, 0, 2)
	s1 = append(s1, 1, 2) // Len=2, Cap=2
	s2 := s1
	s1 = append(s1, 3) // s1 se muda a un array nuevo más grande.
	s1[0] = 888       // s2 ni se entera del cambio.
	fmt.Printf("s1 (mudado): %v, s2 (viejo): %v\n", s1, s2)
}

// STRINGS: Son inmutables. 🎭
// Internamente son un struct con un puntero y un largo. 
// Al copiar un string, solo copias el struct, ¡no el texto completo!
func ShowStrings() {
	s1 := "hello"
	s2 := s1 
	fmt.Printf("String copia comparte array: %s, %s\n", s1, s2)
}

// FOR RANGE: ¿Copia o Índice? 🔄
func RangeSemantics() {
	nums := [3]int{1, 2, 3}
	
	// Value semantic: 'v' es una COPIA del elemento.
	// Útil para leer sin efectos secundarios.
	for i, v := range nums {
		v = 999 
		fmt.Printf("Copia v[%d]: %d (nums[%d] sigue igual)\n", i, v, i)
	}
	
	// Pointer semantic: Usamos solo el índice 'i'.
	// Accedemos directamente a la memoria del array original.
	for i := range nums {
		nums[i] = 999 
	}
}
