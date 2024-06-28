package handlers

import (
	"errors"
	"net/http"
	"time"

	E "github.com/bmg-c/product-diary/errorhandler"
	"github.com/bmg-c/product-diary/logger"
	"github.com/bmg-c/product-diary/schemas"
	"github.com/bmg-c/product-diary/schemas/item_schemas"
	"github.com/bmg-c/product-diary/schemas/user_schemas"
	"github.com/bmg-c/product-diary/util"
	"github.com/bmg-c/product-diary/views/product_views"
)

func NewItemHandler(itemService ItemService, userService UserService) *ItemHandler {
	return &ItemHandler{
		itemService: itemService,
		userService: userService,
	}
}

type ItemHandler struct {
	itemService ItemService
	userService UserService
}

func (ih *ItemHandler) HandleGetItems(w http.ResponseWriter, r *http.Request) {
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
		userDB, err = ih.userService.GetUserBySession(sessionUUID)
		if err != nil {
			logger.Error.Printf("Failure getting user from valid session cookie.\n")
		}
	}

	var input item_schemas.GetItems = item_schemas.GetItems{}

	err = r.ParseForm()
	if err != nil {
		code = http.StatusInternalServerError
		logger.Error.Printf("Error???? %v\n", err)
		return
	}
	input.ItemDate, err = time.Parse("2006-01-02", r.Form.Get("item_date"))
	if err != nil {
		code = http.StatusUnprocessableEntity
		return
	}
	input.UserID = userDB.UserID

	ve := schemas.ValidateStruct(input)
	if ve != nil {
		code = http.StatusUnprocessableEntity
		return
	}

	items, err := ih.itemService.GetItems(input)
	if err != nil {
		switch err {
		case E.ErrUnprocessableEntity:
			code = http.StatusUnprocessableEntity
			return
		default:
			code = http.StatusInternalServerError
			logger.Error.Printf("Server error %v\n", err)
			return
		}
	}

	inputUser := user_schemas.GetUser{
		UserID: userDB.UserID,
	}
	persons, err := ih.userService.GetUserPersons(inputUser)
	if err != nil {
		switch err {
		case E.ErrUnprocessableEntity:
			code = http.StatusUnprocessableEntity
			return
		default:
			code = http.StatusInternalServerError
			logger.Error.Printf("Server error %v\n", err)
			return
		}
	}

	util.RenderComponent(&out, product_views.ItemList(l, items, persons), r)
}

func (ih *ItemHandler) HandleAddItem(w http.ResponseWriter, r *http.Request) {
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
		userDB, err = ih.userService.GetUserBySession(sessionUUID)
		if err != nil {
			logger.Error.Printf("Failure getting user from valid session cookie.\n")
		}
	}

	var input item_schemas.AddItem = item_schemas.AddItem{}

	err = r.ParseForm()
	if err != nil {
		code = http.StatusInternalServerError
		logger.Error.Printf("Error???? %v\n", err)
		return
	}
	input.ProductID, err = util.GetUintFromString(r.Form.Get("product_id"))
	if err != nil {
		code = http.StatusUnprocessableEntity
		return
	}
	input.ItemDate, err = time.Parse("2006-01-02", r.Form.Get("item_date"))
	if err != nil {
		code = http.StatusUnprocessableEntity
		return
	}
	input.UserID = userDB.UserID

	ve := schemas.ValidateStruct(input)
	if ve != nil {
		code = http.StatusUnprocessableEntity
		return
	}

	itemParsed, err := ih.itemService.AddItem(input)
	if err != nil {
		switch err {
		case E.ErrUnprocessableEntity:
			code = http.StatusUnprocessableEntity
			return
		default:
			code = http.StatusInternalServerError
			logger.Error.Printf("Server error %v\n", err)
			return
		}
	}

	inputUser := user_schemas.GetUser{
		UserID: userDB.UserID,
	}
	persons, err := ih.userService.GetUserPersons(inputUser)
	if err != nil {
		switch err {
		case E.ErrUnprocessableEntity:
			code = http.StatusUnprocessableEntity
			return
		default:
			code = http.StatusInternalServerError
			logger.Error.Printf("Server error %v\n", err)
			return
		}
	}

	util.RenderComponent(&out, product_views.Item(l, itemParsed, persons), r)
}

func (ih *ItemHandler) HandleDeleteItem(w http.ResponseWriter, r *http.Request) {
	_ = util.InitHTMLHandler(w, r)
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
		userDB, err = ih.userService.GetUserBySession(sessionUUID)
		if err != nil {
			logger.Error.Printf("Failure getting user from valid session cookie.\n")
		}
	}

	var input item_schemas.DeleteItem = item_schemas.DeleteItem{}

	err = r.ParseForm()
	if err != nil {
		code = http.StatusInternalServerError
		logger.Error.Printf("Error???? %v\n", err)
		return
	}
	input.ItemID, err = util.GetUintFromString(r.Form.Get("item_id"))
	if err != nil {
		code = http.StatusUnprocessableEntity
		return
	}
	input.UserID = userDB.UserID

	ve := schemas.ValidateStruct(input)
	if ve != nil {
		code = http.StatusUnprocessableEntity
		return
	}

	err = ih.itemService.DeleteItem(input)
	if err != nil {
		switch err {
		case E.ErrUnprocessableEntity:
			code = http.StatusUnprocessableEntity
			return
		default:
			code = http.StatusInternalServerError
			logger.Error.Printf("Server error %v\n", err)
			return
		}
	}
}

func (ih *ItemHandler) HandleChangeItem(w http.ResponseWriter, r *http.Request) {
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
		userDB, err = ih.userService.GetUserBySession(sessionUUID)
		if err != nil {
			logger.Error.Printf("Failure getting user from valid session cookie.\n")
		}
	}

	var input item_schemas.ChangeItem = item_schemas.ChangeItem{}

	err = r.ParseForm()
	if err != nil {
		code = http.StatusInternalServerError
		logger.Error.Printf("Error???? %v\n", err)
		return
	}
	input.ItemID, err = util.GetUintFromString(r.Form.Get("item_id"))
	if err != nil {
		code = http.StatusUnprocessableEntity
		return
	}
	input.ProductID, _ = util.GetUintFromString(r.Form.Get("product_id"))
	input.ItemCost, _ = util.GetFloatFromString(r.Form.Get("item_cost"))
	input.ItemAmount, _ = util.GetFloatFromString(r.Form.Get("item_amount"))
	typ, _ := util.GetUintFromString(r.Form.Get("item_type"))
	input.ItemType = uint8(typ)
	input.PersonID, _ = util.GetUintFromString(r.Form.Get("person_id"))
	input.UserID = userDB.UserID

	ve := schemas.ValidateStruct(input)
	if ve != nil {
		code = http.StatusUnprocessableEntity
		return
	}

	itemParsed, err := ih.itemService.ChangeItem(input)
	if err != nil {
		switch err {
		case E.ErrUnprocessableEntity:
			code = http.StatusUnprocessableEntity
			return
		default:
			code = http.StatusInternalServerError
			logger.Error.Printf("Server error %v\n", err)
			return
		}
	}

	inputUser := user_schemas.GetUser{
		UserID: userDB.UserID,
	}
	persons, err := ih.userService.GetUserPersons(inputUser)
	if err != nil {
		switch err {
		case E.ErrUnprocessableEntity:
			code = http.StatusUnprocessableEntity
			return
		default:
			code = http.StatusInternalServerError
			logger.Error.Printf("Server error %v\n", err)
			return
		}
	}

	util.RenderComponent(&out, product_views.Item(l, itemParsed, persons), r)
}
