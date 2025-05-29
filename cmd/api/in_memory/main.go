package main

import (
	httpadapter "hexa-user/internal/adapters/http"
	"hexa-user/internal/adapters/storage"
	"hexa-user/internal/domain"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	userRepo := storage.NewInMemoryUserRepository()
	userService := domain.NewUserService(userRepo)
	userHandler := httpadapter.NewUserHandler(userService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", userHandler.CreateUserHandler)
	mux.HandleFunc("GET /users", userHandler.ListUsersHandler)
	mux.HandleFunc("GET /users/{id}", userHandler.GetUserByIDHandler)
	mux.HandleFunc("PUT /users/{id}", userHandler.UpdateUserHandler)
	mux.HandleFunc("DELETE /users/{id}", userHandler.DeleteUserHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Printf("Server starting on port %s\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
