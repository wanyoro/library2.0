package utils

import (
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// func EncodeAuthToken createsauthenticaton token
func EncodeAuthToken(uid uint) (string, error) {
	claims := jwt.MapClaims{}
	claims["studentID"] = uid
	claims["IssuedAt"] = time.Now().Unix()
	claims["ExpiresAt"] = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
	return token.SignedString([]byte(os.Getenv("SECRET")))
}
