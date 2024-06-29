package services

import (
	"errors"

	E "github.com/bmg-c/product-diary/errorhandler"
	"github.com/bmg-c/product-diary/schemas/item_schemas"
	"github.com/bmg-c/product-diary/schemas/user_schemas"
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
	GetItemsRange(data item_schemas.GetItemsRange) ([]item_schemas.ItemParsed, error)
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

func (is *ItemService) GetAnalyticsRange(data item_schemas.GetItemsRange) (item_schemas.Analytics, error) {
	items, err := is.itemDB.GetItemsRange(data)
	if err != nil {
		return item_schemas.Analytics{}, err
	}

	a := item_schemas.Analytics{
		TotalSpent:    0,
		TotalCalories: 0,
		TotalFats:     0,
		TotalCarbs:    0,
		TotalProteins: 0,
		Persons:       []item_schemas.PersonAnalytics{},
	}
	for _, i := range items {
		switch i.ItemType {
		case item_schemas.ItemTypeMyPurchase:
			a.TotalSpent += i.ItemCost * i.ItemAmount
			a.TotalCalories += float32(i.ProductCalories) * i.ItemAmount
			a.TotalFats += float32(i.ProductFats) * i.ItemAmount
			a.TotalCarbs += float32(i.ProductCarbs) * i.ItemAmount
			a.TotalProteins += float32(i.ProductProteins) * i.ItemAmount
		case item_schemas.ItemTypeFromPersonPurchase:
			a.TotalSpent += i.ItemCost * i.ItemAmount
			a.TotalCalories += float32(i.ProductCalories) * i.ItemAmount
			a.TotalFats += float32(i.ProductFats) * i.ItemAmount
			a.TotalCarbs += float32(i.ProductCarbs) * i.ItemAmount
			a.TotalProteins += float32(i.ProductProteins) * i.ItemAmount
			personInd := -1
			for ind, personDB := range a.Persons {
				if personDB.PersonDB.PersonID == i.PersonID {
					personInd = ind
					break
				}
			}
			if personInd != -1 {
				a.Persons[personInd].TotalDebt += i.ItemCost * i.ItemAmount
			} else {
				a.Persons = append(a.Persons, item_schemas.PersonAnalytics{
					PersonDB: user_schemas.PersonDB{
						PersonID:   i.PersonID,
						UserID:     i.UserID,
						PersonName: i.PersonName,
					},
					TotalDebt: i.ItemCost * i.ItemAmount,
				})
			}
		case item_schemas.ItemTypeToPersonPurchase:
			personInd := -1
			for ind, personDB := range a.Persons {
				if personDB.PersonDB.PersonID == i.PersonID {
					personInd = ind
					break
				}
			}
			if personInd != -1 {
				a.Persons[personInd].TotalDebt -= i.ItemCost * i.ItemAmount
			} else {
				a.Persons = append(a.Persons, item_schemas.PersonAnalytics{
					PersonDB: user_schemas.PersonDB{
						PersonID:   i.PersonID,
						UserID:     i.UserID,
						PersonName: i.PersonName,
					},
					TotalDebt: -(i.ItemCost * i.ItemAmount),
				})
			}
		default:
			// Error in db values?
			return item_schemas.Analytics{}, E.ErrInternalServer
		}
	}
	return a, nil
}

func (is *ItemService) GetAnalytics(data []item_schemas.ItemParsed) (item_schemas.Analytics, error) {
	a := item_schemas.Analytics{
		TotalSpent:    0,
		TotalCalories: 0,
		TotalFats:     0,
		TotalCarbs:    0,
		TotalProteins: 0,
		Persons:       []item_schemas.PersonAnalytics{},
	}
	for _, i := range data {
		switch i.ItemType {
		case item_schemas.ItemTypeMyPurchase:
			a.TotalSpent += i.ItemCost * i.ItemAmount
			a.TotalCalories += float32(i.ProductCalories) * i.ItemAmount
			a.TotalFats += float32(i.ProductFats) * i.ItemAmount
			a.TotalCarbs += float32(i.ProductCarbs) * i.ItemAmount
			a.TotalProteins += float32(i.ProductProteins) * i.ItemAmount
		case item_schemas.ItemTypeFromPersonPurchase:
			a.TotalSpent += i.ItemCost * i.ItemAmount
			a.TotalCalories += float32(i.ProductCalories) * i.ItemAmount
			a.TotalFats += float32(i.ProductFats) * i.ItemAmount
			a.TotalCarbs += float32(i.ProductCarbs) * i.ItemAmount
			a.TotalProteins += float32(i.ProductProteins) * i.ItemAmount
			personInd := -1
			for ind, personDB := range a.Persons {
				if personDB.PersonDB.PersonID == i.PersonID {
					personInd = ind
					break
				}
			}
			if personInd != -1 {
				a.Persons[personInd].TotalDebt += i.ItemCost * i.ItemAmount
			} else {
				a.Persons = append(a.Persons, item_schemas.PersonAnalytics{
					PersonDB: user_schemas.PersonDB{
						PersonID:   i.PersonID,
						UserID:     i.UserID,
						PersonName: i.PersonName,
					},
					TotalDebt: i.ItemCost * i.ItemAmount,
				})
			}
		case item_schemas.ItemTypeToPersonPurchase:
			personInd := -1
			for ind, personDB := range a.Persons {
				if personDB.PersonDB.PersonID == i.PersonID {
					personInd = ind
					break
				}
			}
			if personInd != -1 {
				a.Persons[personInd].TotalDebt -= i.ItemCost * i.ItemAmount
			} else {
				a.Persons = append(a.Persons, item_schemas.PersonAnalytics{
					PersonDB: user_schemas.PersonDB{
						PersonID:   i.PersonID,
						UserID:     i.UserID,
						PersonName: i.PersonName,
					},
					TotalDebt: -(i.ItemCost * i.ItemAmount),
				})
			}
		default:
			// Error in db values?
			return item_schemas.Analytics{}, E.ErrInternalServer
		}
	}
	return a, nil
}
