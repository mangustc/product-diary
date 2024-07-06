package product_db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/bmg-c/product-diary/db"
	E "github.com/bmg-c/product-diary/errorhandler"
	"github.com/bmg-c/product-diary/schemas/product_schemas"
	"github.com/bmg-c/product-diary/util"
	"github.com/mattn/go-sqlite3"
)

type ProductDB struct {
	productStore *db.Store
}

func NewProductDB(productStore *db.Store) (*ProductDB, error) {
	if productStore == nil {
		return nil, fmt.Errorf("Error creating ProductDB instance, one of the stores is nil")
	}
	return &ProductDB{
		productStore: productStore,
	}, nil
}

func (pdb *ProductDB) AddProduct(data product_schemas.AddProduct) (product_schemas.ProductDB, error) {
	query := `INSERT INTO ` + pdb.productStore.TableName + `
        (product_id, product_title, product_calories, product_fats, product_carbs, product_proteins, user_id, is_deleted)
        VALUES (NULL, ?, ?, ?, ?, ?, ?, FALSE)`

	stmt, err := pdb.productStore.DB.Prepare(query)
	defer stmt.Close()
	if err != nil {
		return product_schemas.ProductDB{}, E.ErrInternalServer
	}
	res, err := stmt.Exec(
		data.ProductTitle,
		data.ProductCalories,
		data.ProductFats,
		data.ProductCarbs,
		data.ProductProteins,
		data.UserID,
	)
	if err != nil {
		if util.IsErrorSQL(err, sqlite3.ErrConstraint) {
			return product_schemas.ProductDB{}, E.ErrUnprocessableEntity
		}
		return product_schemas.ProductDB{}, E.ErrInternalServer
	}
	createdID, err := res.LastInsertId()
	if err != nil {
		return product_schemas.ProductDB{}, E.ErrInternalServer
	}
	productDB := product_schemas.ProductDB{
		ProductID:       uint(createdID),
		ProductTitle:    data.ProductTitle,
		ProductCalories: data.ProductCalories,
		ProductFats:     data.ProductFats,
		ProductCarbs:    data.ProductCarbs,
		ProductProteins: data.ProductProteins,
		UserID:          data.UserID,
		IsDeleted:       false,
	}

	return productDB, nil
}

func (pdb *ProductDB) GetProducts(data product_schemas.GetProducts) ([]product_schemas.ProductDB, error) {
	var productDB product_schemas.ProductDB = product_schemas.ProductDB{}
	query := `SELECT product_id, product_title, product_calories, product_fats, product_carbs, product_proteins, user_id, is_deleted
        FROM ` + pdb.productStore.TableName + `
        WHERE length(trim(replace(lower(?), ' ', ''), replace(lower(product_title || product_calories 
    || product_fats || product_carbs || product_proteins), ' ', ''))) < 1 AND
            is_deleted = FALSE`

	rows, err := pdb.productStore.DB.Query(query, data.SearchQuery)
	if err != nil {
		fmt.Printf("%v\n", err)
		if errors.Is(err, sql.ErrNoRows) {
			return []product_schemas.ProductDB{}, nil
		}
		return []product_schemas.ProductDB{}, E.ErrInternalServer
	}
	defer rows.Close()

	products := []product_schemas.ProductDB{}
	for rows.Next() {
		err = rows.Scan(
			&productDB.ProductID,
			&productDB.ProductTitle,
			&productDB.ProductCalories,
			&productDB.ProductFats,
			&productDB.ProductCarbs,
			&productDB.ProductProteins,
			&productDB.UserID,
			&productDB.IsDeleted,
		)
		if err != nil {
			return []product_schemas.ProductDB{}, E.ErrInternalServer
		}
		products = append(products, productDB)
	}

	return products, nil
}

func (pdb *ProductDB) GetProduct(data product_schemas.GetProduct) (product_schemas.ProductDB, error) {
	var productDB product_schemas.ProductDB = product_schemas.ProductDB{}
	query := `SELECT product_id, product_title, product_calories, product_fats,
        product_carbs, product_proteins, user_id, is_deleted FROM ` + pdb.productStore.TableName + `
		WHERE product_id = ? AND is_deleted = FALSE`

	stmt, err := pdb.productStore.DB.Prepare(query)
	if err != nil {
		return product_schemas.ProductDB{}, E.ErrInternalServer
	}
	defer stmt.Close()

	err = stmt.QueryRow(
		data.ProductID,
	).Scan(
		&productDB.ProductID,
		&productDB.ProductTitle,
		&productDB.ProductCalories,
		&productDB.ProductFats,
		&productDB.ProductCarbs,
		&productDB.ProductProteins,
		&productDB.UserID,
		&productDB.IsDeleted,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return product_schemas.ProductDB{}, E.ErrNotFound
		}
		return product_schemas.ProductDB{}, E.ErrInternalServer
	}

	return productDB, nil
}

func (pdb *ProductDB) DeleteProduct(data product_schemas.DeleteProduct) error {
	query := `UPDATE ` + pdb.productStore.TableName + ` 
        SET is_deleted = TRUE
        WHERE product_id = ? AND user_id = ?`

	stmt, err := pdb.productStore.DB.Prepare(query)
	defer stmt.Close()
	if err != nil {
		return E.ErrInternalServer
	}
	_, err = stmt.Exec(data.ProductID, data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return E.ErrNotFound
		}
		return E.ErrInternalServer
	}

	return nil
}
