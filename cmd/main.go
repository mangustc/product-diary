package main

import (
	"fmt"
	"net/http"

	"github.com/bmg-c/product-diary/handlers"
	"github.com/bmg-c/product-diary/services"
	// "github.com/bmg-c/product-diary/views/test_views"
)


func main() {
    router := http.NewServeMux()

    us := services.NewUserService(services.UserPublic{})
    uh := handlers.NewUserHandler(us)
    router.HandleFunc("GET /users/byid", uh.HandleGetUserByID)
    router.HandleFunc("GET /users", uh.HandleGetUsersAll)

    port := ":1323"
    server := http.Server {
        Addr: port,
        Handler: router,
    }

    fmt.Println("Server is listening on port " + port)
    server.ListenAndServe()
}
