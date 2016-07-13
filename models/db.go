package models

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var db *sql.DB

func InitDB(dsn string) *sql.DB{
	var err error
	db, err = sql.Open("mysql", dsn)
	if err = db.Ping(); err != nil {
		log.Panic(err)
	}
	log.Println("db = ", db)
	return db //return db so in main can call defer db.Close()
}
