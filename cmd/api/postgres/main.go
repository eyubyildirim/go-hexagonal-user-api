package main

import (
	"context"
	httpadapter "hexa-user/internal/adapters/http"
	"hexa-user/internal/adapters/storage"
	"hexa-user/internal/domain"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://eyub:1234@localhost:5432/user_service_db?sslmode=disable"
	}
	log.Printf("Using database URL: %s", databaseURL)

	dbPool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer dbPool.Close()

	userRepo, err := storage.NewPostgresUserRepository(dbPool)
	if err != nil {
		log.Fatalf("Failed to create user repository: %v", err)
	}

	if err := userRepo.CreateTableIfNotExists(ctx); err != nil {
		log.Fatalf("Failed to ensure users table exists: %v", err)
	}

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
	log.Println("Server stopped gracefully")
}
