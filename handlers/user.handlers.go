package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bmg-c/product-diary/services"
)

type UserService interface {
	GetUserByID(id int) (services.UserPublic, error)
	GetUsersAll() ([]services.UserPublic, error)
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
	users, err := uh.UserService.GetUsersAll()
	if err != nil {
		fmt.Printf("user handler: Failed to get users from the database (%s)", err.Error())
	}
	out, err := json.MarshalIndent(users, "", " ")
	if err != nil {
		fmt.Printf("user handler: Failed to convert users to json (%s)", err.Error())
	}
	w.Write(out)
}

func (uh *UserHandler) HandleGetUserByID(w http.ResponseWriter, r *http.Request) {
	user, err := uh.UserService.GetUserByID(0)
	if err != nil {
		fmt.Printf("user handler: Failed to get users from the database (%s)", err)
	}
	out, err := json.MarshalIndent(user, "", " ")
	if err != nil {
		fmt.Printf("user handler: Failed to convert users to json (%s)", err.Error())
	}
	w.Write(out)
}
