package Users

import (
	"github.com/sql_chat/models"
)

type UserUseCase interface {
	GetUser(id int) (models.Users, error)
	GetUserByName(name string) (models.Users, error)
	GetAllUser() ([]models.Users, error)
	CreateUser(user models.Users) (models.Users, error)
	UpdateUser(id int, user models.Users) (models.Users, error)
	DeleteUser(id int) error
	LoginUserCase(pass string, user models.Users) error
	AddFriend(models.Friends) error
	LookFriends(id int) ([]models.Users, error)
}