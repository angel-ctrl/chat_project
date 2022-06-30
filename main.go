package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	security "github.com/sql_chat/Security"
	user_handler "github.com/sql_chat/Users/handler"
	"github.com/sql_chat/handlers"
	"github.com/sql_chat/webSocketChat"
)

func main() {

	router := mux.NewRouter()

	//routes
	user_handler.CreateUserHandler(router, handlers.UserUseCase)

	PORT := os.Getenv("PORT")

	if PORT == "" {
		PORT = "8080"
	}

	router.HandleFunc("/API/SendPublicKey", security.SendPublicKey).Methods("GET")

	webSocketChat.NewHub()

	router.HandleFunc("/API/ws", webSocketChat.Ws.HandlerWebSocket).Methods("GET")
	go webSocketChat.Ws.UsersManager()

	cor := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://25.12.136.60:3000", "http://127.0.0.1:8000"},
		AllowedHeaders:   []string{"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token", "Authorization", "Accept", "Accept-Language"},
		AllowedMethods:   []string{"GET", "PATCH", "POST", "PUT", "OPTIONS", "DELETE", "COPY"},
		Debug:            true,
		AllowCredentials: true,
	})

	handler := cor.Handler(router)

	log.Fatal(http.ListenAndServe(":"+PORT, handler))
}
