package service_test

import (
	"testing"
	"github.com/user/golang_basics/internal/domain"
)

// 🤡 5. Consumer-Defined Mocking
// Tú (como consumidor del repositorio) defines la interface pequeña que necesitas.
type UserSaver interface {
	Save(user domain.User) error
}

// Mock simple
type MockRepo struct {
	saveCalled bool
}

func (m *MockRepo) Save(user domain.User) error {
	m.saveCalled = true
	return nil
}

// Nota: Para que UserService pueda usar este Mock, tendría que haber sido diseñado 
// usando una Interface en lugar de un Concrete Type. 
// Pero siguiendo la regla "Concrete First", UserService usa *repository.PostgresUserRepository.
// Para testing real con mocks, refactorizaríamos a Interface SOLO cuando sea necesario.

func TestRegisterUser(t *testing.T) {
    // Aquí iría la lógica de test
}
