package services

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/mail"

	"github.com/bmg-c/product-diary/db"
	"github.com/bmg-c/product-diary/errorhandler"
	"github.com/bmg-c/product-diary/schemas/user_schemas"
)

func NewUserService(userStore *db.Store, codeStore *db.Store) *UserService {
	return &UserService{
		UserStore: userStore,
		CodeStore: codeStore,
	}
}

type UserService struct {
	UserStore *db.Store
	CodeStore *db.Store
}

func (us *UserService) SigninUser(ur user_schemas.UserSignin) error {
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

func (us *UserService) ConfirmSignin(ucr user_schemas.UserConfirmSignin) error {
	err := us.deleteExpiredCodes()
	if err != nil {
		return err
	}

	var email string = ""
	var code string = ""
	query := `SELECT code FROM ` + us.CodeStore.TableName + ` 
		WHERE email = ?`

	stmt, err := us.UserStore.DB.Prepare(query)
	if err != nil {
		return errorhandler.StatusError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	email = ucr.Email
	err = stmt.QueryRow(
		email,
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

func (us *UserService) GetUserByID(id uint) (user_schemas.UserPublic, error) {
	var up user_schemas.UserPublic = user_schemas.UserPublic{}
	query := `SELECT user_id, username, email, created_at FROM ` + us.UserStore.TableName + ` 
		WHERE user_id = ?`

	stmt, err := us.UserStore.DB.Prepare(query)
	if err != nil {
		return user_schemas.UserPublic{}, errorhandler.StatusError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	defer stmt.Close()

	up.UserID = id
	err = stmt.QueryRow(
		up.UserID,
	).Scan(
		&up.UserID,
		&up.Username,
		&up.Email,
		&up.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return user_schemas.UserPublic{}, errorhandler.StatusError{
				Err:  err,
				Code: http.StatusNotFound,
			}
		}
		return user_schemas.UserPublic{}, errorhandler.StatusError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	return up, nil
}

func (us *UserService) GetUsersAll() ([]user_schemas.UserPublic, error) {
	var up user_schemas.UserPublic = user_schemas.UserPublic{}
	query := `SELECT user_id, username, email, created_at FROM ` + us.UserStore.TableName + ` ORDER BY created_at DESC`

	rows, err := us.UserStore.DB.Query(query)
	if err != nil {
		return []user_schemas.UserPublic{}, errorhandler.StatusError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	defer rows.Close()

	users := []user_schemas.UserPublic{}
	for rows.Next() {
		err = rows.Scan(
			&up.UserID,
			&up.Username,
			&up.Email,
			&up.CreatedAt,
		)
		if err != nil {
			return []user_schemas.UserPublic{}, errorhandler.StatusError{
				Err:  err,
				Code: http.StatusInternalServerError,
			}
		}

		users = append(users, up)
	}

	return users, nil
}

func (us *UserService) sendUserLogin(ul user_schemas.UserLogin) error {
	return nil
}

func (us *UserService) addUserToDB(email string) (user_schemas.UserLogin, error) {
	password := "awooga"

	query := `INSERT INTO ` + us.UserStore.TableName + `(user_id, username, email, password, created_at)
        VALUES (NULL, ?, ?, ?, datetime('now'))`

	stmt, err := us.CodeStore.DB.Prepare(query)
	defer stmt.Close()
	if err != nil {
		return user_schemas.UserLogin{}, errorhandler.StatusError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	_, err = stmt.Exec("master", email, password)
	if err != nil {
		return user_schemas.UserLogin{}, errorhandler.StatusError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	return user_schemas.UserLogin{
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
