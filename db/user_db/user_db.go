package user_db

import (
	"database/sql"
	"fmt"

	"github.com/bmg-c/product-diary/db"
	E "github.com/bmg-c/product-diary/errorhandler"
	"github.com/bmg-c/product-diary/logger"
	"github.com/bmg-c/product-diary/schemas"
	"github.com/bmg-c/product-diary/schemas/user_schemas"
	"github.com/google/uuid"
)

type UserDB struct {
	userStore    *db.Store
	codeStore    *db.Store
	sessionStore *db.Store
}

func NewUserDB(userStore *db.Store, codeStore *db.Store, sessionStore *db.Store) (*UserDB, error) {
	if userStore == nil || codeStore == nil || sessionStore == nil {
		return nil, fmt.Errorf("Error creating UserDB instance, one of the stores is nil")
	}
	return &UserDB{
		userStore:    userStore,
		codeStore:    codeStore,
		sessionStore: sessionStore,
	}, nil
}

func (udb *UserDB) AddCode(email string) error {
	err := udb.deleteExpiredCodes()
	if err != nil {
		return E.ErrInternalServer
	}

	query := `INSERT INTO ` + udb.codeStore.TableName + `(code_id, email, code, created_at)
        VALUES (NULL, ?, ?, datetime('now'))`

	stmt, err := udb.codeStore.DB.Prepare(query)
	defer stmt.Close()
	if err != nil {
		return E.ErrInternalServer
	}
	_, err = stmt.Exec(email, "000000")
	if err != nil {
		return E.ErrInternalServer
	}

	return nil
}

func (udb *UserDB) GetCode(email string) (string, error) {
	err := udb.deleteExpiredCodes()
	if err != nil {
		return "", E.ErrInternalServer
	}

	var code string = ""
	query := `SELECT code FROM ` + udb.codeStore.TableName + ` WHERE email = ?`

	err = udb.codeStore.DB.QueryRow(query, email).Scan(&code)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", E.ErrNotFound
		}
		return "", E.ErrInternalServer
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
		return E.ErrInternalServer
	}
	_, err = stmt.Exec(username, email, password)
	if err != nil {
		return E.ErrInternalServer
	}

	return nil
}

func (udb *UserDB) GetUser(userInfo user_schemas.GetUser) (user_schemas.UserDB, error) {
	var userDB user_schemas.UserDB = user_schemas.UserDB{}

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
		return user_schemas.UserDB{}, E.ErrUnprocessableEntity
	}

	stmt, err := udb.userStore.DB.Prepare(query)
	if err != nil {
		return user_schemas.UserDB{}, E.ErrInternalServer
	}
	defer stmt.Close()

	err = stmt.QueryRow(
		arg,
	).Scan(
		&userDB.UserID,
		&userDB.Username,
		&userDB.Email,
		&userDB.Password,
		&userDB.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return user_schemas.UserDB{}, E.ErrNotFound
		}
		return user_schemas.UserDB{}, E.ErrInternalServer
	}

	ve := schemas.ValidateStruct(userDB)
	if ve != nil {
		logger.Error.Printf("Invalid user in database %#v. Errors: %s", userDB, ve.Error())
		return user_schemas.UserDB{}, E.ErrInternalServer
	}

	return userDB, nil
}

func (udb *UserDB) GetUsersAll() ([]user_schemas.UserDB, error) {
	var userDB user_schemas.UserDB = user_schemas.UserDB{}
	query := `SELECT user_id, username, email, password, created_at FROM ` + udb.userStore.TableName +
		` ORDER BY created_at DESC`

	rows, err := udb.userStore.DB.Query(query)
	if err != nil {
		return []user_schemas.UserDB{}, E.ErrInternalServer
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
			return []user_schemas.UserDB{}, E.ErrInternalServer
		}

		ve := schemas.ValidateStruct(userDB)
		if ve == nil {
			users = append(users, userDB)
		} else {
			logger.Error.Printf("Invalid user in database %#v. Errors: %s", userDB, ve.Error())
		}
	}

	return users, nil
}

func (udb *UserDB) GetSession(sessionUUID uuid.UUID) (user_schemas.SessionDB, error) {
	var sessionDB user_schemas.SessionDB = user_schemas.SessionDB{}

	var query string = ""
	var arg any
	if !schemas.IsZero(sessionUUID) {
		query = `SELECT session_uuid, user_id FROM ` + udb.sessionStore.TableName + ` 
		    WHERE session_uuid = ?`
		arg = sessionUUID
	} else {
		return user_schemas.SessionDB{}, E.ErrUnprocessableEntity
	}

	stmt, err := udb.sessionStore.DB.Prepare(query)
	if err != nil {
		return user_schemas.SessionDB{}, E.ErrInternalServer
	}
	defer stmt.Close()

	err = stmt.QueryRow(
		arg,
	).Scan(
		&sessionDB.SessionUUID,
		&sessionDB.UserID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return user_schemas.SessionDB{}, E.ErrNotFound
		}
		return user_schemas.SessionDB{}, E.ErrInternalServer
	}

	ve := schemas.ValidateStruct(sessionDB)
	if ve != nil {
		logger.Error.Printf("Invalid session in database %#v. Errors: %s", sessionDB, ve.Error())
		return user_schemas.SessionDB{}, E.ErrInternalServer
	}

	return sessionDB, nil
}

func (udb *UserDB) AddSession(userID uint) (uuid.UUID, error) {
	query := `INSERT INTO ` + udb.sessionStore.TableName + `(session_uuid, user_id)
        VALUES (?, ?)`
	sessionUUID := uuid.New()
	sessionUUIDStr := sessionUUID.String()
	if schemas.IsZero(sessionUUID) {
		return uuid.UUID{}, E.ErrInternalServer
	}

	stmt, err := udb.sessionStore.DB.Prepare(query)
	defer stmt.Close()
	if err != nil {
		return uuid.UUID{}, E.ErrInternalServer
	}
	_, err = stmt.Exec(sessionUUIDStr, userID)
	if err != nil {
		return uuid.UUID{}, E.ErrInternalServer
	}

	return sessionUUID, nil
}

func (udb *UserDB) deleteExpiredCodes() error {
	query := `DELETE FROM ` + udb.codeStore.TableName + ` 
        WHERE created_at <= datetime('now', '-5 minutes')`

	stmt, err := udb.codeStore.DB.Prepare(query)
	defer stmt.Close()
	if err != nil {
		return E.ErrInternalServer
	}
	_, err = stmt.Exec()
	if err != nil {
		return E.ErrInternalServer
	}

	return nil
}
