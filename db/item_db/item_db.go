package item_db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/bmg-c/product-diary/db"
	E "github.com/bmg-c/product-diary/errorhandler"
	"github.com/bmg-c/product-diary/logger"
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
	cols := []string{}
	argsStr := []string{}
	args := []any{}

	cols = append(cols, "user_id")
	args = append(args, data.UserID)
	argsStr = append(argsStr, "?")
	cols = append(cols, "product_id")
	args = append(args, data.ProductID)
	argsStr = append(argsStr, "?")
	cols = append(cols, "item_date")
	args = append(args, data.ItemDate.Format("2006-01-02"))
	argsStr = append(argsStr, "?")
	if !schemas.IsZero(data.ItemCost) {
		cols = append(cols, "item_cost")
		args = append(args, data.ItemCost)
		argsStr = append(argsStr, "?")
	}
	if !schemas.IsZero(data.ItemAmount) {
		cols = append(cols, "item_amount")
		args = append(args, data.ItemAmount)
		argsStr = append(argsStr, "?")
	}
	if !schemas.IsZero(data.ItemType) {
		cols = append(cols, "item_type")
		args = append(args, data.ItemType)
		argsStr = append(argsStr, "?")
	}
	if !schemas.IsZero(data.PersonID) {
		cols = append(cols, "person_id")
		args = append(args, data.PersonID)
		argsStr = append(argsStr, "?")
	}

	query := `INSERT INTO ` + idb.itemStore.TableName + `
        (` + strings.Join(cols, ", ") + `)
        VALUES (` + strings.Join(argsStr, ", ") + `)
        RETURNING *`

	stmt, err := idb.itemStore.DB.Prepare(query)
	defer stmt.Close()
	if err != nil {
		return item_schemas.ItemDB{}, E.ErrInternalServer
	}

	nullPersonID := sql.NullInt64{}
	itemDB := item_schemas.ItemDB{}
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
		&nullPersonID,
	)
	if err != nil {
		if util.IsErrorSQL(err, sqlite3.ErrConstraint) {
			return item_schemas.ItemDB{}, E.ErrUnprocessableEntity
		}
		return item_schemas.ItemDB{}, E.ErrInternalServer
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
            INNER JOIN %[3]s ON %[1]s.user_id = %[3]s.user_id)
        WHERE (
            (%[1]s.user_id = ? AND %[1]s.item_date = ?) AND 
            (length(trim(
                replace(lower(?), ' ', ''),
                replace(lower(
                    %[2]s.product_title ||
                    %[2]s.product_calories ||
                    %[2]s.product_fats ||
                    %[2]s.product_carbs ||
                    %[2]s.product_proteins), ' ', '')
            )) < 3))
        GROUP BY %[1]s.item_id`,
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

	personIDNull := sql.NullInt64{}
	personNameNull := sql.NullString{}
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
			&personIDNull,
			&itemParsed.ProductTitle,
			&itemParsed.ProductCalories,
			&itemParsed.ProductFats,
			&itemParsed.ProductCarbs,
			&itemParsed.ProductProteins,
			&personNameNull,
		)
		if personIDNull.Valid {
			itemParsed.PersonID = uint(personIDNull.Int64)
			itemParsed.PersonName = personNameNull.String
		} else {
			itemParsed.PersonID = 0
			itemParsed.PersonName = ""
		}
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
            INNER JOIN %[3]s ON %[1]s.user_id = %[3]s.user_id)
        WHERE 
            (%[1]s.item_id = ? AND %[1]s.user_id = ?)
        GROUP BY %[1]s.item_id`,
		idb.itemStore.TableName,
		idb.productStore.TableName,
		idb.personStore.TableName,
	)

	stmt, err := idb.itemStore.DB.Prepare(query)
	if err != nil {
		return item_schemas.ItemParsed{}, E.ErrInternalServer
	}
	defer stmt.Close()

	personIDNull := sql.NullInt64{}
	personNameNull := sql.NullString{}
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
		&personIDNull,
		&itemParsed.ProductTitle,
		&itemParsed.ProductCalories,
		&itemParsed.ProductFats,
		&itemParsed.ProductCarbs,
		&itemParsed.ProductProteins,
		&personNameNull,
	)
	if personIDNull.Valid {
		itemParsed.PersonID = uint(personIDNull.Int64)
		itemParsed.PersonName = personNameNull.String
	} else {
		itemParsed.PersonID = 0
		itemParsed.PersonName = ""
	}
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
	setOptions := []string{}
	args := []any{}
	if !schemas.IsZero(data.ProductID) {
		setOptions = append(setOptions, "product_id = ?")
		args = append(args, data.ProductID)
	}
	if !schemas.IsZero(data.ItemCost) {
		setOptions = append(setOptions, "item_cost = ?")
		args = append(args, data.ItemCost)
	}
	if !schemas.IsZero(data.ItemAmount) {
		setOptions = append(setOptions, "item_amount = ?")
		args = append(args, data.ItemAmount)
	}
	if !schemas.IsZero(data.ItemType) {
		setOptions = append(setOptions, "item_type = ?")
		args = append(args, data.ItemType)
	}
	if !schemas.IsZero(data.PersonID) {
		setOptions = append(setOptions, "person_id = ?")
		args = append(args, data.PersonID)
	}
	query := `UPDATE ` + idb.itemStore.TableName + "\nSET " +
		strings.Join(setOptions, ", ") + `
        WHERE item_id = ? AND user_id = ?
        RETURNING *`
	args = append(args, data.ItemID, data.UserID)
	// logger.Info.Println(query)

	stmt, err := idb.itemStore.DB.Prepare(query)
	if err != nil {
		return item_schemas.ItemDB{}, E.ErrInternalServer
	}
	defer stmt.Close()

	personIDNull := sql.NullInt64{}
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
		&personIDNull,
	)
	if personIDNull.Valid {
		itemDB.PersonID = uint(personIDNull.Int64)
	}
	if err != nil {
		logger.Info.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return item_schemas.ItemDB{}, E.ErrNotFound
		}
		if util.IsErrorSQL(err, sqlite3.ErrConstraint) {
			return item_schemas.ItemDB{}, E.ErrUnprocessableEntity
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
	_, err = stmt.Exec(data.ItemID, data.UserID)
	if err != nil {
		if util.IsErrorSQL(err, sqlite3.ErrNotFound) {
			return E.ErrNotFound
		}
		return E.ErrInternalServer
	}

	return nil
}

func (idb *ItemDB) GetItemsRange(data item_schemas.GetItemsRange) ([]item_schemas.ItemParsed, error) {
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
            INNER JOIN %[3]s ON %[1]s.user_id = %[3]s.user_id)
        WHERE 
            (%[1]s.user_id = ? AND (%[1]s.item_date >= ? AND %[1]s.item_date <= ?))
        GROUP BY %[1]s.item_id`,
		idb.itemStore.TableName,
		idb.productStore.TableName,
		idb.personStore.TableName,
	)

	rows, err := idb.itemStore.DB.Query(query,
		data.UserID,
		data.ItemDateFrom.Format("2006-01-02"),
		data.ItemDateTo.Format("2006-01-02"),
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []item_schemas.ItemParsed{}, nil
		}
		return []item_schemas.ItemParsed{}, E.ErrInternalServer
	}
	defer rows.Close()

	personIDNull := sql.NullInt64{}
	personNameNull := sql.NullString{}
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
			&personIDNull,
			&itemParsed.ProductTitle,
			&itemParsed.ProductCalories,
			&itemParsed.ProductFats,
			&itemParsed.ProductCarbs,
			&itemParsed.ProductProteins,
			&personNameNull,
		)
		if personIDNull.Valid {
			itemParsed.PersonID = uint(personIDNull.Int64)
			itemParsed.PersonName = personNameNull.String
		} else {
			itemParsed.PersonID = 0
			itemParsed.PersonName = ""
		}
		if err != nil {
			return []item_schemas.ItemParsed{}, E.ErrInternalServer
		}
		items = append(items, itemParsed)
	}

	return items, nil
}
