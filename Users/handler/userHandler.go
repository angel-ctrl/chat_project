package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	security "github.com/sql_chat/Security"
	"github.com/sql_chat/Users"
	"github.com/sql_chat/middlewares"
	"github.com/sql_chat/models"
	"github.com/sql_chat/utils"
)

type UserHandler struct {
	Allfuncs Users.UserUseCase
}

func CreateUserHandler(router *mux.Router, userUseCase Users.UserUseCase) {

	userHandler := UserHandler{userUseCase}

	router.HandleFunc("/API/users/v1", middlewares.ValidoJWT(userHandler.UsersController)).Methods("POST", "GET", "PUT", "DELETE", "PATCH")
	router.HandleFunc("/API/login/v1", userHandler.Login).Methods("POST")
	router.HandleFunc("/API/addfriend/v1", middlewares.ValidoJWT(userHandler.AddFriendHD)).Methods("POST")
	router.HandleFunc("/API/look_friend/v1", middlewares.ValidoJWT(userHandler.LookFriends)).Methods("GET")

}

func (e *UserHandler) UsersController(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))

	switch r.Method {

	case "GET":

		ID := r.URL.Query().Get("id")
		OPC := r.URL.Query().Get("opc")

		id, err := strconv.Atoi(ID)

		if err != nil {
			http.Error(w, "error convirtiendo string a int", 400)
			return
		}

		switch OPC {

		case "ByID":

			r, err := e.Allfuncs.GetUser(id)

			if err != nil {
				http.Error(w, "error leeyendo: "+err.Error(), 400)
				return
			}

			w.WriteHeader(http.StatusOK)

			json.NewEncoder(w).Encode(r)

		case "ALL":

			r, err := e.Allfuncs.GetAllUser()

			if err != nil {
				http.Error(w, "error leeyendo: "+err.Error(), 400)
				return
			}

			w.WriteHeader(http.StatusOK)

			json.NewEncoder(w).Encode(r)

		default:

			http.Error(w, "Opcion no valida: ", http.StatusUnprocessableEntity)
			return

		}

	case "POST":

		var u models.Users

		err := json.NewDecoder(r.Body).Decode(&u)

		if err != nil {
			http.Error(w, "error en los datos recibidos: "+err.Error(), 400)
			return
		}

		md, err := e.Allfuncs.CreateUser(u)

		if err != nil {
			http.Error(w, "error creando: "+err.Error(), 400)
			return
		}

		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(md)

	case "PUT":

		ID := r.URL.Query().Get("id")

		var u models.Users

		err := json.NewDecoder(r.Body).Decode(&u)

		if err != nil {
			http.Error(w, "error en los datos recibidos: "+err.Error(), 400)
			return
		}

		if len(u.Username) < 1 {
			http.Error(w, "falta el campo nombre", http.StatusUnprocessableEntity)
			return
		}

		if len(u.Lastname) < 1 {
			http.Error(w, "falta el campo apellido", http.StatusUnprocessableEntity)
			return
		}

		if len(u.Password) < 1 {
			http.Error(w, "falta el campo contraseña", http.StatusUnprocessableEntity)
			return
		}

		id, err := strconv.Atoi(ID)

		if err != nil {
			http.Error(w, "error convirtiendo string a int", 400)
			return
		}

		_, err = e.Allfuncs.UpdateUser(id, u)

		if err != nil {
			http.Error(w, "error actualizando:  "+err.Error(), 400)
			return
		}

		w.WriteHeader(http.StatusOK)

	case "DELETE":

		ID := r.URL.Query().Get("id")

		id, err := strconv.Atoi(ID)

		if err != nil {
			http.Error(w, "error convirtiendo string a int", 400)
			return
		}

		e.Allfuncs.DeleteUser(id)

		if err != nil {
			http.Error(w, "error borrando: "+err.Error(), 400)
			return
		}

		w.WriteHeader(http.StatusOK)

	default:

		http.Error(w, "method not supported", http.StatusMethodNotAllowed)
		return

	}
}

func (e *UserHandler) Login(w http.ResponseWriter, r *http.Request) {

	var u models.Users

	err := json.NewDecoder(r.Body).Decode(&u)

	u.Password = security.DescryptMessage(u.Password)

	if err != nil {
		http.Error(w, "error en los datos recibidos: "+err.Error(), 400)
		return
	}

	userDB, err := e.Allfuncs.GetUserByName(u.Username)

	if err != nil {
		http.Error(w, "error: "+err.Error(), 400)
		return
	}

	err = e.Allfuncs.LoginUserCase(u.Password, userDB)

	if err != nil {
		http.Error(w, "error: "+err.Error(), 400)
		return
	}

	jwtKey, err := utils.GeneroJWT(userDB)

	if err != nil {
		http.Error(w, "error: "+err.Error(), 400)
		return
	}

	w.WriteHeader(http.StatusOK)

	type loginr struct {
		Token string `json:"token"`
	}

	var urr loginr

	urr.Token = jwtKey

	json.NewEncoder(w).Encode(urr)
}

func (e *UserHandler) AddFriendHD(w http.ResponseWriter, r *http.Request) {

	var u models.Friends

	err := json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		http.Error(w, "error en los datos recibidos: "+err.Error(), 400)
		return
	}

	err = e.Allfuncs.AddFriend(u)

	if err != nil {
		http.Error(w, "error: "+err.Error(), 400)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func (e *UserHandler) LookFriends(w http.ResponseWriter, r *http.Request) {

	ID := r.URL.Query().Get("id")

	id, err := strconv.Atoi(ID)

	if err != nil {
		http.Error(w, "error convirtiendo string a int", 400)
		return
	}

	lst, err := e.Allfuncs.LookFriends(id)

	if err != nil {
		http.Error(w, "error borrando: "+err.Error(), 400)
		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(lst)

}
