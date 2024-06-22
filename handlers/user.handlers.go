package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/bmg-c/product-diary/errorhandler"
	L "github.com/bmg-c/product-diary/localization"
	"github.com/bmg-c/product-diary/logger"
	"github.com/bmg-c/product-diary/schemas"
	"github.com/bmg-c/product-diary/schemas/user_schemas"
	"github.com/bmg-c/product-diary/views"
	"github.com/bmg-c/product-diary/views/user_views"
	"github.com/google/uuid"
)

const locale = L.LocaleEnUS

type UserService interface {
	GetUser(userInfo user_schemas.GetUser) (user_schemas.UserPublic, error)
	GetUsersAll() ([]user_schemas.UserPublic, error)
	SigninUser(ur user_schemas.UserSignin) error
	ConfirmSignin(ucr user_schemas.UserConfirmSignin) error
	LoginUser(ul user_schemas.UserLogin) (uuid.UUID, error)
	GetUserBySession(sessionInfo user_schemas.GetSession) (user_schemas.UserDB, error)
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

	sessionUUID, err := uh.getUserSessionCookieValue(r)
	if err != nil {
		errText := l.Localize(err.Error())
		code := errorhandler.GetStatusCode(err)
		if code == http.StatusUnprocessableEntity {
			uh.deleteUserSessionCookie(w)
		} else {
			w.WriteHeader(code)
			w.Write([]byte(errText))
			return
		}
	}
	sessionInfo := user_schemas.GetSession{
		SessionUUID: sessionUUID,
	}
	userDB, _ := uh.UserService.GetUserBySession(sessionInfo)
	authorized := false
	if !schemas.IsZero(userDB) {
		authorized = true
	}

	err = user_views.UserControls(l, authorized).Render(r.Context(), w)
	if err != nil {
		errText := l.Localize(err.Error())
		code := errorhandler.GetStatusCode(err)
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

func (uh *UserHandler) HandleLoginIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	l := L.NewLocilizer(locale)

	err := user_views.LoginIndex(l).Render(r.Context(), w)
	if err != nil {
		errText := l.Localize(err.Error())
		code := http.StatusInternalServerError
		w.WriteHeader(code)
		w.Write([]byte(errText))
		return
	}
}

func (uh *UserHandler) HandleLoginLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	l := L.NewLocilizer(locale)

	var input user_schemas.UserLogin = user_schemas.UserLogin{}

	err := r.ParseForm()
	input.Password = r.Form.Get("password")
	input.Email = r.Form.Get("email")
	ve := schemas.ValidateStruct(input)
	if ve != nil || schemas.IsZero(input) || err != nil {
		err = L.GetError(L.MsgErrorGetUserNotFound)
		errText := l.Localize(err.Error())
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(errText))
		return
	}

	sessionUUID, err := uh.UserService.LoginUser(input)
	if err != nil {
		errText := l.Localize(err.Error())
		code := errorhandler.GetStatusCode(err)
		if code >= 500 {
			logger.Error.Println("Failed to login user: " + errText)
		}
		w.WriteHeader(code)
		w.Write([]byte(errText))
		return
	}
	uh.setUserSessionCookie(w, sessionUUID)

	user_views.UsersPage(l).Render(r.Context(), w)
}

func (uh *UserHandler) HandleProfileIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	l := L.NewLocilizer(locale)

	sessionUUID, err := uh.getUserSessionCookieValue(r)
	if err != nil {
		errText := l.Localize(err.Error())
		code := errorhandler.GetStatusCode(err)
		if code == http.StatusUnprocessableEntity {
			uh.deleteUserSessionCookie(w)
		} else {
			w.WriteHeader(code)
			w.Write([]byte(errText))
			return
		}
	}
	sessionInfo := user_schemas.GetSession{
		SessionUUID: sessionUUID,
	}
	userDB, err := uh.UserService.GetUserBySession(sessionInfo)
	up := user_schemas.UserPublic{
		UserID:    userDB.UserID,
		Email:     userDB.Email,
		Username:  userDB.Username,
		CreatedAt: userDB.CreatedAt,
	}

	err = user_views.User(up).Render(r.Context(), w)
	if err != nil {
		errText := l.Localize(err.Error())
		code := http.StatusInternalServerError
		w.WriteHeader(code)
		w.Write([]byte(errText))
		return
	}
}

func (uh *UserHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	l := L.NewLocilizer(locale)

	uh.deleteUserSessionCookie(w)

	err := user_views.UsersPage(l).Render(r.Context(), w)
	if err != nil {
		errText := l.Localize(err.Error())
		code := http.StatusInternalServerError
		w.WriteHeader(code)
		w.Write([]byte(errText))
		return
	}
}

func (uh *UserHandler) setUserSessionCookie(w http.ResponseWriter, sessionUUID uuid.UUID) {
	cookie := http.Cookie{
		Name:     "session",
		Value:    sessionUUID.String(),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)
}

func (uh *UserHandler) getUserSessionCookieValue(r *http.Request) (uuid.UUID, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return uuid.UUID{}, nil
		}
		return uuid.UUID{}, errorhandler.StatusError{
			Err:  L.GetError(L.MsgErrorInternalServer),
			Code: http.StatusInternalServerError,
		}
	}
	sessionUUID, err := uuid.Parse(cookie.Value)
	if err != nil {
		return uuid.UUID{}, errorhandler.StatusError{
			Err:  err,
			Code: http.StatusUnprocessableEntity,
		}
	}
	return sessionUUID, nil
}

func (uh *UserHandler) deleteUserSessionCookie(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	}

	http.SetCookie(w, &cookie)
}
