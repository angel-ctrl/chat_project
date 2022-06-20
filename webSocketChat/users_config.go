package webSocketChat

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sql_chat/handlers"
	"github.com/sql_chat/models"
)

const (
	pongWait = 30 * time.Second

	pingPeriod = (pongWait * 9) / 10

	maxMessageSize = 512
)

type Client struct {
	UserID     string
	Username   string
	Connection *websocket.Conn
}

func NewClient(id string, name string, conn *websocket.Conn) *Client {
	return &Client{
		UserID:     id,
		Username:   name,
		Connection: conn,
	}
}

func (u *Client) warnfriendsAndMe() []models.Users {

	id, _ := strconv.Atoi(u.UserID)

	friends, err := handlers.UserUseCase.LookFriends(id)

	if err != nil {
		log.Println(err)
		return nil
	}

	var userConecteds []models.Users

	for _, element := range friends {

		if userCnn, ok := Ws.clients[element.Username]; ok {

			userCnn.Connection.WriteMessage(websocket.TextMessage, []byte("{{"+"uc:"+u.Username+"}}"))

			var user models.Users
			user.Id = element.Id
			user.Username = element.Username
			userConecteds = append(userConecteds, user)
		}
	}

	return userConecteds

}

func (u *Client) OnLine() {

	u.Connection.SetReadLimit(maxMessageSize)
	u.Connection.SetReadDeadline(time.Now().Add(pongWait))
	u.Connection.SetPongHandler(func(string) error { u.Connection.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	lst := u.warnfriendsAndMe()

	pingTicker := time.NewTicker(pingPeriod)

	if data, err := json.Marshal(lst); err != nil {

		log.Println(err)

		return

	} else {

		u.Connection.WriteMessage(websocket.TextMessage, data)

	}

	go u.aliveConection(pingTicker)

	for {

		if _, message, err := u.Connection.ReadMessage(); err != nil { 
			
			log.Println("Error on read message: ", err.Error())
			pingTicker.Stop()
			u.Connection.Close()
			break

		} else {

			sms := &models.Message{}

			fmt.Println("Data: ", string(message))

			if err := json.Unmarshal(message, sms); err != nil {
				continue
			} else {
				log.Println(sms)
				Ws.chanel.messageChanel <- sms
			}
		}

	}

	Ws.chanel.leaveChanel <- u
}

func (u *Client) aliveConection(pingTicker *time.Ticker) {

	for {

		<-pingTicker.C

		if err := u.Connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
			fmt.Println("ping: ", err)
		}

	}

}

func (u *Client) DisconectUserAndFriends() error {

	id, _ := strconv.Atoi(u.UserID)

	friends, err := handlers.UserUseCase.LookFriends(id)

	if err != nil {
		log.Println(err)
		return err
	}

	for _, element := range friends {

		if userCnn, ok := Ws.clients[element.Username]; ok {
			userCnn.Connection.WriteMessage(websocket.TextMessage, []byte("{{"+"Disconected:"+u.Username+"}}"))
		}
	}

	return nil

}

func (u *Client) SendMessage(message *models.Message) error {

	if data, err := json.Marshal(message); err != nil {
		return err
	} else {
		err = u.Connection.WriteMessage(websocket.TextMessage, data)
		log.Printf("Message send: from %s to %s", message.UserSender, message.UserDestination)
		return err
	}
}
