package model

import (
	//"database/sql"
	"log"
	"github.com/jmoiron/sqlx"
)

func NewDB(dsn string) (*sqlx.DB, error){
	db := sqlx.MustConnect("mysql", dsn)

	//if err != nil {
	//	log.Panic("sql.Open() Error>>", err)
	//	return nil, err
	//}
	//if err = db.Ping(); err != nil {
	//	log.Panic("db.Ping() Error>>", err)
	//	return nil, err
	//}
	//db.SetConnMaxLifetime(0)
	//db.SetMaxIdleConns(100)
	//db.SetMaxOpenConns(0)
	log.Println("db = ", db)
	return db, nil //return db so in main can call defer db.Close()
}

