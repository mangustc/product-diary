package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/bmg-c/product-diary/errorhandler"
	"github.com/bmg-c/product-diary/logger"
	"github.com/bmg-c/product-diary/schemas"
	"github.com/bmg-c/product-diary/schemas/user_schemas"
	"github.com/bmg-c/product-diary/views"
	"github.com/bmg-c/product-diary/views/user_views"
)

type UserService interface {
	GetUser(userInfo user_schemas.GetUser) (user_schemas.UserPublic, error)
	GetUsersAll() ([]user_schemas.UserPublic, error)
	SigninUser(ur user_schemas.UserSignin) error
	ConfirmSignin(ucr user_schemas.UserConfirmSignin) error
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
}

func (uh *UserHandler) HandleSigninSignin(w http.ResponseWriter, r *http.Request) {
	var input user_schemas.UserSignin = user_schemas.UserSignin{}
	w.Header().Set("Content-Type", "text/html")

	err := r.ParseForm()
	input.Email = r.Form.Get("email")
	ve := schemas.ValidateStruct(input)
	if ve != nil {
		err = fmt.Errorf("Email should have format like mail@example.com")
	}
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(err.Error()))
		return
	}

	err = uh.UserService.SigninUser(input)
	if err != nil {
		code := errorhandler.GetStatusCode(err)
		if code >= 500 {
			logger.Error.Println("Failed to register user: " + err.Error())
		}
		w.WriteHeader(code)
		w.Write([]byte(err.Error()))
		return
	}
	err = user_views.ConfirmSignin(input.Email).Render(r.Context(), w)
	if err != nil {
		code := http.StatusInternalServerError
		w.WriteHeader(code)
		w.Write([]byte(err.Error()))
		return
	}
}

func (uh *UserHandler) HandleConfirmSignin(w http.ResponseWriter, r *http.Request) {
	var input user_schemas.UserConfirmSignin = user_schemas.UserConfirmSignin{}
	w.Header().Set("Content-Type", "text/html")

	err := r.ParseForm()
	input.Email = r.Form.Get("email")
	input.Code = r.Form.Get("code")
	ve := schemas.ValidateStruct(input)
	if ve != nil {
		err = fmt.Errorf("Confirmation codes do not match")
	}
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(err.Error()))
		return
	}

	err = uh.UserService.ConfirmSignin(input)
	if err != nil {
		code := errorhandler.GetStatusCode(err)
		if code >= 500 {
			logger.Error.Println("Failed to confirm confirmation code: " + err.Error())
		}
		w.WriteHeader(code)
		w.Write([]byte(err.Error()))
		return
	}
	err = user_views.EndSignin(input.Email).Render(r.Context(), w)
	if err != nil {
		code := http.StatusInternalServerError
		w.WriteHeader(code)
		w.Write([]byte(err.Error()))
		return
	}
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
	var input user_schemas.GetUser = user_schemas.GetUser{}
	w.Header().Set("Content-Type", "text/html")

	err := r.ParseForm()
	id64, err := strconv.ParseUint(r.Form.Get("id"), 10, 0)
	input.UserID = uint(id64)
	input.Email = r.Form.Get("email")
	ve := schemas.ValidateStruct(input)
	if ve != nil {
		err = fmt.Errorf("ID should be greater than 0")
	}
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(err.Error()))
		return
	}

	user, err := uh.UserService.GetUser(input)
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
}
