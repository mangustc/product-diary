package services

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/mail"
	"time"

	"github.com/bmg-c/product-diary/db"
	"github.com/bmg-c/product-diary/errorhandler"
)

type UserPublic struct {
	ID        int       `json:"user_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type UserSignin struct {
	Email string `json:"email"`
}

type UserConfirmSignin struct {
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

func (us *UserService) SigninUser(ur UserSignin) error {
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
		return nil
	}

	err = us.sendCode(ur.Email)
	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) ConfirmSignin(ucr UserConfirmSignin) error {
	err := us.deleteExpiredCodes()
	if err != nil {
		return err
	}

	var code string
	query := `SELECT code FROM ` + us.CodeStore.TableName + ` 
		WHERE email = ?`

	stmt, err := us.UserStore.DB.Prepare(query)
	if err != nil {
		return errorhandler.StatusError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	us.UserPublic.Email = ucr.Email
	err = stmt.QueryRow(
		us.UserPublic.Email,
	).Scan(
		&code,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return errorhandler.StatusError{
				Err:  err,
				Code: http.StatusNotFound,
			}
		}
		return errorhandler.StatusError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	if code != ucr.Code {
		return errorhandler.StatusError{
			Err:  fmt.Errorf("Confirmation codes do not match"),
			Code: http.StatusUnprocessableEntity,
		}
	}

	stmt.Close()

	ul, err := us.addUserToDB(ucr.Email)
	if err != nil {
		return err
	}
	err = us.sendUserLogin(ul)
	return err
}

func (us *UserService) GetUserByID(id int) (UserPublic, error) {
	query := `SELECT user_id, username, email, created_at FROM ` + us.UserStore.TableName + ` 
		WHERE user_id = ?`

	stmt, err := us.UserStore.DB.Prepare(query)
	if err != nil {
		return UserPublic{}, errorhandler.StatusError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
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
		if err == sql.ErrNoRows {
			return UserPublic{}, errorhandler.StatusError{
				Err:  err,
				Code: http.StatusNotFound,
			}
		}
		return UserPublic{}, errorhandler.StatusError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	return us.UserPublic, nil
}

func (us *UserService) GetUsersAll() ([]UserPublic, error) {
	query := `SELECT user_id, username, email, created_at FROM ` + us.UserStore.TableName + ` ORDER BY created_at DESC`

	rows, err := us.UserStore.DB.Query(query)
	if err != nil {
		return []UserPublic{}, errorhandler.StatusError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	defer rows.Close()

	users := []UserPublic{}
	for rows.Next() {
		err = rows.Scan(
			&us.UserPublic.ID,
			&us.UserPublic.Username,
			&us.UserPublic.Email,
			&us.UserPublic.CreatedAt,
		)
		if err != nil {
			return []UserPublic{}, errorhandler.StatusError{
				Err:  err,
				Code: http.StatusInternalServerError,
			}
		}

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
        VALUES (NULL, ?, ?, ?, datetime('now'))`

	stmt, err := us.CodeStore.DB.Prepare(query)
	defer stmt.Close()
	if err != nil {
		return UserLogin{}, errorhandler.StatusError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	_, err = stmt.Exec("master", email, password)
	if err != nil {
		return UserLogin{}, errorhandler.StatusError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	return UserLogin{
		Email:    email,
		Password: password,
	}, nil
}

func (us *UserService) isCodeSent(email string) (bool, error) {
	var id int
	query := `SELECT code_id FROM ` + us.CodeStore.TableName + ` WHERE email = ?`

	err := us.CodeStore.DB.QueryRow(query, email).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return true, errorhandler.StatusError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	return true, nil
}

func (us *UserService) deleteExpiredCodes() error {
	query := `DELETE FROM ` + us.CodeStore.TableName + ` 
        WHERE created_at <= datetime('now', '-5 minutes')`

	stmt, err := us.CodeStore.DB.Prepare(query)
	defer stmt.Close()
	if err != nil {
		return errorhandler.StatusError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	_, err = stmt.Exec()
	if err != nil {
		return errorhandler.StatusError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	return nil
}

func (us *UserService) sendCode(email string) error {
	query := `INSERT INTO ` + us.CodeStore.TableName + `(code_id, email, code, created_at)
        VALUES (NULL, ?, ?, datetime('now'))`

	stmt, err := us.CodeStore.DB.Prepare(query)
	defer stmt.Close()
	if err != nil {
		return errorhandler.StatusError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	_, err = stmt.Exec(email, "000000")
	if err != nil {
		return errorhandler.StatusError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	return nil
}
