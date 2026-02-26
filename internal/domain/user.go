package domain

import "fmt"

// --- DECOUPLING Y VISIBILIDAD (EXPORTING) ---
// En Go, no hay public/private. Si empieza por Mayúscula, se ve fuera. Si no, es timidez absoluta. 🛑

// User representa un usuario. 
// Como empieza con 'U', es EXPORTADO. Todo el mundo puede usar este tipo.
type User struct {
	ID    int    // Público: Cualquiera puede leer su ID.
	Name  string // Público.
	Email string // Público.
	
	// passwordHash es PRIVADO (Unexported). 🔒
	// Solo el código dentro de la carpeta 'domain' puede tocar esto.
	// ¡Adiós a que main imprima contraseñas por error!
	passwordHash string 
}

// NewUser es nuestra Factory Function (el "Constructor"). 🏗️
// Es la forma oficial de crear un usuario. Nos asegura que el estado inicial es válido.
func NewUser(name, email, password string) *User {
	return &User{
		Name:         name,
		Email:        email,
		passwordHash: "hash_" + password, // Aplicamos lógica interna (hashing)
	}
}

// CheckPassword: Sigue el principio "Tell, Don't Ask". 🗣️
// No le pidas el hash al usuario para compararlo tú; dile que él se valide solo.
func (u *User) CheckPassword(input string) bool {
	return u.passwordHash == "hash_"+input
}

// UpdatePassword: Control total sobre los cambios.
// Nadie puede poner passwordHash a "" desde fuera. Tienen que pasar por aquí.
func (u *User) UpdatePassword(old, new string) error {
	if !u.CheckPassword(old) {
		return fmt.Errorf("invalid password")
	}
	u.passwordHash = "hash_" + new
	return nil
}

// Mover define comportamiento
type Mover interface {
	Move() string
}
