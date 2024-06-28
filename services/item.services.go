package services

import (
	"errors"

	E "github.com/bmg-c/product-diary/errorhandler"
	"github.com/bmg-c/product-diary/schemas/item_schemas"
)

func NewItemService(itemDB ItemDB) *ItemService {
	return &ItemService{
		itemDB: itemDB,
	}
}

type ItemService struct {
	itemDB ItemDB
}

type ItemDB interface {
	AddItem(data item_schemas.AddItem) (item_schemas.ItemDB, error)
	DeleteItem(data item_schemas.DeleteItem) error
	GetItem(data item_schemas.GetItem) (item_schemas.ItemParsed, error)
	GetItems(data item_schemas.GetItems) ([]item_schemas.ItemParsed, error)
	ChangeItem(data item_schemas.ChangeItem) (item_schemas.ItemDB, error)
}

func (is *ItemService) AddItem(data item_schemas.AddItem) (item_schemas.ItemParsed, error) {
	itemDB, err := is.itemDB.AddItem(data)
	if err != nil {
		return item_schemas.ItemParsed{}, err
	}

	getItem := item_schemas.GetItem{
		ItemID: itemDB.ItemID,
		UserID: itemDB.UserID,
	}
	itemParsed, err := is.itemDB.GetItem(getItem)
	if err != nil {
		if errors.Is(err, E.ErrNotFound) {
			return item_schemas.ItemParsed{}, E.ErrInternalServer
		}
		return item_schemas.ItemParsed{}, err
	}
	return itemParsed, nil
}

func (is *ItemService) DeleteItem(data item_schemas.DeleteItem) error {
	err := is.itemDB.DeleteItem(data)
	if err != nil {
		if errors.Is(err, E.ErrNotFound) {
			return E.ErrUnprocessableEntity
		}
		return err
	}
	return nil
}

func (is *ItemService) GetItems(data item_schemas.GetItems) ([]item_schemas.ItemParsed, error) {
	items, err := is.itemDB.GetItems(data)
	if err != nil {
		return []item_schemas.ItemParsed{}, err
	}

	return items, nil
}

func (is *ItemService) ChangeItem(data item_schemas.ChangeItem) (item_schemas.ItemParsed, error) {
	itemDB, err := is.itemDB.ChangeItem(data)
	if err != nil {
		if errors.Is(err, E.ErrNotFound) {
			return item_schemas.ItemParsed{}, E.ErrUnprocessableEntity
		}
		return item_schemas.ItemParsed{}, err
	}

	getItem := item_schemas.GetItem{
		ItemID: itemDB.ItemID,
		UserID: itemDB.UserID,
	}
	itemParsed, err := is.itemDB.GetItem(getItem)
	if err != nil {
		if errors.Is(err, E.ErrNotFound) {
			return item_schemas.ItemParsed{}, E.ErrInternalServer
		}
		return item_schemas.ItemParsed{}, err
	}
	return itemParsed, nil
}
