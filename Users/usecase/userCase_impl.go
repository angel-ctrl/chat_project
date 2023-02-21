package usecase

import (
	"errors"

	"github.com/sql_chat/Users"
	"github.com/sql_chat/models"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCaseImpl struct {
	user Users.UserRepo
}

func CreateUserUseCase(userRepo Users.UserRepo) Users.UserUseCase {
	return &UserUseCaseImpl{userRepo}
}

func (e *UserUseCaseImpl) GetUser(id int) (models.Users, error) {
	return e.user.GetUser_repo(id)
}

func (e *UserUseCaseImpl) GetUserByName(name string) (models.Users, error) {
	return e.user.GetUserByName_repo(name)
}

func (e *UserUseCaseImpl) GetAllUser() ([]models.Users, error) {
	return e.user.GetAllUser_repo()
}

func (e *UserUseCaseImpl) CreateUser(user models.Users) (models.Users, error) {
	return e.user.CreateUser_repo(user)
}

func (e *UserUseCaseImpl) UpdateUser(id int, user models.Users) (models.Users, error) {
	return e.user.UpdateUser_repo(id, user)
}

func (e *UserUseCaseImpl) DeleteUser(id int) error {
	return e.user.DeleteUser_repo(id)
}

func (e *UserUseCaseImpl) LoginUserCase(pass string, u models.Users) error {

	passwordBytes := []byte(pass)

	passwordBD := []byte(u.Password)

	err := bcrypt.CompareHashAndPassword(passwordBD, passwordBytes)

	if err != nil {
		return errors.New("contraseña incorrecta")
	}

	return nil
}

func (e *UserUseCaseImpl) AddFriend(u models.Friends) error {
	return e.user.AddFriend_repo(u)
}

func (e *UserUseCaseImpl) LookFriends(id int) ([]models.Users, error) {
	return e.user.LookFriends_repo(id)
}

func (e *UserUseCaseImpl) GetUserName(name string) (models.Users, error) {
	return e.user.GetUserName_repo(name)
}