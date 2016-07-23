package models

import (
	"database/sql"
	"log"
)

func NewDB(dsn string) (*sql.DB, error){
	db, err := sql.Open("mysql", dsn)
	if err != nil {

		log.Panic("sql.Open() Error>>", err)
		return nil, err
	}
	if err = db.Ping(); err != nil {
		log.Panic("db.Ping() Error>>", err)
		return nil, err
	}
	log.Println("db = ", db)
	return db, nil //return db so in main can call defer db.Close()
}

