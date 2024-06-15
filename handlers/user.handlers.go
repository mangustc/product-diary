package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bmg-c/product-diary/errorhandler"
	"github.com/bmg-c/product-diary/logger"
	"github.com/bmg-c/product-diary/services"
	"github.com/bmg-c/product-diary/views"
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
		code := errorhandler.GetStatusCode(err)
		if code >= 500 {
			logger.Error.Println("Failed to decode request body: " + err.Error())
		}
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "text/html")
		views.ErrorIndex(code, http.StatusText(code)).Render(r.Context(), w)
		return
	}
	err = uh.UserService.RegisterUser(ur)
	if err != nil {
		code := errorhandler.GetStatusCode(err)
		if code >= 500 {
			logger.Error.Println("Failed to register user: " + err.Error())
		}
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "text/html")
		views.ErrorIndex(code, http.StatusText(code)).Render(r.Context(), w)
		return
	}
	w.Write([]byte("Successfully sent confirmation code to " + ur.Email))
}

func (uh *UserHandler) HandleConfirmRegister(w http.ResponseWriter, r *http.Request) {
	var ucr services.UserConfirmRegister
	err := json.NewDecoder(r.Body).Decode(&ucr)
	if err != nil {
		code := errorhandler.GetStatusCode(err)
		if code >= 500 {
			logger.Error.Println("Failed to decode request body: " + err.Error())
		}
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "text/html")
		views.ErrorIndex(code, http.StatusText(code)).Render(r.Context(), w)
		return
	}
	err = uh.UserService.ConfirmRegister(ucr)
	if err != nil {
		code := errorhandler.GetStatusCode(err)
		if code >= 500 {
			logger.Error.Println("Failed to confirm confirmation code: " + err.Error())
		}
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "text/html")
		views.ErrorIndex(code, http.StatusText(code)).Render(r.Context(), w)
		return
	}
	w.Write([]byte("User login info was sent to " + ucr.Email))
}

func (uh *UserHandler) HandleGetUsersAll(w http.ResponseWriter, r *http.Request) {
	users, err := uh.UserService.GetUsersAll()
	if err != nil {
		code := errorhandler.GetStatusCode(err)
		if code >= 500 {
			logger.Error.Println("Failed to get users from the database: " + err.Error())
		}
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "text/html")
		views.ErrorIndex(code, http.StatusText(code)).Render(r.Context(), w)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	views.UserListIndex(users).Render(r.Context(), w)
}

func (uh *UserHandler) HandleGetUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		code := http.StatusUnprocessableEntity
		if code >= 500 {
			logger.Error.Println("Failed to get id from request: " + err.Error())
		}
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "text/html")
		views.ErrorIndex(code, http.StatusText(code)).Render(r.Context(), w)
		return
	}

	user, err := uh.UserService.GetUserByID(id)
	if err != nil {
		code := errorhandler.GetStatusCode(err)
		if code >= 500 {
			logger.Error.Println("Failed to get user from the database: " + err.Error())
		}
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "text/html")
		views.ErrorIndex(code, http.StatusText(code)).Render(r.Context(), w)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	views.UserIndex(user).Render(r.Context(), w)
}
