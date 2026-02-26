package repository

import (
	"github.com/user/golang_basics/internal/domain"
)

// PostgresUserRepository es un Concrete Type. 🏗️
// En Go, solemos escribir primero el tipo concreto y luego extraer interfaces si es necesario.
// Rule: "Concrete First".
type PostgresUserRepository struct {
	users map[int]domain.User // Simulación de base de datos
}

func NewPostgresUserRepository() *PostgresUserRepository {
	return &PostgresUserRepository{
		users: make(map[int]domain.User),
	}
}

// Save implementa el CÓMO se guardan los datos. 💾
func (r *PostgresUserRepository) Save(user domain.User) error {
	r.users[user.ID] = user
	return nil
}

func (r *PostgresUserRepository) GetByID(id int) (domain.User, error) {
	user, ok := r.users[id]
	if !ok {
		// El repositorio es quien sabe que algo no existe. 
		// Por eso GENERA el error Sentinel de dominio. 🚩
		return domain.User{}, domain.ErrNotFound
	}
	return user, nil
}
