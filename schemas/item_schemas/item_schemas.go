package item_schemas

import (
	"time"

	"github.com/bmg-c/product-diary/schemas/user_schemas"
)

const (
	// Purchase made by the user
	ItemTypeMyPurchase uint8 = iota + 1
	// Purchase made by some user person to user
	ItemTypeFromPersonPurchase
	// Purchase made by user to some person
	ItemTypeToPersonPurchase
)

type ItemDB struct {
	ItemID     uint      `json:"item_id" format:"id"`
	UserID     uint      `json:"user_id" format:"id"`
	ProductID  uint      `json:"product_id" format:"id"`
	ItemDate   time.Time `json:"item_date"`
	ItemCost   float32   `json:"item_cost" format:"item_cost"`
	ItemAmount float32   `json:"item_amount" format:"item_amount"`
	ItemType   uint8     `json:"item_type" format:"item_type"`
	PersonID   uint      `json:"person_id" format:"id" validate:"omitzero"`
}

type AddItem struct {
	UserID    uint      `json:"user_id" format:"id"`
	ProductID uint      `json:"product_id" format:"id"`
	ItemDate  time.Time `json:"item_date"`
	// ItemCost   float32 `json:"item_cost" format:"item_cost"`
	// ItemAmount float32 `json:"item_amount" format:"item_amount"`
	// ItemType   uint8   `json:"item_type" format:"item_type"`
	// PersonID   uint    `json:"person_id" format:"id"`
}

type DeleteItem struct {
	ItemID uint `json:"item_id" format:"id"`
	UserID uint `json:"user_id" format:"id"`
}

type ChangeItem struct {
	ItemID     uint    `json:"item_id" format:"id"`
	UserID     uint    `json:"user_id" format:"id"`
	ProductID  uint    `json:"product_id" format:"id" validate:"omitzero"`
	ItemCost   float32 `json:"item_cost" format:"item_cost" validate:"omitzero"`
	ItemAmount float32 `json:"item_amount" format:"item_amount" validate:"omitzero"`
	ItemType   uint8   `json:"item_type" format:"item_type" validate:"omitzero"`
	PersonID   uint    `json:"person_id" format:"id" validate:"omitzero"`
}

type GetItems struct {
	UserID      uint      `json:"user_id" format:"id"`
	ItemDate    time.Time `json:"item_date"`
	SearchQuery string    `json:"search_query"`
}

type GetItem struct {
	ItemID uint `json:"item_id" format:"id"`
	UserID uint `json:"user_id" format:"id"`
}

type ItemParsed struct {
	ItemID     uint      `json:"item_id" format:"id"`
	UserID     uint      `json:"user_id" format:"id"`
	ProductID  uint      `json:"product_id" format:"id"`
	ItemDate   time.Time `json:"item_date"`
	ItemCost   float32   `json:"item_cost" format:"item_cost"`
	ItemAmount float32   `json:"item_amount" format:"item_amount"`
	ItemType   uint8     `json:"item_type" format:"item_type"`
	PersonID   uint      `json:"person_id" format:"id" validate:"omitzero"`
	// Parsed info
	ProductTitle    string  `json:"product_title" format:"product_title"`
	ProductCalories float32 `json:"product_calories" format:"product_calories"`
	ProductFats     float32 `json:"product_fats" format:"product_nutrient"`
	ProductCarbs    float32 `json:"product_carbs" format:"product_nutrient"`
	ProductProteins float32 `json:"product_proteins" format:"product_nutrient"`
	PersonName      string  `json:"person_name" format:"username" validate:"omitzero"`
	PersonIsHidden  bool    `json:"person_is_hidden"`
}

type GetItemsRange struct {
	UserID       uint      `json:"user_id" format:"id"`
	ItemDateFrom time.Time `json:"item_date_from"`
	ItemDateTo   time.Time `json:"item_date_to"`
}

type PersonAnalytics struct {
	PersonDB  user_schemas.PersonDB `json:"person_db"`
	TotalDebt float32               `json:"total_debt"`
}

type Analytics struct {
	TotalSpent    float32           `json:"total_spent"`
	Persons       []PersonAnalytics `json:"persons"`
	TotalCalories float32           `json:"total_calories"`
	TotalFats     float32           `json:"total_fats"`
	TotalCarbs    float32           `json:"total_carbs"`
	TotalProteins float32           `json:"total_preteins"`
}
