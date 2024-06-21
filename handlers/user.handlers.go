package handlers

import (
	"net/http"
	"strconv"

	"github.com/bmg-c/product-diary/errorhandler"
	L "github.com/bmg-c/product-diary/localization"
	"github.com/bmg-c/product-diary/logger"
	"github.com/bmg-c/product-diary/schemas"
	"github.com/bmg-c/product-diary/schemas/user_schemas"
	"github.com/bmg-c/product-diary/views"
	"github.com/bmg-c/product-diary/views/user_views"
)

const locale = L.LocaleRuRU

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
	l := L.NewLocilizer(locale)

	err := user_views.UsersPage(l).Render(r.Context(), w)
	if err != nil {
		errText := l.Localize(err.Error())
		code := http.StatusInternalServerError
		w.WriteHeader(code)
		views.ErrorIndex(code, errText).Render(r.Context(), w)
		return
	}
}

func (uh *UserHandler) HandleSigninIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	l := L.NewLocilizer(locale)

	err := user_views.SigninIndex(l).Render(r.Context(), w)
	if err != nil {
		errText := l.Localize(err.Error())
		code := http.StatusInternalServerError
		w.WriteHeader(code)
		w.Write([]byte(errText))
		return
	}
}

func (uh *UserHandler) HandleControlsIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	l := L.NewLocilizer(locale)

	err := user_views.UserControls(l, false).Render(r.Context(), w)
	if err != nil {
		errText := l.Localize(err.Error())
		code := http.StatusInternalServerError
		w.WriteHeader(code)
		w.Write([]byte(errText))
		return
	}
}

func (uh *UserHandler) HandleSigninSignin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	l := L.NewLocilizer(locale)

	var input user_schemas.UserSignin = user_schemas.UserSignin{}

	err := r.ParseForm()
	input.Email = r.Form.Get("email")
	ve := schemas.ValidateStruct(input)
	if ve != nil {
		err = L.GetError(L.MsgErrorEmailWrong)
	}
	if err != nil {
		err = L.GetError(L.MsgErrorEmailWrong)
		errText := l.Localize(err.Error())
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(errText))
		return
	}

	err = uh.UserService.SigninUser(input)
	if err != nil {
		errText := l.Localize(err.Error())
		code := errorhandler.GetStatusCode(err)
		if code >= 500 {
			logger.Error.Println("Failed to register user: " + errText)
		}
		w.WriteHeader(code)
		w.Write([]byte(errText))
		return
	}
	err = user_views.ConfirmSignin(l, input.Email).Render(r.Context(), w)
	if err != nil {
		errText := l.Localize(err.Error())
		code := http.StatusInternalServerError
		w.WriteHeader(code)
		w.Write([]byte(errText))
		return
	}
}

func (uh *UserHandler) HandleConfirmSignin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	l := L.NewLocilizer(locale)

	var input user_schemas.UserConfirmSignin = user_schemas.UserConfirmSignin{}

	err := r.ParseForm()
	input.Email = r.Form.Get("email")
	input.Code = r.Form.Get("code")
	ve := schemas.ValidateStruct(input)
	if ve != nil {
		err = L.GetError(L.MsgErrorCodeWrong)
	}
	if err != nil {
		err = L.GetError(L.MsgErrorCodeWrong)
		errText := l.Localize(err.Error())
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(errText))
		return
	}

	err = uh.UserService.ConfirmSignin(input)
	if err != nil {
		errText := l.Localize(err.Error())
		code := errorhandler.GetStatusCode(err)
		if code >= 500 {
			logger.Error.Println("Failed to confirm confirmation code: " + errText)
		}
		w.WriteHeader(code)
		w.Write([]byte(errText))
		return
	}
	err = user_views.EndSignin(l, input.Email).Render(r.Context(), w)
	if err != nil {
		errText := l.Localize(err.Error())
		code := http.StatusInternalServerError
		w.WriteHeader(code)
		w.Write([]byte(errText))
		return
	}
}

func (uh *UserHandler) HandleGetUsersAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	l := L.NewLocilizer(locale)

	users, err := uh.UserService.GetUsersAll()
	if err != nil {
		errText := l.Localize(err.Error())
		code := errorhandler.GetStatusCode(err)
		if code >= 500 {
			logger.Error.Println("Failed to get users from the database: " + errText)
		}
		w.WriteHeader(code)
		w.Write([]byte(errText))
		return
	}

	user_views.UserlistIndex(users).Render(r.Context(), w)
}

func (uh *UserHandler) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	l := L.NewLocilizer(locale)

	var input user_schemas.GetUser = user_schemas.GetUser{}

	err := r.ParseForm()
	id64, err := strconv.ParseUint(r.Form.Get("id"), 10, 0)
	input.UserID = uint(id64)
	input.Email = r.Form.Get("email")
	ve := schemas.ValidateStruct(input)
	if ve != nil || schemas.IsZero(input) || err != nil {
		err = L.GetError(L.MsgErrorGetUserNotFound)
		errText := l.Localize(err.Error())
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(errText))
		return
	}

	user, err := uh.UserService.GetUser(input)
	if err != nil {
		errText := l.Localize(err.Error())
		code := errorhandler.GetStatusCode(err)
		if code >= 500 {
			logger.Error.Println("Failed to get user from the database: " + errText)
		}
		w.WriteHeader(code)
		w.Write([]byte(errText))
		return
	}

	user_views.User(user).Render(r.Context(), w)
}

func (uh *UserHandler) HandleUserIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	l := L.NewLocilizer(locale)

	err := user_views.UserIndex(l).Render(r.Context(), w)
	if err != nil {
		errText := l.Localize(err.Error())
		code := http.StatusInternalServerError
		w.WriteHeader(code)
		w.Write([]byte(errText))
		return
	}
}
