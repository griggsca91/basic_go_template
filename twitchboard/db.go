package twitchboard

import "github.com/go-pg/pg"

func DB() *pg.DB {
	db := pg.Connect(&pg.Options{
		User: "postgres",
	})

	return db
}
