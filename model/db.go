package model

import (
	//"database/sql"
	"log"
	"github.com/jmoiron/sqlx"
)

func NewDB(dsn string) (*sqlx.DB, error){
	db := sqlx.MustConnect("mysql", dsn)
	log.Println("db = ", db)
	return db, nil
}

