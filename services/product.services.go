package services

import (
	"errors"

	E "github.com/bmg-c/product-diary/errorhandler"
	"github.com/bmg-c/product-diary/schemas/product_schemas"
)

func NewProductService(productDB ProductDB) *ProductService {
	return &ProductService{
		productDB: productDB,
	}
}

type ProductService struct {
	productDB ProductDB
}

type ProductDB interface {
	AddProduct(data product_schemas.AddProduct) (product_schemas.ProductDB, error)
	GetProducts(data product_schemas.GetProducts) ([]product_schemas.ProductDB, error)
	GetProduct(data product_schemas.GetProduct) (product_schemas.ProductDB, error)
	DeleteProduct(data product_schemas.DeleteProduct) error
}

func (ps *ProductService) AddProduct(data product_schemas.AddProduct) (product_schemas.ProductDB, error) {
	productDB, err := ps.productDB.AddProduct(data)
	if err != nil {
		return product_schemas.ProductDB{}, err
	}

	return productDB, nil
}

func (ps *ProductService) GetProducts(data product_schemas.GetProducts) ([]product_schemas.ProductDB, error) {
	products, err := ps.productDB.GetProducts(data)
	if err != nil {
		return []product_schemas.ProductDB{}, err
	}

	return products, nil
}

func (ps *ProductService) GetProduct(data product_schemas.GetProduct) (product_schemas.ProductDB, error) {
	productDB, err := ps.productDB.GetProduct(data)
	if err != nil {
		return product_schemas.ProductDB{}, err
	}

	return productDB, nil
}

func (ps *ProductService) DeleteProduct(data product_schemas.DeleteProduct) error {
	err := ps.productDB.DeleteProduct(data)
	if err != nil {
		if errors.Is(err, E.ErrNotFound) {
			return E.ErrUnprocessableEntity
		}
		return err
	}

	return nil
}
