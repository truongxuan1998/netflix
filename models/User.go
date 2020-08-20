package models

import (
	"html"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type User struct {
	Email    string `db:"email"`
	Password string `db:"password"`
}
type Movie struct {
	ID    int    `json:"id"		db:"id"`
	Name  string `json:"name" 	db:"name"`
	Genre string `json:"genre"	db:"genre"`
}
type ResponseMovie struct {
	Movies []Movie `json:"movies"`
}
type List struct {
	Name  string `db:"name" 	json:"name"`
	ID    int    `db:"id" 		json:"id"`
	Genre string `db: "genre"	json:"genre"`
}
type ResponseList struct {
	Lists []List `json:"movies"`
}
type Person struct {
	Name string `json:"name"	db:"name"`
	ID   int    `json:"id"		db:"id"`
}
type ResponsePerson struct {
	Persons []Person `json:"persons"`
}

func IsEmpty(data string) bool {
	if len(data) == 0 {
		return true
	} else {
		return false
	}
}
func Santize(data string) string {
	data = html.EscapeString(strings.TrimSpace(data))
	return data
}
func CreateJWT(email string) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["email"] = email
	claims["exp"] = time.Now().Add(24 * time.Hour).Unix()
	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := rawToken.SignedString([]byte("SecretJWT"))
	if err != nil {
		return "", err
	}
	return token, nil
}
