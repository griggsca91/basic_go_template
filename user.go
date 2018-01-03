package main

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

type User struct {
	Id             int
	Username       string `sql:",pk"`
	HashedPassword string
	Password       string `-`
	Email          string `sql:",pk"`
	Games          []*TrackedGame
}

func ValidateCredentials(username, password string) bool {
	user := GetUser(username)

	err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	return err == nil
}

func GetUser(username string) *User {
	return nil
}

func (u *User) TrackGame(g *Game) {
	trackedGame := &TrackedGame{
		GameId: g.Id,
		UserId: u.Id,
	}
	db := DB()

	err := db.Insert(trackedGame)
	if err != nil {
		panic(err)
	}
}
