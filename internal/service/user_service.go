package service

import (
	"fmt"
	"github.com/user/golang_basics/internal/domain"
)

// UserRepository define QUÉ necesita el service (el contrato). 📜
// Al Service no le importa si usas Postgres, MongoDB o un papel y boli.
type UserRepository interface {
	Save(user domain.User) error
	GetByID(id int) (domain.User, error)
}

// UserService: El cerebro de la operación. 🧠
// Aquí reside la "Lógica Pura". Orquesta lo que debe pasar sin ensuciarse con detalles técnicos.
type UserService struct {
	repo UserRepository // Inyección de dependencia: recibimos una interfaz.
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterUser(user domain.User) error {
	// 1. El Service VALIDA las reglas de negocio. ✅
	if user.Name == "" {
		return &domain.ValidationError{Field: "Name", Message: "name is required"}
	}
	
	// 2. El Service delega el almacenamiento al repositorio.
	return s.repo.Save(user)
}

func (s *UserService) GetUser(id int) (domain.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		// El Service ENVUELVE errores con %w. 🎁
		// Esto permite añadir contexto ("service failed") sin perder el error original (ErrNotFound).
		return domain.User{}, fmt.Errorf("service failed: %w", err)
	}
	return user, nil
}
