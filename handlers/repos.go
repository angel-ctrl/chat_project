package handlers

import (
	"database/sql"

	"github.com/sql_chat/Users"
	user_repo "github.com/sql_chat/Users/repo"
	user_usecase "github.com/sql_chat/Users/usecase"
	"github.com/sql_chat/database"
)

var DB *sql.DB = database.ConnectDB()

//repositorys
var UserRepo Users.UserRepo = user_repo.CreateUserRepo(DB)

//usecases
var UserUseCase Users.UserUseCase = user_usecase.CreateUserUseCase(UserRepo)