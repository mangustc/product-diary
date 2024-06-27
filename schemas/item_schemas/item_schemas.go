package item_schemas

const (
	// Purchase made by the user
	ItemTypeMyPurchase uint8 = iota + 1
	// Purchase made by some user person to user
	ItemTypeFromPersonPurchase
	// Purchase made by user to some person
	ItemTypeToPersonPurchase
)

type ItemDB struct {
	ItemID     uint    `json:"item_id" format:"id"`
	UserID     uint    `json:"user_id" format:"id"`
	ProductID  uint    `json:"product_id" format:"id"`
	ItemCost   float32 `json:"item_cost" format:"itemcost"`
	ItemAmount float32 `json:"item_amount" format:"itemamount"`
	ItemType   uint8   `json:"item_type" format:"itemtype"`
	PersonID   uint    `json:"person_id" format:"id"`
}

// `CREATE TABLE IF NOT EXISTS items (
//       item_id INTEGER PRIMARY KEY AUTOINCREMENT,
//       user_id INTEGER NOT NULL,
//       product_id INTEGER NOT NULL,
//       item_cost REAL DEFAULT 0,
//       item_type INTEGER NOT NULL,
//       person_id INTEGER DEFAULT NULL,
//       CHECK (item_type >= 1 AND item_type <= 3),
//       FOREIGN KEY (user_id) REFERENCES `+userStore.TableName+` (user_id) ON DELETE RESTRICT,
//       FOREIGN KEY (product_id) REFERENCES `+productStore.TableName+` (product_id) ON DELETE RESTRICT,
//       FOREIGN KEY (person_id) REFERENCES `+personStore.TableName+` (person_id) ON DELETE RESTRICT
//   );`)
