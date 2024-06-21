package user_db

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/bmg-c/product-diary/db"
	"github.com/bmg-c/product-diary/errorhandler"
	L "github.com/bmg-c/product-diary/localization"
	"github.com/bmg-c/product-diary/schemas"
	"github.com/bmg-c/product-diary/schemas/user_schemas"
)

type UserDB struct {
	userStore *db.Store
	codeStore *db.Store
}

func NewUserDB(userStore *db.Store, codeStore *db.Store) (*UserDB, error) {
	if userStore == nil || codeStore == nil {
		return nil, fmt.Errorf("Error creating UserDB instance, one of the stores is nil")
	}
	return &UserDB{
		userStore: userStore,
		codeStore: codeStore,
	}, nil
}

func (udb *UserDB) AddCode(email string) error {
	err := udb.deleteExpiredCodes()
	if err != nil {
		return errorhandler.StatusError{
			Err:  L.GetError(L.MsgErrorInternalServer),
			Code: http.StatusInternalServerError,
		}
	}

	query := `INSERT INTO ` + udb.codeStore.TableName + `(code_id, email, code, created_at)
        VALUES (NULL, ?, ?, datetime('now'))`

	stmt, err := udb.codeStore.DB.Prepare(query)
	defer stmt.Close()
	if err != nil {
		return errorhandler.StatusError{
			Err:  L.GetError(L.MsgErrorInternalServer),
			Code: http.StatusInternalServerError,
		}
	}
	_, err = stmt.Exec(email, "000000")
	if err != nil {
		return errorhandler.StatusError{
			Err:  L.GetError(L.MsgErrorInternalServer),
			Code: http.StatusInternalServerError,
		}
	}

	return nil
}

func (udb *UserDB) GetCode(email string) (string, error) {
	err := udb.deleteExpiredCodes()
	if err != nil {
		return "", errorhandler.StatusError{
			Err:  L.GetError(L.MsgErrorInternalServer),
			Code: http.StatusInternalServerError,
		}
	}

	var code string = ""
	query := `SELECT code FROM ` + udb.codeStore.TableName + ` WHERE email = ?`

	err = udb.codeStore.DB.QueryRow(query, email).Scan(&code)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", errorhandler.StatusError{
			Err:  L.GetError(L.MsgErrorInternalServer),
			Code: http.StatusInternalServerError,
		}
	}
	return code, nil
}

func (udb *UserDB) AddUser(email string) error {
	username := "master"
	password := "awooga"

	query := `INSERT INTO ` + udb.userStore.TableName + `(user_id, username, email, password, created_at)
        VALUES (NULL, ?, ?, ?, datetime('now'))`

	stmt, err := udb.userStore.DB.Prepare(query)
	defer stmt.Close()
	if err != nil {
		return errorhandler.StatusError{
			Err:  L.GetError(L.MsgErrorInternalServer),
			Code: http.StatusInternalServerError,
		}
	}
	_, err = stmt.Exec(username, email, password)
	if err != nil {
		return errorhandler.StatusError{
			Err:  L.GetError(L.MsgErrorInternalServer),
			Code: http.StatusInternalServerError,
		}
	}

	return nil
}

func (udb *UserDB) GetUser(userInfo user_schemas.GetUser) (user_schemas.UserDB, error) {
	var up user_schemas.UserDB = user_schemas.UserDB{}

	var query string = ""
	var arg any
	if !schemas.IsZero(userInfo.UserID) {
		query = `SELECT user_id, username, email, password, created_at FROM ` + udb.userStore.TableName + ` 
		    WHERE user_id = ?`
		arg = userInfo.UserID
	} else if !schemas.IsZero(userInfo.Email) {
		query = `SELECT user_id, username, email, password, created_at FROM ` + udb.userStore.TableName + ` 
		    WHERE email = ?`
		arg = userInfo.Email
	} else {
		return user_schemas.UserDB{}, errorhandler.StatusError{
			Err:  L.GetError(L.MsgErrorGetUserNoInfo),
			Code: http.StatusUnprocessableEntity,
		}
	}

	stmt, err := udb.userStore.DB.Prepare(query)
	if err != nil {
		return user_schemas.UserDB{}, errorhandler.StatusError{
			Err:  L.GetError(L.MsgErrorInternalServer),
			Code: http.StatusInternalServerError,
		}
	}
	defer stmt.Close()

	err = stmt.QueryRow(
		arg,
	).Scan(
		&up.UserID,
		&up.Username,
		&up.Email,
		&up.Password,
		&up.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return user_schemas.UserDB{}, errorhandler.StatusError{
				Err:  L.GetError(L.MsgErrorGetUserNotFound),
				Code: http.StatusNotFound,
			}
		}
		return user_schemas.UserDB{}, errorhandler.StatusError{
			Err:  L.GetError(L.MsgErrorInternalServer),
			Code: http.StatusInternalServerError,
		}
	}

	ve := schemas.ValidateStruct(up)
	if ve != nil {
		return user_schemas.UserDB{}, errorhandler.StatusError{
			Err:  L.GetError(L.MsgErrorInternalServer),
			Code: http.StatusInternalServerError,
		}
	}

	return up, nil
}

func (udb *UserDB) GetUsersAll() ([]user_schemas.UserDB, error) {
	var userDB user_schemas.UserDB = user_schemas.UserDB{}
	query := `SELECT user_id, username, email, password, created_at FROM ` + udb.userStore.TableName +
		` ORDER BY created_at DESC`

	rows, err := udb.userStore.DB.Query(query)
	if err != nil {
		return []user_schemas.UserDB{}, errorhandler.StatusError{
			Err:  L.GetError(L.MsgErrorInternalServer),
			Code: http.StatusInternalServerError,
		}
	}
	defer rows.Close()

	users := []user_schemas.UserDB{}
	for rows.Next() {
		err = rows.Scan(
			&userDB.UserID,
			&userDB.Username,
			&userDB.Email,
			&userDB.Password,
			&userDB.CreatedAt,
		)
		if err != nil {
			return []user_schemas.UserDB{}, errorhandler.StatusError{
				Err:  L.GetError(L.MsgErrorInternalServer),
				Code: http.StatusInternalServerError,
			}
		}

		users = append(users, userDB)
	}

	return users, nil
}

func (udb *UserDB) deleteExpiredCodes() error {
	query := `DELETE FROM ` + udb.codeStore.TableName + ` 
        WHERE created_at <= datetime('now', '-5 minutes')`

	stmt, err := udb.codeStore.DB.Prepare(query)
	defer stmt.Close()
	if err != nil {
		return errorhandler.StatusError{
			Err:  L.GetError(L.MsgErrorInternalServer),
			Code: http.StatusInternalServerError,
		}
	}
	_, err = stmt.Exec()
	if err != nil {
		return errorhandler.StatusError{
			Err:  L.GetError(L.MsgErrorInternalServer),
			Code: http.StatusInternalServerError,
		}
	}

	return nil
}
