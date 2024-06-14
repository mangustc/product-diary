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
	RegisterUser(ur services.UserRegister) error
	ConfirmRegister(ucr services.UserConfirmRegister) error
}

func NewUserHandler(us UserService) *UserHandler {
	return &UserHandler{
		UserService: us,
	}
}

type UserHandler struct {
	UserService UserService
}

func (uh *UserHandler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var ur services.UserRegister
	err := json.NewDecoder(r.Body).Decode(&ur)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = uh.UserService.RegisterUser(ur)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte("Successfully sent confirmation code to " + ur.Email))
}

func (uh *UserHandler) HandleConfirmRegister(w http.ResponseWriter, r *http.Request) {
	var ucr services.UserConfirmRegister
	err := json.NewDecoder(r.Body).Decode(&ucr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = uh.UserService.ConfirmRegister(ucr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte("User login info was sent to " + ucr.Email))
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
		fmt.Printf("user handler: Failed to get users from the database (%s)", err.Error())
	}
	out, err := json.MarshalIndent(user, "", " ")
	if err != nil {
		fmt.Printf("user handler: Failed to convert users to json (%s)", err.Error())
	}
	w.Write(out)
}
