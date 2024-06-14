package main

import (
	"fmt"
	"net/http"

	"github.com/bmg-c/product-diary/db"
	"github.com/bmg-c/product-diary/handlers"
	"github.com/bmg-c/product-diary/services"
	// "github.com/bmg-c/product-diary/views/test_views"
)

func main() {
	router := http.NewServeMux()

	userStore, err := db.NewStore("database.db", "users",
        `CREATE TABLE IF NOT EXISTS users (
        user_id INTEGER PRIMARY KEY AUTOINCREMENT,
        username VARCHAR(64) NOT NULL,
        email VARCHAR(255) NOT NULL UNIQUE,
        password VARCHAR(255) NOT NULL,
        created_at DATETIME default CURRENT_TIMESTAMP
    );`)
	if err != nil {
		fmt.Printf("user store: %s\n", err.Error())
	}
	codeStore, err := db.NewStore("database.db", "codes",
        `CREATE TABLE IF NOT EXISTS codes (
        code_id INTEGER PRIMARY KEY AUTOINCREMENT,
        email VARCHAR(255) NOT NULL UNIQUE,
        code VARCHAR(6),
        created_at DATETIME default CURRENT_TIMESTAMP
    );`)
	if err != nil {
		fmt.Printf("code store: %s\n", err.Error())
	}
	us := services.NewUserService(services.UserPublic{}, userStore, codeStore)
	uh := handlers.NewUserHandler(us)
	router.HandleFunc("GET /users/{id}", uh.HandleGetUserByID)
	router.HandleFunc("GET /users", uh.HandleGetUsersAll)
	router.HandleFunc("POST /users/register", uh.HandleRegisterUser)
	router.HandleFunc("POST /users/confirmregister", uh.HandleConfirmRegister)

	port := ":1323"
	server := http.Server{
		Addr:    port,
		Handler: router,
	}

	fmt.Println("Server is listening on port " + port)
	server.ListenAndServe()
}
