package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bmg-c/product-diary/services"
)

type UserService interface {
	GetUserByID(id int)(services.UserPublic)
	GetUsersAll() ([]services.UserPublic)
}

func NewUserHandler(us UserService) *UserHandler {
	return &UserHandler{
		UserService: us,
	}
}

type UserHandler struct {
	UserService UserService
}

func (uh *UserHandler) HandleGetUsersAll(w http.ResponseWriter, r *http.Request) {
    users := uh.UserService.GetUsersAll()
    out, err := json.MarshalIndent(users, "", " ")
    if err != nil {
        fmt.Println("error at handle get users all")
    }
    w.Write(out)
}

func (uh *UserHandler) HandleGetUserByID(w http.ResponseWriter, r *http.Request) {
    user := uh.UserService.GetUserByID(0)
    out, err := json.MarshalIndent(user, "", " ")
    if err != nil {
        fmt.Println("error at handle get user by id")
    }
    w.Write(out)
}
