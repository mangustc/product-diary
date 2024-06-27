package item_db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/bmg-c/product-diary/db"
	E "github.com/bmg-c/product-diary/errorhandler"
	"github.com/bmg-c/product-diary/schemas"
	"github.com/bmg-c/product-diary/schemas/item_schemas"
	"github.com/bmg-c/product-diary/util"
	"github.com/mattn/go-sqlite3"
)

type ItemDB struct {
	itemStore    *db.Store
	productStore *db.Store
	personStore  *db.Store
}

func NewItemDB(itemStore *db.Store, productStore *db.Store, personStore *db.Store) (*ItemDB, error) {
	if itemStore == nil || productStore == nil || personStore == nil {
		return nil, fmt.Errorf("Error creating ItemDB instance, one of the stores is nil")
	}
	return &ItemDB{
		itemStore:    itemStore,
		productStore: productStore,
		personStore:  personStore,
	}, nil
}

func (idb *ItemDB) AddItem(data item_schemas.AddItem) (item_schemas.ItemDB, error) {
	query := `INSERT INTO ` + idb.itemStore.TableName + `
        (item_id, user_id, product_id, item_date, item_cost, item_amount, item_type, person_id)
        VALUES (NULL, ?, ?, date('now'), ?, ?, ?, ?)`
	stmt, err := idb.itemStore.DB.Prepare(query)
	defer stmt.Close()
	if err != nil {
		return item_schemas.ItemDB{}, E.ErrInternalServer
	}
	res, err := stmt.Exec(
		&data.UserID,
		&data.ProductID,
		&data.ItemCost,
		&data.ItemAmount,
		&data.ItemType,
		&data.PersonID,
	)
	if err != nil {
		if util.IsErrorSQL(err, sqlite3.ErrConstraint) {
			return item_schemas.ItemDB{}, E.ErrUnprocessableEntity
		}
		return item_schemas.ItemDB{}, E.ErrInternalServer
	}
	createdID, err := res.LastInsertId()
	if err != nil {
		return item_schemas.ItemDB{}, E.ErrInternalServer
	}
	itemDB := item_schemas.ItemDB{
		ItemID:     uint(createdID),
		UserID:     data.UserID,
		ProductID:  data.PersonID,
		ItemDate:   time.Now(),
		ItemCost:   data.ItemCost,
		ItemAmount: data.ItemAmount,
		ItemType:   data.ItemType,
		PersonID:   data.PersonID,
	}

	return itemDB, nil
}

func (idb *ItemDB) GetItems(data item_schemas.GetItems) ([]item_schemas.ItemParsed, error) {
	var itemParsed item_schemas.ItemParsed = item_schemas.ItemParsed{}
	query := fmt.Sprintf(`
        SELECT
            %[1]s.item_id,
            %[1]s.user_id,
            %[1]s.product_id,
            %[1]s.item_date,
            %[1]s.item_cost,
            %[1]s.item_amount,
            %[1]s.item_type,
            %[1]s.person_id,
            %[2]s.product_title,
            %[2]s.product_calories,
            %[2]s.product_fats,
            %[2]s.product_carbs,
            %[2]s.product_proteins,
            %[3]s.person_name
        FROM ((%[1]s
            INNER JOIN %[2]s ON %[1]s.product_id = %[2]s.product_id) 
            INNER JOIN %[3]s ON %[3]s.person_id = %[3]s.person_id)
        WHERE (
            (%[1]s.user_id = ? AND %[1]s.date = ?) AND 
            (length(trim(
                replace(lower(?), ' ', ''),
                replace(lower(
                    product_title ||
                    product_calories ||
                    product_fats ||
                    product_carbs ||
                    product_proteins), ' ', '')
            )) < 3)`,
		idb.itemStore.TableName,
		idb.productStore.TableName,
		idb.personStore.TableName,
	)

	rows, err := idb.itemStore.DB.Query(query, data.UserID, data.ItemDate.Format("2006-01-02"), data.SearchQuery)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []item_schemas.ItemParsed{}, nil
		}
		return []item_schemas.ItemParsed{}, E.ErrInternalServer
	}
	defer rows.Close()

	items := []item_schemas.ItemParsed{}
	for rows.Next() {
		err = rows.Scan(
			&itemParsed.ItemID,
			&itemParsed.UserID,
			&itemParsed.ProductID,
			&itemParsed.ItemDate,
			&itemParsed.ItemCost,
			&itemParsed.ItemAmount,
			&itemParsed.ItemType,
			&itemParsed.PersonID,
			&itemParsed.ProductTitle,
			&itemParsed.ProductCalories,
			&itemParsed.ProductFats,
			&itemParsed.ProductCarbs,
			&itemParsed.ProductProteins,
			&itemParsed.PersonName,
		)
		if err != nil {
			return []item_schemas.ItemParsed{}, E.ErrInternalServer
		}
		items = append(items, itemParsed)
	}

	return items, nil
}

func (idb *ItemDB) GetItem(data item_schemas.GetItem) (item_schemas.ItemParsed, error) {
	var itemParsed item_schemas.ItemParsed = item_schemas.ItemParsed{}
	query := fmt.Sprintf(`
        SELECT
            %[1]s.item_id,
            %[1]s.user_id,
            %[1]s.product_id,
            %[1]s.item_date,
            %[1]s.item_cost,
            %[1]s.item_amount,
            %[1]s.item_type,
            %[1]s.person_id,
            %[2]s.product_title,
            %[2]s.product_calories,
            %[2]s.product_fats,
            %[2]s.product_carbs,
            %[2]s.product_proteins,
            %[3]s.person_name
        FROM ((%[1]s
            INNER JOIN %[2]s ON %[1]s.product_id = %[2]s.product_id) 
            INNER JOIN %[3]s ON %[3]s.person_id = %[3]s.person_id)
        WHERE (
            (%[1]s.item_id = ? AND %[1]s.user_id = ?)`,
		idb.itemStore.TableName,
		idb.productStore.TableName,
		idb.personStore.TableName,
	)

	stmt, err := idb.itemStore.DB.Prepare(query)
	if err != nil {
		return item_schemas.ItemParsed{}, E.ErrInternalServer
	}
	defer stmt.Close()

	err = stmt.QueryRow(
		data.ItemID,
		data.UserID,
	).Scan(
		&itemParsed.ItemID,
		&itemParsed.UserID,
		&itemParsed.ProductID,
		&itemParsed.ItemDate,
		&itemParsed.ItemCost,
		&itemParsed.ItemAmount,
		&itemParsed.ItemType,
		&itemParsed.PersonID,
		&itemParsed.ProductTitle,
		&itemParsed.ProductCalories,
		&itemParsed.ProductFats,
		&itemParsed.ProductCarbs,
		&itemParsed.ProductProteins,
		&itemParsed.PersonName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return item_schemas.ItemParsed{}, E.ErrNotFound
		}
		return item_schemas.ItemParsed{}, E.ErrInternalServer
	}

	return itemParsed, nil
}

func (idb *ItemDB) ChangeItem(data item_schemas.ChangeItem) (item_schemas.ItemDB, error) {
	var itemDB item_schemas.ItemDB = item_schemas.ItemDB{}
	setOptions := ""
	args := []any{}
	if !schemas.IsZero(data.ProductID) {
		setOptions += "SET product_id = ?\n"
		args = append(args, data.ProductID)
	}
	if !schemas.IsZero(data.ItemCost) {
		setOptions += "SET item_cost = ?\n"
		args = append(args, data.ItemCost)
	}
	if !schemas.IsZero(data.ItemAmount) {
		setOptions += "SET item_amount = ?\n"
		args = append(args, data.ItemAmount)
	}
	if !schemas.IsZero(data.ItemType) {
		setOptions += "SET item_type = ?\n"
		args = append(args, data.ItemType)
	}
	if !schemas.IsZero(data.PersonID) {
		setOptions += "SET person_id = ?\n"
		args = append(args, data.PersonID)
	}
	query := `UPDATE ` + idb.itemStore.TableName +
		setOptions + `
        WHERE item_id = ? AND user_id = ?
        RETURNING *`
	args = append(args, data.ItemID, data.UserID)

	stmt, err := idb.itemStore.DB.Prepare(query)
	if err != nil {
		return item_schemas.ItemDB{}, E.ErrInternalServer
	}
	defer stmt.Close()

	err = stmt.QueryRow(
		args...,
	).Scan(
		&itemDB.ItemID,
		&itemDB.UserID,
		&itemDB.ProductID,
		&itemDB.ItemDate,
		&itemDB.ItemCost,
		&itemDB.ItemAmount,
		&itemDB.ItemType,
		&itemDB.PersonID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return item_schemas.ItemDB{}, E.ErrNotFound
		}
		return item_schemas.ItemDB{}, E.ErrInternalServer
	}

	return itemDB, nil
}

func (idb *ItemDB) DeleteItem(data item_schemas.DeleteItem) error {
	query := `DELETE FROM ` + idb.itemStore.TableName + ` 
        WHERE item_id = ? AND user_id = ?`

	stmt, err := idb.itemStore.DB.Prepare(query)
	defer stmt.Close()
	if err != nil {
		return E.ErrInternalServer
	}
	_, err = stmt.Exec()
	if err != nil {
		if util.IsErrorSQL(err, sqlite3.ErrNotFound) {
			return E.ErrNotFound
		}
		return E.ErrInternalServer
	}

	return nil
}
