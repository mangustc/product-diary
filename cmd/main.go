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

	uStore, err := db.NewStore("database.db", "users",
		`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username VARCHAR(64) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		created_at DATETIME default CURRENT_TIMESTAMP);`)
	if err != nil {
		fmt.Printf("user store: %s", err.Error())
	}
	us := services.NewUserService(services.UserPublic{}, uStore)
	uh := handlers.NewUserHandler(us)
	router.HandleFunc("GET /users/byid", uh.HandleGetUserByID)
	router.HandleFunc("GET /users", uh.HandleGetUsersAll)

	port := ":1323"
	server := http.Server{
		Addr:    port,
		Handler: router,
	}

	fmt.Println("Server is listening on port " + port)
	server.ListenAndServe()
}
