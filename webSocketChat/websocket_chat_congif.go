package webSocketChat

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sql_chat/models"
	"github.com/sql_chat/utils"
)

var Ws = NewHub()

type MessageChanel chan *models.Message
type UserChanel chan *Client
type Join chan *Client

type Chanel struct {
	messageChanel MessageChanel
	leaveChanel   UserChanel
	join          Join
}

type Hub struct {
	clients map[string]*Client

	chanel *Chanel
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]*Client),
		chanel: &Chanel{
			messageChanel: make(MessageChanel),
			leaveChanel:   make(UserChanel),
			join:          make(chan *Client),
		},
	}
}

func check(r *http.Request) bool {
	log.Printf("%s %s%s %v", r.Method, r.Host, r.RequestURI, r.Proto)
	return r.Method == http.MethodGet
}

var upGradeWebSocket = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     check,
}

func (w *Hub) HandlerWebSocket(rw http.ResponseWriter, r *http.Request) {
	connection, err := upGradeWebSocket.Upgrade(rw, r, nil)

	if err != nil {
		log.Println("No se abrió la connextion")
		fmt.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "Connection failed.")
		return
	}

	_, message, _ := connection.NextReader()

	bytes, _ := ioutil.ReadAll(message)
	
	s := string(bytes)

	dataPublicKey := models.PublicKeyClient{}

	json.Unmarshal([]byte(s), &dataPublicKey)

	cookie, err := r.Cookie("token")

	if err != nil {
		fmt.Println(err)
	}

	token, _ := utils.ExtractClaims(cookie.Value)

	str := fmt.Sprint(token["_id"])
	strName := fmt.Sprint(token["name"])

	fmt.Println(strName, ":  ", dataPublicKey.PublicKey)

	u := NewClient(str, strName, dataPublicKey.PublicKey, connection)
	w.chanel.join <- u

	u.OnLine()
}

func (w *Hub) UsersManager() {
	for {

		select {

		case userChat := <-w.chanel.join:

			if _, ok := w.clients[userChat.Username]; !ok {
				w.clients[userChat.Username] = userChat
			}

			fmt.Println(w.clients)

		case message := <-w.chanel.messageChanel:

			w.ProcessMessage(message)

		case user := <-w.chanel.leaveChanel:

			w.DisconnectUser(user.Username)

		}
	}
}

func (w *Hub) ProcessMessage(message *models.Message) {

	if user, ok := w.clients[message.UserDestination]; ok {
		if err := user.SendMessage(message); err != nil {
			log.Printf("No se pudo mandar el mensaje a %s", message.UserSender)
		}
	}
}

func (w *Hub) DisconnectUser(username string) {
	if user, ok := w.clients[username]; ok {

		user.DisconectUserAndFriends()

		defer user.Connection.Close()
		delete(w.clients, username)
		fmt.Println(w.clients)

	}
}
