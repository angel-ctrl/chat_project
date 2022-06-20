package utils

import (
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
)

func ExtractClaims(tokenStr string) (jwt.MapClaims, bool) {
	hmacSecret := []byte("hackernoob")
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		 // check token signing method etc
		 return hmacSecret, nil
	})

	if err != nil {
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		fmt.Printf("Invalid JWT Token")
		return nil, false
	}
}