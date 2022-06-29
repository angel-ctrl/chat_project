package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	user_handler "github.com/sql_chat/Users/handler"
	"github.com/sql_chat/handlers"
	"github.com/sql_chat/webSocketChat"
)

func main() {

	/*a := models.Message{
		Id:              1,
		UserSender:      "Juan",
		Msg:             "Hola",
		UserDestination: "gary",
	}

	p, _ := json.Marshal(a)

	client.Set("key", p, 30*time.Second)*/

	//a := []int{1, 2, 3}

	//p, _ := json.Marshal(a)

	/*t := 30*time.Second

	client.Set("key", p, 30*time.Second)*/

	/*val, err := client.Get("key").Result()
	if err != nil {
		fmt.Println("no existe papa")
	}
	fmt.Println("key: ", val)*/

	/*a := models.Message{
		Id:              1,
		UserSender:      "Juan",
		Msg:             "Hola",
		UserDestination: "gary",
	}*/

	/*aa := models.Message{
		Id:              1,
		UserSender:      "Juan",
		Msg:             "Hola",
		UserDestination: "gary",
	}

	ab := models.Message{
		Id:              1,
		UserSender:      "Juan",
		Msg:             "Hola",
		UserDestination: "gary",
	}

	a := []models.Message{aa, ab}

	err := redis.Set("key1", a, 30*time.Second)

	if err != nil {
		log.Fatal(err)
	}

	val, err := redis.Get("key1")

	if err != nil {
		fmt.Println("error", err.Error())
	}

	var ints []models.Message

	err = json.Unmarshal([]byte(val.(string)), &ints)

	fmt.Println("key: ", ints[0].Msg)*/

	/*data := models.Message{}
	json.Unmarshal([]byte(val.(string)), &data)

	fmt.Println("key: ", data)*/

	router := mux.NewRouter()

	//routes
	user_handler.CreateUserHandler(router, handlers.UserUseCase)

	PORT := os.Getenv("PORT")

	if PORT == "" {
		PORT = "8080"
	}

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
