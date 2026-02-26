package main

import (
	"fmt"
	"net/http"
	"github.com/user/golang_basics/internal/handler"
	"github.com/user/golang_basics/internal/repository"
	"github.com/user/golang_basics/internal/service"
)

func main() {
	// 1. El Orquestador: Aquí es donde ocurre la magia del "cableado". 🎼
	// Creamos las piezas concretas y las inyectamos.
	// Fíjate cómo el repo (concreto) se le pasa al service (que pide una interfaz).
	
	repo := repository.NewPostgresUserRepository()
	svc := service.NewUserService(repo) // Inyección de dependencia (Rule: main CABLEA todo)
	hdl := handler.NewUserHandler(svc)

	// 2. Setup Routes: Registramos los puntos de entrada. 🛣️
	http.HandleFunc("/users", hdl.CreateUser)
	http.HandleFunc("/user", hdl.GetUser)

	// 3. Start Server: ¡Encendemos los motores! 🚀
	fmt.Println("🚀 Server running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
