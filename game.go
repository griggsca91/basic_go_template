package main

import (
	"time"
)

type Platform struct {
	Id   int
	Name string
}

type Game struct {
	Id          int
	Name        string
	ReleaseDate time.Time
	ImageURL    string
}

type TrackedGame struct {
	Id             int
	GameId         int
	Game           *Game
	UserId         int
	User           *User
	PlatformId     int
	Platform       *Platform
	Status         string `sql:"default:'NotStarted'"`
	CompletionDate time.Time
	DateAdded      time.Time `sql:"default:Now()"`
}

func GetGame(username string) *Game {
	return nil
}
