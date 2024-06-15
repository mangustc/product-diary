package main

import (
	"net/http"

	"github.com/bmg-c/product-diary/db"
	"github.com/bmg-c/product-diary/handlers"
	"github.com/bmg-c/product-diary/logger"
	"github.com/bmg-c/product-diary/middleware"
	"github.com/bmg-c/product-diary/services"
)

func main() {
	router := http.NewServeMux()

	userStore, err := db.NewStore("database.db", "users",
		`CREATE TABLE IF NOT EXISTS users (
        user_id INTEGER PRIMARY KEY AUTOINCREMENT,
        username VARCHAR(64) NOT NULL,
        email VARCHAR(255) NOT NULL UNIQUE,
        password VARCHAR(255) NOT NULL,
        created_at DATETIME default (datetime('now'))
    );`)
	if err != nil {
		logger.Error.Println("Error creating user store: " + err.Error())
	} else {
		logger.Info.Println("Successfully connected user store")
	}
	codeStore, err := db.NewStore("database.db", "codes",
		`CREATE TABLE IF NOT EXISTS codes (
        code_id INTEGER PRIMARY KEY AUTOINCREMENT,
        email VARCHAR(255) NOT NULL UNIQUE,
        code VARCHAR(6),
        created_at DATETIME default (datetime('now'))
    );`)
	if err != nil {
		logger.Error.Println("Error creating code store: " + err.Error())
	} else {
		logger.Info.Println("Successfully connected code store")
	}
	us := services.NewUserService(services.UserPublic{}, userStore, codeStore)
	uh := handlers.NewUserHandler(us)
	router.HandleFunc("GET /users/{id}", uh.HandleGetUserByID)
	router.HandleFunc("GET /users", uh.HandleGetUsersAll)
	router.HandleFunc("POST /users/register", uh.HandleRegisterUser)
	router.HandleFunc("POST /users/confirmregister", uh.HandleConfirmRegister)

	port := ":1323"
	middlewareStack := middleware.CreateStack(
		middleware.Logging,
		middleware.StripSlash,
	)
	server := http.Server{
		Addr:    port,
		Handler: middlewareStack(router),
	}

	logger.Info.Println("Server is listening on port " + port)
	server.ListenAndServe()
}
