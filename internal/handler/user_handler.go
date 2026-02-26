package handler

import (
	"encoding/json"
	"errors"
	"github.com/user/golang_basics/internal/domain"
	"github.com/user/golang_basics/internal/service"
	"net/http"
	"strconv"
)

// UserHandler adapta el mundo exterior (HTTP) a nuestro Service. ⚗️
// Se encarga de traducir JSON a structs y, lo más importante, DESENVOLVER errores.
type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.RegisterUser(user); err != nil {
		// El Handler DESENVUELVE el error para saber qué código HTTP enviar. 🔍
		
		// 🏷️ Caso 2: Error de TIPO (ValidationError)
		var valErr *domain.ValidationError
		if errors.As(err, &valErr) { 
			http.Error(w, valErr.Message, http.StatusBadRequest)
			return
		}
		
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	user, err := h.service.GetUser(id)
	if err != nil {
		// 🚩 Caso 1: Error SENTINEL (ErrNotFound)
		// Usamos errors.Is para buscar la señal específica en la cadena de errores.
		if errors.Is(err, domain.ErrNotFound) { 
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		// 🔄 Caso 3: Error de COMPORTAMIENTO (Temporary)
		// No sabemos qué error es, pero le preguntamos si es temporal.
		var tempErr domain.Temporary
		if errors.As(err, &tempErr) && tempErr.IsTemporary() {
			http.Error(w, "Try again later", http.StatusServiceUnavailable)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}
