package Users

import (
	"github.com/sql_chat/models"
)

type UserRepo interface {
	GetUser_repo(id int) (models.Users, error)
	GetUserByName_repo(name string) (models.Users, error)
	GetAllUser_repo() ([]models.Users, error)
	CreateUser_repo(user models.Users) (models.Users, error)
	UpdateUser_repo(id int, user models.Users) (models.Users, error)
	DeleteUser_repo(id int) error
	AddFriend_repo(models.Friends) error
	LookFriends_repo(id int) ([]models.Users, error)
}