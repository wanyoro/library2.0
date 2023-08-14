package utils

import (
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// func EncodeAuthToken createsauthenticaton token
func EncodeAuthTokenStudent(uid uint) (string, error) {
	claims := jwt.MapClaims{}
	claims["studentID"] = uid
	claims["IssuedAt"] = time.Now().Unix()
	claims["ExpiresAt"] = time.Now().Add(time.Hour * 2).Unix()
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
	return token.SignedString([]byte(os.Getenv("SECRET")))
}

func EncodeAuthTokenTeacher(uid uint) (string, error) {
	claims := jwt.MapClaims{}
	claims["teacherID"] = uid
	claims["IssuedAt"] = time.Now().Unix()
	claims["ExpiresAt"] = time.Now().Add(time.Hour * 2).Unix()
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
	return token.SignedString([]byte(os.Getenv("SECRET")))
}
