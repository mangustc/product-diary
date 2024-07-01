package tests

import (
	"fmt"

	"github.com/bmg-c/product-diary/schemas"
)

func TestValidation() error {
	type User struct {
		UserID   uint   `json:"user_id" format:"id"`
		Email    string `json:"email" format:"email"`
		Username string `json:"username" format:"username"`
		Password string `json:"password" format:"password"`
		Code     string `json:"code" format:"code"`
		Useless  string `json:"useless" validate:"omitzero,min_length=4"`
	}
	var u User
	var ve schemas.ValidationErrors

	u = User{
		UserID:   0,        // should start at 1
		Email:    "let",    // not an email format
		Username: "A",      // length less than 3
		Password: "pass",   // length less than 6
		Code:     "13U1Ua", // should not contain lowercase latin symbols
		Useless:  "bad",    // omit if zero but still less than 4
	}
	ve = schemas.ValidateStruct(u)
	// fmt.Println("Errors: " + ve.Error())
	if ve == nil {
		return fmt.Errorf("Struct %#v should not be valid", u)
	}

	u = User{
		UserID:   1,
		Email:    "test@gmail.com",
		Username: "Molodez",
		Password: "awooga",
		Code:     "UEOA1O",
	}
	ve = schemas.ValidateStruct(u)
	if ve != nil {
		// fmt.Println("Errors: " + ve.Error())
		return fmt.Errorf("Struct %#v should be valid", u)
	}

	return nil
}
