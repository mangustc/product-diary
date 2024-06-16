package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/bmg-c/product-diary/errorhandler"
	"github.com/bmg-c/product-diary/logger"
	"github.com/bmg-c/product-diary/services"
	"github.com/bmg-c/product-diary/views"
	"github.com/bmg-c/product-diary/views/user_views"
)

type UserService interface {
	GetUserByID(id int) (services.UserPublic, error)
	GetUsersAll() ([]services.UserPublic, error)
	SigninUser(ur services.UserSignin) error
	ConfirmSignin(ucr services.UserConfirmSignin) error
}

func NewUserHandler(us UserService) *UserHandler {
	return &UserHandler{
		UserService: us,
	}
}

type UserHandler struct {
	UserService UserService
}

func (uh *UserHandler) HandleUsersPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := user_views.UsersPage().Render(r.Context(), w)
	if err != nil {
		code := http.StatusInternalServerError
		w.WriteHeader(code)
		views.ErrorIndex(code, http.StatusText(code)).Render(r.Context(), w)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (uh *UserHandler) HandleSigninIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	err := user_views.SigninIndex().Render(r.Context(), w)
	if err != nil {
		code := http.StatusInternalServerError
		w.WriteHeader(code)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (uh *UserHandler) HandleControlsIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	err := user_views.UserControls(false).Render(r.Context(), w)
	if err != nil {
		code := http.StatusInternalServerError
		w.WriteHeader(code)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (uh *UserHandler) HandleSigninSignin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	err := r.ParseForm()
	email := r.Form.Get("email")
	if email == "" {
		err = fmt.Errorf("Email is not provided")
	}
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(err.Error()))
		return
	}

	err = uh.UserService.SigninUser(services.UserSignin{
		Email: email,
	})
	if err != nil {
		code := errorhandler.GetStatusCode(err)
		if code >= 500 {
			logger.Error.Println("Failed to register user: " + err.Error())
		}
		w.WriteHeader(code)
		w.Write([]byte(err.Error()))
		return
	}
	err = user_views.ConfirmSignin(email).Render(r.Context(), w)
	if err != nil {
		code := http.StatusInternalServerError
		w.WriteHeader(code)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (uh *UserHandler) HandleConfirmSignin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	err := r.ParseForm()
	email := r.Form.Get("email")
	code := r.Form.Get("code")
	if email == "" || code == "" {
		err = fmt.Errorf("Code is not provided")
	}
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(err.Error()))
		return
	}

	err = uh.UserService.ConfirmSignin(services.UserConfirmSignin{
		Email: email,
		Code:  code,
	})
	if err != nil {
		code := errorhandler.GetStatusCode(err)
		if code >= 500 {
			logger.Error.Println("Failed to confirm confirmation code: " + err.Error())
		}
		w.WriteHeader(code)
		w.Write([]byte(err.Error()))
		return
	}
	err = user_views.EndSignin(email).Render(r.Context(), w)
	if err != nil {
		code := http.StatusInternalServerError
		w.WriteHeader(code)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (uh *UserHandler) HandleGetUsersAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	users, err := uh.UserService.GetUsersAll()
	if err != nil {
		code := errorhandler.GetStatusCode(err)
		if code >= 500 {
			logger.Error.Println("Failed to get users from the database: " + err.Error())
		}
		w.WriteHeader(code)
		w.Write([]byte(err.Error()))
		return
	}

	user_views.UserlistIndex(users).Render(r.Context(), w)
}

func (uh *UserHandler) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	err := r.ParseForm()
	id, err := strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(err.Error()))
		return
	}

	user, err := uh.UserService.GetUserByID(id)
	if err != nil {
		code := errorhandler.GetStatusCode(err)
		if code >= 500 {
			logger.Error.Println("Failed to get user from the database: " + err.Error())
		}
		w.WriteHeader(code)
		w.Write([]byte(err.Error()))
		return
	}

	user_views.User(user).Render(r.Context(), w)
}

func (uh *UserHandler) HandleUserIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := user_views.UserIndex().Render(r.Context(), w)
	if err != nil {
		code := http.StatusInternalServerError
		w.WriteHeader(code)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}
