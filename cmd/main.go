package main

import (
	// "database/sql"
	"net/http"

	"github.com/bmg-c/product-diary/db"
	"github.com/bmg-c/product-diary/db/item_db"
	"github.com/bmg-c/product-diary/db/product_db"
	"github.com/bmg-c/product-diary/db/user_db"
	"github.com/bmg-c/product-diary/handlers"
	"github.com/bmg-c/product-diary/logger"
	"github.com/bmg-c/product-diary/middleware"
	"github.com/bmg-c/product-diary/services"
	"github.com/bmg-c/product-diary/tests"
	// "github.com/mattn/go-sqlite3"
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
	personStore, err := db.NewStore("database.db", "persons",
		`CREATE TABLE IF NOT EXISTS persons (
        person_id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        person_name VARCHAR(64) NOT NULL,
        is_hidden INTEGER NOT NULL DEFAULT FALSE,
        FOREIGN KEY (user_id) REFERENCES `+userStore.TableName+` (user_id) ON DELETE RESTRICT,
        UNIQUE(user_id, person_name)
    );`)
	if err != nil {
		logger.Error.Println("Error creating person store: " + err.Error())
		panic(err.Error())
	} else {
		logger.Info.Println("Successfully connected person store")
	}
	udb, err := user_db.NewUserDB(userStore, codeStore, sessionStore, personStore)
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
	router.HandleFunc("POST /api/users/person/togglehidden", uh.HandleTogglePerson)
	router.HandleFunc("POST /api/users/person/addperson", uh.HandleAddPerson)

	productStore, err := db.NewStore("database.db", "products",
		`CREATE TABLE IF NOT EXISTS products (
        product_id INTEGER PRIMARY KEY AUTOINCREMENT,
        product_title VARCHAR(128) NOT NULL,
        product_calories INTEGER DEFAULT 0,
        product_fats INTEGER DEFAULT 0,
        product_carbs INTEGER DEFAULT 0,
        product_proteins INTEGER DEFAULT 0,
        user_id INTEGER NOT NULL,
        is_deleted INTEGER NOT NULL DEFAULT FALSE,
        CHECK (product_fats + product_carbs + product_proteins <= 100),
        CHECK (length(product_title) >= 4 AND length(product_title) <= 128),
        FOREIGN KEY (user_id) REFERENCES `+userStore.TableName+` (user_id) ON DELETE RESTRICT
    );`)
	if err != nil {
		logger.Error.Println("Error creating product store: " + err.Error())
		panic(err.Error())
	} else {
		logger.Info.Println("Successfully connected product store")
	}
	pdb, err := product_db.NewProductDB(productStore)
	if err != nil {
		logger.Error.Println("Error creating product database layer: " + err.Error())
	}
	ps := services.NewProductService(pdb)
	ph := handlers.NewProductHandler(ps, us)
	router.HandleFunc("GET /products", ph.HandleProductsPage)
	router.HandleFunc("POST /api/products/addproduct", ph.HandleAddProduct)
	router.HandleFunc("POST /api/products/getproducts", ph.HandleGetProducts)
	router.HandleFunc("POST /api/products/copyproduct", ph.HandleCopyProduct)
	router.HandleFunc("POST /api/products/deleteproduct", ph.HandleDeleteProduct)

	itemStore, err := db.NewStore("database.db", "items",
		`CREATE TABLE IF NOT EXISTS items (
        item_id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        product_id INTEGER NOT NULL,
        item_date DATE NOT NULL,
        item_cost REAL DEFAULT 0,
        item_amount REAL DEFAULT 0,
        item_type INTEGER NOT NULL DEFAULT 1,
        person_id INTEGER DEFAULT NULL,
        CHECK (item_type >= 1 AND item_type <= 3),
        CHECK (item_cost >= 0),
        CHECK (item_amount >= 0),
        FOREIGN KEY (user_id) REFERENCES `+userStore.TableName+` (user_id) ON DELETE RESTRICT,
        FOREIGN KEY (product_id) REFERENCES `+productStore.TableName+` (product_id) ON DELETE RESTRICT,
        FOREIGN KEY (person_id) REFERENCES `+personStore.TableName+` (person_id) ON DELETE RESTRICT
    );`)
	if err != nil {
		logger.Error.Println("Error creating item store: " + err.Error())
		panic(err.Error())
	} else {
		logger.Info.Println("Successfully connected item store")
	}
	idb, err := item_db.NewItemDB(itemStore, productStore, personStore)
	if err != nil {
		logger.Error.Println("Error creating item database layer: " + err.Error())
	}
	is := services.NewItemService(idb)
	ih := handlers.NewItemHandler(is, us)
	router.HandleFunc("GET /analytics", ih.HandleAnalyticsPage)
	router.HandleFunc("POST /api/items/getitems", ih.HandleGetItems)
	router.HandleFunc("POST /api/items/additem", ih.HandleAddItem)
	router.HandleFunc("POST /api/items/deleteitem", ih.HandleDeleteItem)
	router.HandleFunc("POST /api/items/changeitem", ih.HandleChangeItem)
	router.HandleFunc("POST /api/items/getanalyticsrange", ih.HandleGetAnalyticsRange)

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
