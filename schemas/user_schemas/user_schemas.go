package user_schemas

import (
	"time"

	"github.com/google/uuid"
)

type UserPublic struct {
	UserID    uint      `json:"user_id" format:"id"`
	Username  string    `json:"username" format:"username"`
	Email     string    `json:"email" format:"email"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type UserSignin struct {
	Email string `json:"email" format:"email"`
}

type UserConfirmSignin struct {
	Email string `json:"email" format:"email"`
	Code  string `json:"code" format:"code"`
}

type UserGetByID struct {
	UserID uint `json:"user_id" format:"id"`
}

type UserLogin struct {
	Email    string `json:"email" format:"email"`
	Password string `json:"password" format:"password"`
}

type UserDB struct {
	UserID    uint      `json:"user_id" format:"id"`
	Username  string    `json:"username" format:"username"`
	Email     string    `json:"email" format:"email"`
	Password  string    `json:"password" format:"password"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type GetUser struct {
	UserID uint   `json:"user_id" format:"id" validate:"omitzero"`
	Email  string `json:"email" format:"email" validate:"omitzero"`
}

type SessionDB struct {
	SessionUUID uuid.UUID `json:"session_uuid"`
	UserID      uint      `json:"user_id" format:"id"`
}

type GetSession struct {
	SessionUUID uuid.UUID `json:"session_uuid"`
}

type PersonDB struct {
	PersonID   uint   `json:"person_id" format:"id"`
	UserID     uint   `json:"user_id" format:"id"`
	PersonName string `json:"person_name" format:"username"`
	IsHidden   bool   `json:"is_hidden"`
}

type GetPerson struct {
	UserID     uint   `json:"user_id" format:"id"`
	PersonName string `json:"person_name" format:"username"`
}
