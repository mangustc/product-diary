package services

import (
	"errors"

	E "github.com/bmg-c/product-diary/errorhandler"
	"github.com/bmg-c/product-diary/schemas/user_schemas"
	"github.com/google/uuid"
)

func NewUserService(userDB UserDB) *UserService {
	return &UserService{
		userDB: userDB,
	}
}

type UserService struct {
	userDB UserDB
}

type UserDB interface {
	AddCode(email string) error
	GetCode(email string) (string, error)
	AddUser(email string) error
	GetUser(userInfo user_schemas.GetUser) (user_schemas.UserDB, error)
	GetUsersAll() ([]user_schemas.UserDB, error)
	GetSession(sessionUUID uuid.UUID) (user_schemas.SessionDB, error)
	AddSession(userID uint) (uuid.UUID, error)
}

func (us *UserService) SigninUser(ur user_schemas.UserSignin) error {
	userInfo := user_schemas.GetUser{
		Email: ur.Email,
	}
	_, err := us.userDB.GetUser(userInfo)
	if err == nil {
		return E.ErrUnprocessableEntity
	}
	if !errors.Is(err, E.ErrNotFound) {
		return err
	}

	_, err = us.userDB.GetCode(ur.Email)
	if err != nil {
		if !errors.Is(err, E.ErrNotFound) {
			return err
		}
	} else {
		return nil
	}

	err = us.userDB.AddCode(ur.Email)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) ConfirmSignin(ucr user_schemas.UserConfirmSignin) error {
	code, err := us.userDB.GetCode(ucr.Email)
	if err != nil {
		return err
	}
	if code != ucr.Code {
		return E.ErrUnprocessableEntity
	}

	err = us.userDB.AddUser(ucr.Email)
	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) GetUser(userInfo user_schemas.GetUser) (user_schemas.UserPublic, error) {
	udb, err := us.userDB.GetUser(userInfo)
	if err != nil {
		return user_schemas.UserPublic{}, err
	}

	return user_schemas.UserPublic{
		UserID:    udb.UserID,
		Email:     udb.Email,
		Username:  udb.Username,
		CreatedAt: udb.CreatedAt,
	}, nil
}

func (us *UserService) GetUsersAll() ([]user_schemas.UserPublic, error) {
	users, err := us.userDB.GetUsersAll()
	if err != nil {
		return []user_schemas.UserPublic{}, err
	}

	usersPublic := []user_schemas.UserPublic{}
	for _, user := range users {
		usersPublic = append(usersPublic, user_schemas.UserPublic{
			UserID:    user.UserID,
			Email:     user.Email,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
		})
	}

	return usersPublic, nil
}

func (us *UserService) LoginUser(ul user_schemas.UserLogin) (uuid.UUID, error) {
	userInfo := user_schemas.GetUser{
		Email: ul.Email,
	}
	userDB, err := us.userDB.GetUser(userInfo)
	if err != nil {
		return uuid.UUID{}, err
	}
	if ul.Password != userDB.Password {
		return uuid.UUID{}, E.ErrUnprocessableEntity
	}

	sessionUUID, err := us.userDB.AddSession(userDB.UserID)
	if err != nil {
		return uuid.UUID{}, err
	}
	return sessionUUID, nil
}

func (us *UserService) GetUserBySession(sessionUUID uuid.UUID) (user_schemas.UserDB, error) {
	sessionDB, err := us.userDB.GetSession(sessionUUID)
	if err != nil {
		if errors.Is(err, E.ErrNotFound) {
			err = E.ErrUnprocessableEntity
		}
		return user_schemas.UserDB{}, err
	}

	userInfo := user_schemas.GetUser{
		UserID: sessionDB.UserID,
	}
	userDB, err := us.userDB.GetUser(userInfo)
	if err != nil {
		return user_schemas.UserDB{}, err
	}

	return userDB, nil
}

func (us *UserService) sendUserLogin(ul user_schemas.UserLogin) error {
	return nil
}
