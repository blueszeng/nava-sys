package model

import (
	//"database/sql"
	"log"
	"github.com/jmoiron/sqlx"
)

func NewDB(dsn string) (*sqlx.DB){
	db := sqlx.MustConnect("mysql", dsn)
	log.Println("db = ", db)
	return db
}

