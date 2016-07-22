package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/mrtomyum/nava-api3/models"
	"log"
	"net/http"
	"github.com/mrtomyum/nava-api3/controllers"
	"github.com/gorilla/mux"
)

const (
	DB_HOST = "tcp(nava.work:3306)"
	DB_NAME = "nava"
	DB_USER = "root"
	DB_PASS = "mypass"
)

var dsn = DB_USER + ":" + DB_PASS + "@" + DB_HOST + "/" + DB_NAME + "?parseTime=true"

func main() {
	db, err := models.NewDB(dsn)
	if err != nil {
		log.Panic("NewDB() Error:", err)
	}

	c := &controllers.Env{DB: db}

	defer db.Close()
	log.Println("start NewDB()")

	r := mux.NewRouter()

	// User
	r.HandleFunc("/api/v1/user", c.UserAll).Methods("GET")
	log.Println("start '/api/v1/user' GET UserAll")
	r.HandleFunc("/api/v1/user", c.UserInsert).Methods("POST")
	log.Println("start '/api/v1/user' POST UserNew")
	r.HandleFunc("/api/v1/user/{id:[0-9]+}", c.UserShow).Methods("GET")
	log.Println("start '/api/v1/user/:id' GET UserShow")
	r.HandleFunc("/api/v1/user/{id:[0-9]+}", c.UserUpdate).Methods("PUT")
	log.Println("start'/api/v1/user/:id' PUT UserUpdate ")
	r.HandleFunc("/api/v1/user/search", c.UserSearch).Methods("POST")
	log.Println("start '/api/v1/user/search' POST UserSearch")
	r.HandleFunc("/api/v1/login", c.UserLogin).Methods("POST")
	log.Println("start '/api/v1/login' POST UserLogin")
	r.HandleFunc("/api/v1/user/{id:[0-9]+}", c.UserDelete).Methods("DELETE")
	log.Println("start '/api/v1/user/:id' DELETE UserDelete")
	r.HandleFunc("/api/v1/user/{id:[0-9]+}/undelete", c.UserUndelete).Methods("PUT")
	log.Println("start '/api/v1/user/:id/undelete' PUT UserUndelete")


	// # Stock

	// ## Item

	// ## Location

	http.Handle("/", r)
	http.ListenAndServe(":8000", nil)
}