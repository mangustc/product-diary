package main

import (
	"net/http"

	"github.com/bmg-c/product-diary/db"
	"github.com/bmg-c/product-diary/db/user_db"
	"github.com/bmg-c/product-diary/handlers"
	"github.com/bmg-c/product-diary/logger"
	"github.com/bmg-c/product-diary/middleware"
	"github.com/bmg-c/product-diary/services"
	"github.com/bmg-c/product-diary/tests"
)

func main() {
	err := tests.TestValidation()
	if err != nil {
		logger.Error.Println(err.Error())
	} else {
		logger.Info.Println("Tests passed successfully")
	}

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
		panic(err.Error())
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
		panic(err.Error())
	} else {
		logger.Info.Println("Successfully connected code store")
	}
	sessionStore, err := db.NewStore("database.db", "sessions",
		`CREATE TABLE IF NOT EXISTS sessions (
        session_uuid VARCHAR(32) NOT NULL UNIQUE,
        user_id INTEGER NOT NULL,
        FOREIGN KEY (user_id) REFERENCES `+userStore.TableName+` (user_id) ON DELETE RESTRICT
    );`)
	if err != nil {
		logger.Error.Println("Error creating session store: " + err.Error())
		panic(err.Error())
	} else {
		logger.Info.Println("Successfully connected session store")
	}
	udb, err := user_db.NewUserDB(userStore, codeStore, sessionStore)
	if err != nil {
		logger.Error.Println("Error creating user database layer: " + err.Error())
	}
	us := services.NewUserService(udb)
	uh := handlers.NewUserHandler(us)
	router.HandleFunc("GET /users", uh.HandleUsersPage)
	router.HandleFunc("GET /api/users/controls/index", uh.HandleControlsIndex)
	router.HandleFunc("GET /api/users/userlist/index", uh.HandleGetUsersAll)
	router.HandleFunc("GET /api/users/user/index", uh.HandleUserIndex)
	router.HandleFunc("POST /api/users/user/getuser", uh.HandleGetUser)
	router.HandleFunc("GET /api/users/signin/index", uh.HandleSigninIndex)
	router.HandleFunc("POST /api/users/signin/signin", uh.HandleSigninSignin)
	router.HandleFunc("GET /api/users/login/index", uh.HandleLoginIndex)
	router.HandleFunc("POST /api/users/login/login", uh.HandleLoginLogin)
	router.HandleFunc("GET /api/users/profile/index", uh.HandleProfileIndex)
	router.HandleFunc("POST /api/users/logout/logout", uh.HandleLogout)

	mh := handlers.NewMainHandler()
	router.HandleFunc("GET /api/locale/index", mh.HandleLocale)
	router.HandleFunc("POST /api/locale/setlocale", mh.HandleSetLocale)

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
