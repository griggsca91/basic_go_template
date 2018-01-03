package main

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

func CreateSchema(db *pg.DB) error {
	for _, model := range []interface{}{
		&User{},
		&Game{},
		&TrackedGame{},
	} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func DB() *pg.DB {
	db := pg.Connect(&pg.Options{
		User: "postgres",
	})

	return db
}
