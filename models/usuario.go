package models

type Users struct {
	Id       int    `json:"id"`
	Username string `json:"user"`
	Lastname string `json:"lastname"`
	Password string `json:"password"`
	Active   bool   `json:"active"`
	PublicKey string `json:"publicKey"`
}