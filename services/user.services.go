package services

import "time"

type UserPublic struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

func NewUserService(up UserPublic) *UserService {
	return &UserService{
		UserPublic: up,
	}
}

type UserService struct {
	UserPublic UserPublic
}

func (us *UserService) GetUserByID(id int) UserPublic {
	return UserPublic{
		ID:        0,
		Username:  "sasi",
		Email:     "sasamba@gmail.com",
		CreatedAt: time.Now(),
	}
}

func (us *UserService) GetUsersAll() []UserPublic {
	users := []UserPublic{}
	users = append(users, UserPublic{
		ID:        0,
		Username:  "sasi",
		Email:     "sasamba@gmail.com",
		CreatedAt: time.Now(),
	})
	users = append(users, UserPublic{
		ID:        1,
		Username:  "isas",
		Email:     "abmasas@liamg.moc",
		CreatedAt: time.Now(),
	})
	return users
}
