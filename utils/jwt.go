package utils

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go" //el jwt es un alias
	"github.com/sql_chat/models"
)

func GeneroJWT(t models.Users) (string, error){

	miClave := []byte("hackernoob")

	payload := jwt.MapClaims{
		"name": t.Username,
		"_id": t.Id,
		"exp": time.Now().Add(time.Minute * 60).Unix(), //unix hace que devuelva la vaina con formato long y es muy rapido
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenStr, err := token.SignedString(miClave)

	if err != nil {
		return tokenStr, err
	}

	return tokenStr, nil

}