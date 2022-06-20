package repo

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/sql_chat/Users"
	"github.com/sql_chat/models"
	"github.com/sql_chat/utils"
)

type UserRepoImpl struct {
	DB *sql.DB
}

func CreateUserRepo(DB *sql.DB) Users.UserRepo {
	return &UserRepoImpl{DB}
}

func (e *UserRepoImpl) GetUser_repo(id int) (models.Users, error) {

	var user models.Users

	query := "SELECT * FROM usuarios WHERE id = $1;"

	if err := e.DB.QueryRow(query, id).Scan(&user.Id, &user.Username, &user.Lastname, &user.Password, &user.Active); err != nil {
		return user, errors.New(err.Error())
	}

	return user, nil

}

func (e *UserRepoImpl) GetUserByName_repo(name string) (models.Users, error) {

	var user models.Users

	query := "SELECT * FROM usuarios WHERE username = $1;"

	if err := e.DB.QueryRow(query, name).Scan(&user.Id, &user.Username, &user.Lastname, &user.Password, &user.Active); err != nil {
		return user, errors.New(err.Error())
	}

	return user, nil

}

func (e *UserRepoImpl) GetAllUser_repo() ([]models.Users, error) {

	var list []models.Users

	query := "SELECT * FROM usuarios;"

	m, err := e.DB.Query(query)

	if err != nil {
		return list, errors.New(err.Error())
	}

	defer m.Close()

	for m.Next() {
		var user models.Users
		err := m.Scan(&user.Id, &user.Username, &user.Lastname, &user.Password, &user.Active)
		if err != nil {
			fmt.Println(err)
			break
		}
		list = append(list, user)
	}

	return list, nil

}

func (e *UserRepoImpl) CreateUser_repo(user models.Users) (models.Users, error) {

	query := "INSERT INTO usuarios(username, lastname, password, active) values($1, $2, $3, $4)"

	pass, err := utils.EncriptPass(user.Password)

	if err != nil {
		return models.Users{}, err
	}

	if _, err := e.DB.Exec(query, user.Username, user.Lastname, pass, user.Active); err != nil {
		return models.Users{}, err
	}

	return user, nil

}

func (e *UserRepoImpl) UpdateUser_repo(id int, user models.Users) (models.Users, error) {

	query := "UPDATE usuarios SET username = $1, lastname = $2, password = $3, active = $4 WHERE id = $5;"

	if _, err := e.DB.Exec(query, user.Username, user.Lastname, user.Password, user.Active, id); err != nil {
		return models.Users{}, err
	}

	return user, nil

}

func (e *UserRepoImpl) DeleteUser_repo(id int) error {

	query := "DELETE FROM usuarios WHERE id = $1;"

	if _, err := e.DB.Exec(query, id); err != nil {
		return err
	}

	return nil

}
 
func (e *UserRepoImpl) AddFriend_repo(u models.Friends) error {

	query := "INSERT INTO friends (IDuser1, IDuser2) VALUES ($1, $2) RETURNING id, IDuser1, IDuser2;"

	err := e.DB.QueryRow(query, u.IDuser1, u.IDuser2)

	if err != nil {
		return err.Err()
	}

	return nil

}

func (e *UserRepoImpl) LookFriends_repo(id int) ([]models.Users, error) {

	query := `SELECT usuarios.id, usuarios.username, usuarios.lastname
			  FROM friends 
			  inner join usuarios on friends.iduser2 = usuarios.Id or friends.iduser1 = usuarios.Id
			  WHERE (friends.iduser2 = $1 or friends.iduser1 = $2) and usuarios.id != $3;`


	m, err := e.DB.Query(query, id, id, id)

	var list []models.Users

	if err != nil {
		return list, errors.New(err.Error())
	}

	defer m.Close()

	for m.Next() {

		var userModel models.Users

		err := m.Scan(&userModel.Id, &userModel.Username, &userModel.Lastname)
		if err != nil {
			fmt.Println(err)
			break
		}
		list = append(list, userModel)
	}

	return list, nil

}