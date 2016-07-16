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
	DB_USER = /*"root"*/ "root"
	DB_PASS = /*""*/ "mypass"
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

	http.HandleFunc("/api", controllers.MainIndex)
	//log.Println("start HandleFunc('/api')")
	http.HandleFunc("/api/v1/test", controllers.TestIndex)
	//log.Println("start HandleFunc('/api/v1/test')")

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/user", c.UserAll).Methods("GET")
	log.Println("start HandleFunc('/api/v1/users') GET")
	r.HandleFunc("/api/v1/user", c.UserNew).Methods("POST")
	log.Println("start HandleFunc('/api/v1/users') POST")
	r.HandleFunc("/api/v1/user/{id:[0-9]+}", c.UserShow).Methods("GET")
	log.Println("start HandleFunc('/api/v1/users/:id') GET")
	r.HandleFunc("/api/v1/user/{id:[0-9]+}", c.UserUpdate).Methods("PUT")
	log.Println("start HandleFunc('/api/v1/users/:id') PUT")
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}