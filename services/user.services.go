package services

import (
	"time"

	"github.com/bmg-c/product-diary/db"
)

type UserPublic struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

func NewUserService(up UserPublic, uStore *db.Store) *UserService {
	return &UserService{
		UserPublic: up,
		UserStore:  uStore,
	}
}

type UserService struct {
	UserPublic UserPublic
	UserStore  *db.Store
}

func (us *UserService) GetUserByID(id int) (UserPublic, error) {
	query := `SELECT id, username, email, created_at FROM users
		WHERE id = ?`

	stmt, err := us.UserStore.DB.Prepare(query)
	if err != nil {
		return UserPublic{}, err
	}

	defer stmt.Close()

	us.UserPublic.ID = id
	err = stmt.QueryRow(
		us.UserPublic.ID,
	).Scan(
		&us.UserPublic.ID,
		&us.UserPublic.Username,
		&us.UserPublic.Email,
		&us.UserPublic.CreatedAt,
	)
	if err != nil {
		return UserPublic{}, err
	}

	return us.UserPublic, nil
}

func (us *UserService) GetUsersAll() ([]UserPublic, error) {
	query := `SELECT id, username, email, created_at FROM users ORDER BY created_at DESC`

	rows, err := us.UserStore.DB.Query(query)
	if err != nil {
		return []UserPublic{}, err
	}
	// We close the resource
	defer rows.Close()

	users := []UserPublic{}
	for rows.Next() {
		rows.Scan(
			&us.UserPublic.ID,
			&us.UserPublic.Username,
			&us.UserPublic.Email,
			&us.UserPublic.CreatedAt,
		)

		users = append(users, us.UserPublic)
	}

	return users, nil
}
