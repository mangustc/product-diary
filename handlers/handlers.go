package handlers

import (
	"github.com/bmg-c/product-diary/schemas/product_schemas"
	"github.com/bmg-c/product-diary/schemas/user_schemas"
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

type ProductService interface {
	AddProduct(data product_schemas.AddProduct) (product_schemas.ProductDB, error)
	GetProducts(data product_schemas.GetProducts) ([]product_schemas.ProductDB, error)
	GetProduct(data product_schemas.GetProduct) (product_schemas.ProductDB, error)
	DeleteProduct(data product_schemas.GetProduct) error
}
