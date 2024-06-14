package services

import (
	"database/sql"
	"fmt"
	"net/mail"
	"time"

	"github.com/bmg-c/product-diary/db"
)

type UserPublic struct {
	ID        int       `json:"user_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type UserRegister struct {
	Email string `json:"email"`
}

type UserConfirmRegister struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUserService(up UserPublic, userStore *db.Store, codeStore *db.Store) *UserService {
	return &UserService{
		UserPublic: up,
		UserStore:  userStore,
		CodeStore:  codeStore,
	}
}

type UserService struct {
	UserPublic UserPublic
	UserStore  *db.Store
	CodeStore  *db.Store
}

func (us *UserService) RegisterUser(ur UserRegister) error {
	_, err := mail.ParseAddress(ur.Email)
	if err != nil {
		return err
	}

	err = us.deleteExpiredCodes()
	if err != nil {
		return err
	}

	isCodeSent, err := us.isCodeSent(ur.Email)
	if err != nil {
		return err
	}
	if isCodeSent {
		return fmt.Errorf("Code already has been sent")
	}

	err = us.sendCode(ur.Email)
	return err
}

func (us *UserService) ConfirmRegister(ucr UserConfirmRegister) error {
	err := us.deleteExpiredCodes()
	if err != nil {
		return err
	}

	var code string
	query := `SELECT code FROM ` + us.CodeStore.TableName + ` 
		WHERE email = ?`

	stmt, err := us.UserStore.DB.Prepare(query)
	if err != nil {
		return err
	}

	us.UserPublic.Email = ucr.Email
	err = stmt.QueryRow(
		us.UserPublic.Email,
	).Scan(
		&code,
	)
	if err != nil {
		return err
	}
	if code != ucr.Code {
		return fmt.Errorf("Confirmation codes do not match")
	}

	stmt.Close()

	ul, err := us.addUserToDB(ucr.Email)
	err = us.sendUserLogin(ul)
	return err
}

func (us *UserService) GetUserByID(id int) (UserPublic, error) {
	query := `SELECT user_id, username, email, created_at FROM ` + us.UserStore.TableName + ` 
		WHERE user_id = ?`

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
	query := `SELECT user_id, username, email, created_at FROM ` + us.UserStore.TableName + ` ORDER BY created_at DESC`

	rows, err := us.UserStore.DB.Query(query)
	if err != nil {
		return []UserPublic{}, err
	}
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

func (us *UserService) sendUserLogin(ul UserLogin) error {
	return nil
}

func (us *UserService) addUserToDB(email string) (UserLogin, error) {
	password := "awooga"

	query := `INSERT INTO ` + us.UserStore.TableName + `(user_id, username, email, password, created_at)
        VALUES (?, ?, ?, ?, datetime('now'))`

	stmt, err := us.CodeStore.DB.Prepare(query)
	defer stmt.Close()
	if err != nil {
		return UserLogin{}, err
	}
	_, err = stmt.Exec(nil, "master", email, password)

	return UserLogin{
		Email:    email,
		Password: password,
	}, err
}

func (us *UserService) isCodeSent(email string) (bool, error) {
	var id int
	query := `SELECT code_id FROM ` + us.CodeStore.TableName + ` WHERE email = ?`

	err := us.CodeStore.DB.QueryRow(query, email).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return true, err
	}
	return true, nil
}

func (us *UserService) deleteExpiredCodes() error {
	query := `DELETE FROM ` + us.CodeStore.TableName + ` 
        WHERE created_at <= datetime('now', '-5 minutes')`

	stmt, err := us.CodeStore.DB.Prepare(query)
	defer stmt.Close()
	if err != nil {
		return err
	}
	_, err = stmt.Exec()

	return err
}

func (us *UserService) sendCode(email string) error {
	query := `INSERT INTO ` + us.CodeStore.TableName + `(code_id, email, code, created_at)
        VALUES (?, ?, ?, datetime('now'))`

	stmt, err := us.CodeStore.DB.Prepare(query)
	defer stmt.Close()
	if err != nil {
		return err
	}
	_, err = stmt.Exec(nil, email, "000000")

	return err
}
