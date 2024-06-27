package item_schemas

import "time"

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
	PersonID   uint      `json:"person_id" format:"id"`
}

type AddItem struct {
	UserID     uint    `json:"user_id" format:"id"`
	ProductID  uint    `json:"product_id" format:"id"`
	ItemCost   float32 `json:"item_cost" format:"item_cost"`
	ItemAmount float32 `json:"item_amount" format:"item_amount"`
	ItemType   uint8   `json:"item_type" format:"item_type"`
	PersonID   uint    `json:"person_id" format:"id"`
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
	PersonID   uint      `json:"person_id" format:"id"`
	// Parsed info
	ProductTitle    string `json:"product_title" format:"product_title"`
	ProductCalories uint   `json:"product_calories" format:"product_calories"`
	ProductFats     uint   `json:"product_fats" format:"product_nutrient"`
	ProductCarbs    uint   `json:"product_carbs" format:"product_nutrient"`
	ProductProteins uint   `json:"product_proteins" format:"product_nutrient"`
	PersonName      string `json:"person_name" format:"username"`
}
