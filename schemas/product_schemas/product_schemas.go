package product_schemas

type ProductDB struct {
	ProductID       uint    `json:"product_id" format:"id"`
	ProductTitle    string  `json:"product_title" format:"product_title"`
	ProductCalories float32 `json:"product_calories" format:"product_calories"`
	ProductFats     float32 `json:"product_fats" format:"product_nutrient"`
	ProductCarbs    float32 `json:"product_carbs" format:"product_nutrient"`
	ProductProteins float32 `json:"product_proteins" format:"product_nutrient"`
	UserID          uint    `json:"user_id" format:"id"`
	IsDeleted       bool    `json:"is_deleted"`
}

type AddProduct struct {
	ProductTitle    string  `json:"product_title" format:"product_title"`
	ProductCalories float32 `json:"product_calories" format:"product_calories" validate:"omitzero"`
	ProductFats     float32 `json:"product_fats" format:"product_nutrient" validate:"omitzero"`
	ProductCarbs    float32 `json:"product_carbs" format:"product_nutrient" validate:"omitzero"`
	ProductProteins float32 `json:"product_proteins" format:"product_nutrient" validate:"omitzero"`
	UserID          uint    `json:"user_id" format:"id"`
}

type GetProduct struct {
	ProductID uint `json:"product_id" format:"id"`
}

type DeleteProduct struct {
	ProductID uint `json:"product_id" format:"id"`
	UserID    uint `json:"user_id" format:"id"`
}

type GetProducts struct {
	SearchQuery string `json:"search_query"`
}
