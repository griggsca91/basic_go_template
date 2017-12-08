package twitchboard

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id             int64
	Username       string `sql:",pk"`
	HashedPassword string
	Password       string `-`
	Email          string `sql:",pk"`
}

func ValidateCredentials(username, password string) bool {
	user := GetUser(username)

	err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	return err == nil
}

func GetUser(username string) *User {

	return nil
}
