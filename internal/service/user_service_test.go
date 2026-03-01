package service_test

import (
	"errors"
	"testing"

	"github.com/user/golang_basics/internal/domain"
	"github.com/user/golang_basics/internal/service"
)

// 🤡 MockRepo: Un espía para nuestros tests.
// Implementa service.UserRepository (la interfaz que el service necesita).
// No guarda nada en disco, solo en memoria para verificar qué pasó.
type MockRepo struct {
	users map[int]domain.User
	err   error // Forzamos un error si queremos testear fallos.
}

func NewMockRepo() *MockRepo {
	return &MockRepo{users: make(map[int]domain.User)}
}

func (m *MockRepo) Save(user domain.User) error {
	if m.err != nil {
		return m.err
	}
	m.users[user.ID] = user
	return nil
}

func (m *MockRepo) GetByID(id int) (domain.User, error) {
	if m.err != nil {
		return domain.User{}, m.err
	}
	user, ok := m.users[id]
	if !ok {
		return domain.User{}, domain.ErrNotFound
	}
	return user, nil
}

// 🧪 TestRegisterUser_TableDriven: El estándar de oro en Go.
// Probamos muchos casos de una sola vez usando una "tabla" de datos.
func TestRegisterUser_TableDriven(t *testing.T) {
	tests := []struct {
		name    string
		user    domain.User
		mockErr error
		wantErr bool
	}{
		{
			name:    "Success: Valid User",
			user:    domain.User{ID: 1, Name: "Gopher", Email: "go@lang.org"},
			mockErr: nil,
			wantErr: false,
		},
		{
			name:    "Fail: Missing Name",
			user:    domain.User{ID: 2, Name: "", Email: "no-name@lang.org"},
			mockErr: nil,
			wantErr: true,
		},
		{
			name:    "Fail: Repository Error",
			user:    domain.User{ID: 3, Name: "Broken", Email: "broken@repo.com"},
			mockErr: errors.New("db connection failed"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		// t.Run crea un "sub-test" para cada fila de la tabla. 📉
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepo()
			repo.err = tt.mockErr
			svc := service.NewUserService(repo)

			err := svc.RegisterUser(tt.user)

			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Si no queríamos error, verificamos que el repo realmente tiene al usuario. ✅
			if !tt.wantErr && repo.users[tt.user.ID].Name != tt.user.Name {
				t.Errorf("User was not correctly saved in repo")
			}
		})
	}
}

// 🧪 TestGetUser: Verificando la recuperación y el "wrapping" de errores.
func TestGetUser(t *testing.T) {
	repo := NewMockRepo()
	svc := service.NewUserService(repo)

	// Pre-insertamos un usuario
	repo.Save(domain.User{ID: 42, Name: "Douglas"})

	t.Run("Found", func(t *testing.T) {
		user, err := svc.GetUser(42)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if user.Name != "Douglas" {
			t.Errorf("Expected Douglas, got %s", user.Name)
		}
	})

	t.Run("NotFound_WrappedError", func(t *testing.T) {
		_, err := svc.GetUser(99)
		if err == nil {
			t.Fatal("Expected error for missing user")
		}

		// Verificamos que el error está "envuelto" (wrapped) y contiene el original. 🎁
		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("Expected error to wrap ErrNotFound, but got: %v", err)
		}
	})
}
