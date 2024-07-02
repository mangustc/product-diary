package handlers

import (
	"errors"
	"net/http"

	E "github.com/bmg-c/product-diary/errorhandler"
	L "github.com/bmg-c/product-diary/localization"
	"github.com/bmg-c/product-diary/logger"
	"github.com/bmg-c/product-diary/schemas"
	"github.com/bmg-c/product-diary/schemas/product_schemas"
	"github.com/bmg-c/product-diary/schemas/user_schemas"
	"github.com/bmg-c/product-diary/util"
	"github.com/bmg-c/product-diary/views/product_views"
)

func NewProductHandler(productService ProductService, userService UserService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		userService:    userService,
	}
}

type ProductHandler struct {
	productService ProductService
	userService    UserService
}

func (ph *ProductHandler) HandleProductsPage(w http.ResponseWriter, r *http.Request) {
	l := util.InitHTMLHandler(w, r)
	var code int = http.StatusOK
	var out []byte
	defer util.RespondHTTP(w, &code, &out)

	util.RenderComponent(&out, product_views.ProductsPage(l), r)
}

func (ph *ProductHandler) HandleAddProduct(w http.ResponseWriter, r *http.Request) {
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
		userDB, err = ph.userService.GetUserBySession(sessionUUID)
		if err != nil {
			logger.Error.Printf("Failure getting user from valid session cookie.\n")
		}
	}

	var input product_schemas.AddProduct = product_schemas.AddProduct{}
	var inputErrs product_views.ProductAddRowErrors = product_views.ProductAddRowErrors{}

	err = r.ParseForm()
	if err != nil {
		code = http.StatusInternalServerError
		logger.Error.Printf("Error???? %v\n", err)
		return
	}
	input.ProductTitle = r.Form.Get("product_title")
	input.ProductCalories, err = util.GetFloatFromString(r.Form.Get("product_calories"))
	if err != nil {
		code = http.StatusUnprocessableEntity
		inputErrs.CaloriesErr = L.GetError(L.MsgErrorProductCalories)
	}
	input.ProductFats, err = util.GetFloatFromString(r.Form.Get("product_fats"))
	if err != nil {
		code = http.StatusUnprocessableEntity
		inputErrs.FatsErr = L.GetError(L.MsgErrorProductNutrient)
	}
	input.ProductCarbs, err = util.GetFloatFromString(r.Form.Get("product_carbs"))
	if err != nil {
		code = http.StatusUnprocessableEntity
		inputErrs.CarbsErr = L.GetError(L.MsgErrorProductNutrient)
	}
	input.ProductProteins, err = util.GetFloatFromString(r.Form.Get("product_proteins"))
	if err != nil {
		code = http.StatusUnprocessableEntity
		inputErrs.ProteinsErr = L.GetError(L.MsgErrorProductNutrient)
	}
	input.UserID = userDB.UserID
	ve := schemas.ValidateStruct(input)
	if ve != nil {
		code = http.StatusUnprocessableEntity
		for _, fe := range ve {
			logger.Info.Println(fe.Name())
			switch fe.Name() {
			case "ProductTitle":
				inputErrs.TitleErr = L.GetError(L.MsgErrorProductTitle)
			case "ProductCalories":
				inputErrs.CaloriesErr = L.GetError(L.MsgErrorProductCalories)
			case "ProductFats":
				inputErrs.FatsErr = L.GetError(L.MsgErrorProductNutrient)
			case "ProductCarbs":
				inputErrs.CarbsErr = L.GetError(L.MsgErrorProductNutrient)
			case "ProductProteins":
				inputErrs.ProteinsErr = L.GetError(L.MsgErrorProductNutrient)
			}
		}
	}

	productDB, err := ph.productService.AddProduct(input)
	if err != nil {
		switch err {
		case E.ErrUnprocessableEntity:
			code = http.StatusUnprocessableEntity
		default:
			code = http.StatusInternalServerError
			logger.Error.Printf("Server error %v\n", err)
			return
		}
		util.RenderComponent(&out, product_views.ProductAddRow(l, input, inputErrs), r)
	} else {
		input = product_schemas.AddProduct{}
		util.RenderComponent(&out, product_views.ProductAddRow(l, input, inputErrs), r)
		util.RenderComponent(&out, product_views.Product(l, productDB, userDB.UserID), r)
	}
}

func (ph *ProductHandler) HandleGetProducts(w http.ResponseWriter, r *http.Request) {
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
		userDB, err = ph.userService.GetUserBySession(sessionUUID)
		if err != nil {
			logger.Error.Printf("Failure getting user from valid session cookie.\n")
		}
	}

	var input product_schemas.GetProducts = product_schemas.GetProducts{}

	err = r.ParseForm()
	input.SearchQuery = r.Form.Get("search_query")
	if err != nil {
		code = http.StatusInternalServerError
		logger.Error.Printf("Error???? %v\n", err)
	}

	products, err := ph.productService.GetProducts(input)
	if err != nil {
		code = http.StatusInternalServerError
		logger.Error.Printf("Server error %v\n", err)
	}

	util.RenderComponent(&out, product_views.ProductList(l, products, userDB.UserID), r)
}

func (ph *ProductHandler) HandleCopyProduct(w http.ResponseWriter, r *http.Request) {
	l := util.InitHTMLHandler(w, r)
	var code int = http.StatusOK
	var out []byte
	defer util.RespondHTTP(w, &code, &out)

	var input product_schemas.GetProduct = product_schemas.GetProduct{}

	err := r.ParseForm()
	if err != nil {
		code = http.StatusUnprocessableEntity
		return
	}
	input.ProductID, err = util.GetUintFromString(r.Form.Get("product_id"))
	if err != nil {
		code = http.StatusUnprocessableEntity
		return
	}

	productDB, err := ph.productService.GetProduct(input)
	if err != nil {
		switch err {
		case E.ErrNotFound:
			code = http.StatusUnprocessableEntity
			return
		default:
			code = http.StatusInternalServerError
			logger.Error.Println("Error", err)
			return
		}
	}

	addProduct := product_schemas.AddProduct{
		ProductTitle:    productDB.ProductTitle,
		ProductCalories: productDB.ProductCalories,
		ProductFats:     productDB.ProductFats,
		ProductCarbs:    productDB.ProductCarbs,
		ProductProteins: productDB.ProductProteins,
		UserID:          productDB.UserID,
	}

	util.RenderComponent(&out, product_views.ProductAddRow(l, addProduct, product_views.ProductAddRowErrors{}), r)
}

func (ph *ProductHandler) HandleDeleteProduct(w http.ResponseWriter, r *http.Request) {
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
		userDB, err = ph.userService.GetUserBySession(sessionUUID)
		if err != nil {
			logger.Error.Printf("Failure getting user from valid session cookie.\n")
		}
	}

	var input product_schemas.DeleteProduct = product_schemas.DeleteProduct{}

	err = r.ParseForm()
	if err != nil {
		code = http.StatusUnprocessableEntity
		return
	}
	input.ProductID, err = util.GetUintFromString(r.Form.Get("product_id"))
	if err != nil {
		code = http.StatusUnprocessableEntity
		return
	}
	input.UserID = userDB.UserID

	err = ph.productService.DeleteProduct(input)
	if err != nil {
		switch err {
		case E.ErrNotFound:
			code = http.StatusUnprocessableEntity
			return
		default:
			code = http.StatusInternalServerError
			logger.Error.Println("Error", err)
			return
		}
	}
}
