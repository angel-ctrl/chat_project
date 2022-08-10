package utils

import (
	"errors"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

type Claim struct {
	Username string `json:"user"`
	Id       string `json:"id"`
	jwt.StandardClaims
}

func ProcesoToken(tk string) (*Claim, bool, error) {

	miClave := []byte("hackernoob")

	claims := &Claim{}

	splitToken := strings.Split(tk, "Bearer")

	if len(splitToken) != 2 {
		return claims, false, errors.New("formato de token invalido")
	}

	tk = strings.TrimSpace(splitToken[1])

	tkn, err := jwt.ParseWithClaims(tk, claims, func(token *jwt.Token) (interface{}, error) {
		return miClave, nil
	})

	if !tkn.Valid {
		return claims, false, errors.New("token invalido")
	}

	return claims, false, err
}
