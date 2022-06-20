package models

type Message struct {
	Id              int    `json:"id"`
	UserSender      string `json:"userSender"`
	Msg             string `json:"msg"`
	UserDestination string `json:"userDestination"`
}
