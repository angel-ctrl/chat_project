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
	"github.com/sql_chat/redis"
)

const (
	pongWait = 5 * time.Second

	pingPeriod = (pongWait * 9) / 10

	maxMessageSize = 512
)

type Client struct {
	UserID     string
	Username   string
	Connection *websocket.Conn
	PublicKey  string
}

func NewClient(id string, name string, publicKey string, conn *websocket.Conn) *Client {
	return &Client{
		UserID:     id,
		Username:   name,
		Connection: conn,
		PublicKey:  publicKey,
	}
}

func (u *Client) warnfriendsAndMe() []models.Users {

	val, err := redis.Get(u.Username)

	if err != nil {

		fmt.Println("error ", err.Error())

		id, _ := strconv.Atoi(u.UserID)

		friends, err := handlers.UserUseCase.LookFriends(id)

		if err != nil {
			log.Println(err)
			return nil
		}

		userConecteds := u.verifyUserConecteds(friends)

		err = redis.Set(u.Username, friends, 30*time.Second)

		if err != nil {

			log.Println(err)

		}

		return userConecteds

	} else {

		fmt.Println("si estaba en redis")

		var friends []models.Users

		err = json.Unmarshal([]byte(val.(string)), &friends)

		if err != nil {
			log.Println(err)
		}

		return u.verifyUserConecteds(friends)

	}
}

func (u *Client) verifyUserConecteds(friends []models.Users) []models.Users {

	var userConecteds []models.Users

	var userc models.UserConected

	userc.Username = u.Username
	userc.PublicKey = u.PublicKey

	fmt.Println(u.Username)
	fmt.Println(u.PublicKey)

	for _, element := range friends {

		if userCnn, ok := Ws.clients[element.Username]; ok {

			if data, err := json.Marshal(userc); err == nil {

				userCnn.Connection.WriteMessage(websocket.BinaryMessage, data)

			}

			var user models.Users
			user.Id = element.Id
			user.Username = element.Username
			user.PublicKey = userCnn.PublicKey
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

		u.Connection.WriteMessage(websocket.BinaryMessage, data)

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

			if err := json.Unmarshal(message, sms); err != nil {
				continue
			} else {
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

	val, err := redis.Get(u.Username)

	if err != nil {

		id, _ := strconv.Atoi(u.UserID)

		friends, err := handlers.UserUseCase.LookFriends(id)

		if err != nil {
			log.Println(err)
			return err
		}

		u.DisconectUserOfFriends(friends)

		return nil

	} else {

		fmt.Println("si estaba en redis")

		var friends []models.Users

		err = json.Unmarshal([]byte(val.(string)), &friends)

		if err != nil {
			log.Println(err)
		}

		u.DisconectUserOfFriends(friends)

		return nil

	}

}

func (u *Client) DisconectUserOfFriends(friends []models.Users) {

	for _, element := range friends {

		if userCnn, ok := Ws.clients[element.Username]; ok {
			userCnn.Connection.WriteMessage(websocket.BinaryMessage, []byte("{\""+"Disconected\":\""+u.Username+"\"}"))
		}
	}

}

func (u *Client) SendMessage(message *models.Message) error {

	if data, err := json.Marshal(message); err != nil {
		return err
	} else {
		err = u.Connection.WriteMessage(websocket.BinaryMessage, data)
		return err
	}
}
