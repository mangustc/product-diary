package handlers

import (
	"errors"
	"net/http"
	"strconv"

	E "github.com/bmg-c/product-diary/errorhandler"
	L "github.com/bmg-c/product-diary/localization"
	"github.com/bmg-c/product-diary/logger"
	"github.com/bmg-c/product-diary/schemas"
	"github.com/bmg-c/product-diary/schemas/user_schemas"
	"github.com/bmg-c/product-diary/util"
	"github.com/bmg-c/product-diary/views/user_views"
	"github.com/google/uuid"
)

type UserService interface {
	GetUser(userInfo user_schemas.GetUser) (user_schemas.UserPublic, error)
	GetUsersAll() ([]user_schemas.UserPublic, error)
	SigninUser(ur user_schemas.UserSignin) error
	ConfirmSignin(ucr user_schemas.UserConfirmSignin) error
	LoginUser(ul user_schemas.UserLogin) (uuid.UUID, error)
	GetUserBySession(sessionUUID uuid.UUID) (user_schemas.UserDB, error)
	AddPerson(personInfo user_schemas.GetPerson) (user_schemas.PersonDB, error)
	GetUserPersons(userInfo user_schemas.GetUser) ([]user_schemas.PersonDB, error)
	ToggleHiddenPerson(personInfo user_schemas.GetPerson) (user_schemas.PersonDB, error)
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
	l := util.InitHTMLHandler(w, r)
	var code int = http.StatusOK
	var out []byte
	defer util.RespondHTTP(w, &code, &out)

	util.RenderComponent(&out, user_views.UsersPage(l), r)
}

func (uh *UserHandler) HandleControlsIndex(w http.ResponseWriter, r *http.Request) {
	l := util.InitHTMLHandler(w, r)
	var code int = http.StatusOK
	var out []byte
	defer util.RespondHTTP(w, &code, &out)

	var authorized bool = false
	// var userDB user_schemas.UserDB
	sessionUUID, err := util.GetUserSessionCookieValue(w, r)
	if err != nil {
		if errors.Is(err, E.ErrInternalServer) {
			logger.Error.Printf("Failure getting session cookie.\n")
		}
	} else {
		// userDB, err = uh.UserService.GetUserBySession(sessionUUID)
		_, err = uh.UserService.GetUserBySession(sessionUUID)
		if err != nil {
			logger.Error.Printf("Failure getting user from valid session cookie.\n")
		} else {
			authorized = true
		}
	}

	util.RenderComponent(&out, user_views.UserControls(l, authorized), r)
}

func (uh *UserHandler) HandleSigninIndex(w http.ResponseWriter, r *http.Request) {
	l := util.InitHTMLHandler(w, r)
	var code int = http.StatusOK
	var out []byte
	defer util.RespondHTTP(w, &code, &out)

	data := user_views.SigninData{
		CodeSent: false,
		Email:    "",
	}
	util.RenderComponent(&out, user_views.SigninIndex(l, data), r)
}

func (uh *UserHandler) HandleSigninSignin(w http.ResponseWriter, r *http.Request) {
	l := util.InitHTMLHandler(w, r)
	var code int = http.StatusOK
	var out []byte
	defer util.RespondHTTP(w, &code, &out)

	var input user_schemas.UserSignin = user_schemas.UserSignin{}
	var inputConfirm user_schemas.UserConfirmSignin = user_schemas.UserConfirmSignin{}

	var hasCode bool = false
	err := r.ParseForm()
	if err != nil {
		err = E.ErrUnprocessableEntity
		data := user_views.SigninData{
			CodeSent: false,
			Email:    input.Email,
			Err:      L.GetError(L.MsgErrorEmailWrong),
		}
		util.RenderComponent(&out, user_views.SigninIndex(l, data), r)
		return
	}
	hasCode = r.Form.Has("code")
	if hasCode {
		inputConfirm.Email = r.Form.Get("email")
		inputConfirm.Code = r.Form.Get("code")
		ve := schemas.ValidateStruct(inputConfirm)
		if ve != nil {
			err = E.ErrUnprocessableEntity
		}
	} else {
		input.Email = r.Form.Get("email")
		ve := schemas.ValidateStruct(input)
		if ve != nil {
			err = E.ErrUnprocessableEntity
		}
	}
	if err != nil {
		switch err {
		case E.ErrUnprocessableEntity:
			code = http.StatusUnprocessableEntity
			var email string
			if hasCode {
				err = L.GetError(L.MsgErrorCodeWrong)
				email = inputConfirm.Email
			} else {
				err = L.GetError(L.MsgErrorEmailWrong)
				email = input.Email
			}
			data := user_views.SigninData{
				CodeSent: hasCode,
				Email:    email,
				Err:      err,
			}
			util.RenderComponent(&out, user_views.SigninIndex(l, data), r)
			return
		default:
			code = http.StatusInternalServerError
			return
		}
	}

	if hasCode {
		err = uh.UserService.ConfirmSignin(inputConfirm)
		if err != nil {
			switch err {
			case E.ErrUnprocessableEntity:
				code = http.StatusUnprocessableEntity
				data := user_views.SigninData{
					CodeSent: hasCode,
					Email:    inputConfirm.Email,
					Err:      L.GetError(L.MsgErrorCodeWrong),
				}
				util.RenderComponent(&out, user_views.SigninIndex(l, data), r)
				return
			default:
				code = http.StatusInternalServerError
				logger.Error.Println("Error confirming confirmation code from user")
				return
			}
		}

		util.RenderComponent(&out, user_views.EndSignin(l, inputConfirm.Email), r)
	} else {
		err = uh.UserService.SigninUser(input)
		if err != nil {
			switch err {
			case E.ErrUnprocessableEntity:
				code = http.StatusUnprocessableEntity
				data := user_views.SigninData{
					CodeSent: hasCode,
					Email:    inputConfirm.Email,
					Err:      L.GetError(L.MsgEmailExists),
				}
				util.RenderComponent(&out, user_views.SigninIndex(l, data), r)
				return
			default:
				code = http.StatusInternalServerError
				logger.Error.Println("Error confirming confirmation code from user")
				return
			}
		}

		data := user_views.SigninData{
			CodeSent: true,
			Email:    input.Email,
			Err:      nil,
		}
		util.RenderComponent(&out, user_views.SigninIndex(l, data), r)
		return
	}
}

func (uh *UserHandler) HandleGetUsersAll(w http.ResponseWriter, r *http.Request) {
	_ = util.InitHTMLHandler(w, r)
	var code int = http.StatusOK
	var out []byte
	defer util.RespondHTTP(w, &code, &out)

	users, err := uh.UserService.GetUsersAll()
	if err != nil {
		code = http.StatusInternalServerError
		logger.Error.Println("Failure getting users from the database.")
		return
	}

	util.RenderComponent(&out, user_views.UserlistIndex(users), r)
}

func (uh *UserHandler) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	l := util.InitHTMLHandler(w, r)
	var code int = http.StatusOK
	var out []byte
	defer util.RespondHTTP(w, &code, &out)

	var input user_schemas.GetUser = user_schemas.GetUser{}

	err := r.ParseForm()
	id64, err := strconv.ParseUint(r.Form.Get("id"), 10, 0)
	input.UserID = uint(id64)
	input.Email = r.Form.Get("email")
	ve := schemas.ValidateStruct(input)
	if ve != nil || schemas.IsZero(input) || err != nil {
		code = http.StatusNotFound
		util.RenderComponent(&out, user_views.ErrorMsg(l, L.GetError(L.MsgErrorGetUserNotFound)), r)
		return
	}

	user, err := uh.UserService.GetUser(input)
	if err != nil {
		code = http.StatusNotFound
		util.RenderComponent(&out, user_views.ErrorMsg(l, L.GetError(L.MsgErrorGetUserNotFound)), r)
		return
	}

	util.RenderComponent(&out, user_views.User(user), r)
}

func (uh *UserHandler) HandleUserIndex(w http.ResponseWriter, r *http.Request) {
	l := util.InitHTMLHandler(w, r)
	var code int = http.StatusOK
	var out []byte
	defer util.RespondHTTP(w, &code, &out)

	util.RenderComponent(&out, user_views.UserIndex(l), r)
}

func (uh *UserHandler) HandleLoginIndex(w http.ResponseWriter, r *http.Request) {
	l := util.InitHTMLHandler(w, r)
	var code int = http.StatusOK
	var out []byte
	defer util.RespondHTTP(w, &code, &out)

	data := user_views.LoginData{
		Email:        "",
		Password:     "",
		SuccessLogin: false,
	}
	util.RenderComponent(&out, user_views.LoginIndex(l, data), r)
}

func (uh *UserHandler) HandleLoginLogin(w http.ResponseWriter, r *http.Request) {
	l := util.InitHTMLHandler(w, r)
	var code int = http.StatusOK
	var out []byte
	defer util.RespondHTTP(w, &code, &out)

	var input user_schemas.UserLogin = user_schemas.UserLogin{}

	err := r.ParseForm()
	input.Password = r.Form.Get("password")
	input.Email = r.Form.Get("email")
	data := user_views.LoginData{
		Email:    input.Email,
		Password: input.Password,
	}
	ve := schemas.ValidateStruct(input)
	if ve != nil || schemas.IsZero(input) || err != nil {
		code = http.StatusNotFound
		util.RenderComponent(&out, user_views.LoginIndex(l, data), r)
		util.RenderComponent(&out, user_views.ErrorMsg(l, L.GetError(L.MsgErrorGetUserNotFound)), r)
		return
	}

	sessionUUID, err := uh.UserService.LoginUser(input)
	if err != nil {
		code = http.StatusNotFound
		util.RenderComponent(&out, user_views.LoginIndex(l, data), r)
		util.RenderComponent(&out, user_views.ErrorMsg(l, L.GetError(L.MsgErrorGetUserNotFound)), r)
		return
	}
	util.SetUserSessionCookie(w, sessionUUID)

	w.Header().Set("HX-Redirect", r.Header.Get("Referer"))
}

func (uh *UserHandler) HandleProfileIndex(w http.ResponseWriter, r *http.Request) {
	l := util.InitHTMLHandler(w, r)
	var code int = http.StatusOK
	var out []byte
	defer util.RespondHTTP(w, &code, &out)

	var userDB user_schemas.UserDB
	sessionUUID, err := util.GetUserSessionCookieValue(w, r)
	if err != nil {
		if errors.Is(err, E.ErrInternalServer) {
			logger.Error.Printf("Failure getting session cookie.\n")
		}
	} else {
		userDB, err = uh.UserService.GetUserBySession(sessionUUID)
		if err != nil {
			logger.Error.Printf("Failure getting user from valid session cookie.\n")
		}
	}
	up := user_schemas.UserPublic{
		UserID:    userDB.UserID,
		Email:     userDB.Email,
		Username:  userDB.Username,
		CreatedAt: userDB.CreatedAt,
	}

	userInfo := user_schemas.GetUser{
		UserID: userDB.UserID,
	}

	persons, err := uh.UserService.GetUserPersons(userInfo)
	if err != nil {
		logger.Error.Printf("Erorr: %v\n", err)
	}

	util.RenderComponent(&out, user_views.ProfileBlock(l, up, persons), r)
}

func (uh *UserHandler) HandleTogglePerson(w http.ResponseWriter, r *http.Request) {
	l := util.InitHTMLHandler(w, r)
	var code int = http.StatusOK
	var out []byte
	defer util.RespondHTTP(w, &code, &out)

	var userDB user_schemas.UserDB
	sessionUUID, err := util.GetUserSessionCookieValue(w, r)
	if err != nil {
		if errors.Is(err, E.ErrInternalServer) {
			logger.Error.Printf("Failure getting session cookie.\n")
		}
	} else {
		userDB, err = uh.UserService.GetUserBySession(sessionUUID)
		if err != nil {
			logger.Error.Printf("Failure getting user from valid session cookie.\n")
		}
	}

	var input user_schemas.GetPerson = user_schemas.GetPerson{}
	err = r.ParseForm()
	input.UserID = userDB.UserID
	input.PersonName = r.Form.Get("person_name")
	ve := schemas.ValidateStruct(input)
	if ve != nil || schemas.IsZero(input) || err != nil {
		code = http.StatusUnprocessableEntity
		util.RenderComponent(&out, user_views.ErrorMsg(l, L.GetError(L.MsgErrorUsernameEmpty)), r)
		return
	}

	personDB, err := uh.UserService.ToggleHiddenPerson(input)
	if err != nil {
		switch err {
		case E.ErrUnprocessableEntity:
			code = http.StatusUnprocessableEntity
			util.RenderComponent(&out, user_views.ErrorMsg(l, L.GetError(L.MsgErrorGetUserNotFound)), r)
			return
		default:
			code = http.StatusInternalServerError
			logger.Error.Println("Failure getting users from the database.")
			return
		}
	}

	util.RenderComponent(&out, user_views.Person(l, personDB), r)
}

func (uh *UserHandler) HandleAddPerson(w http.ResponseWriter, r *http.Request) {
	l := util.InitHTMLHandler(w, r)
	var code int = http.StatusOK
	var out []byte
	defer util.RespondHTTP(w, &code, &out)

	var userDB user_schemas.UserDB
	sessionUUID, err := util.GetUserSessionCookieValue(w, r)
	if err != nil {
		if errors.Is(err, E.ErrInternalServer) {
			logger.Error.Printf("Failure getting session cookie.\n")
		}
	} else {
		userDB, err = uh.UserService.GetUserBySession(sessionUUID)
		if err != nil {
			logger.Error.Printf("Failure getting user from valid session cookie.\n")
		}
	}

	var input user_schemas.GetPerson = user_schemas.GetPerson{}
	err = r.ParseForm()
	input.UserID = userDB.UserID
	input.PersonName = r.Form.Get("person_name")
	ve := schemas.ValidateStruct(input)
	if ve != nil || schemas.IsZero(input) || err != nil {
		code = http.StatusUnprocessableEntity
		util.RenderComponent(&out, user_views.ErrorMsg(l, L.GetError(L.MsgErrorUsernameEmpty)), r)
		return
	}

	personDB, err := uh.UserService.AddPerson(input)
	if err != nil {
		switch err {
		case E.ErrUnprocessableEntity:
			code = http.StatusUnprocessableEntity
			util.RenderComponent(&out, user_views.ErrorMsg(l, L.GetError(L.MsgErrorUsernameAlreadyExists)), r)
			return
		default:
			code = http.StatusInternalServerError
			logger.Error.Println("Failure getting users from the database.")
			return
		}
	}

	util.RenderComponent(&out, user_views.Person(l, personDB), r)
}

func (uh *UserHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	_ = util.InitHTMLHandler(w, r)
	var code int = http.StatusOK
	var out []byte
	defer util.RespondHTTP(w, &code, &out)

	util.DeleteUserSessionCookie(w)

	w.Header().Set("HX-Redirect", r.Header.Get("Referer"))
}
