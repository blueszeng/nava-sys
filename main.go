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
	DB_HOST = "tcp(172.17.0.3:3306)"
	DB_NAME = "nava"
	DB_USER = "root"
	DB_PASS = "mypass"
)

var dsn = DB_USER + ":" + DB_PASS + "@" + DB_HOST + "/" + DB_NAME + "?charset=utf8"

func main() {
	db, err := models.NewDB(dsn)
	if err != nil {
		log.Panic("NewDB() Error:", err)
	}

	c := &controllers.Env{DB: db}

	defer db.Close()
	log.Println("start NewDB()")

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/user", c.UserAll).Methods("GET")
	log.Println("start HandleFunc('/api/v1/users') GET UserAll")
	r.HandleFunc("/api/v1/user", c.UserInsert).Methods("POST")
	log.Println("start HandleFunc('/api/v1/users') POST UserNew")
	r.HandleFunc("/api/v1/user/{id:[0-9]+}", c.UserShow).Methods("GET")
	log.Println("start HandleFunc('/api/v1/users/:id') GET UserShow")
	r.HandleFunc("/api/v1/user", c.UserUpdate).Methods("PUT")
	log.Println("start HandleFunc('/api/v1/users/:id') PUT UserUpdate ")
	r.HandleFunc("/api/v1/user/search", c.UserSearch).Methods("POST")
	log.Println("start HandleFunc('/api/v1/user/search') POST UserSearch")
	r.HandleFunc("/api/v1/login", c.UserLogin).Methods("POST")
	log.Println("start HandleFunc('/api/v1/login') POST UserLogin")
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}