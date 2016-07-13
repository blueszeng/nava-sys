package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/mrtomyum/nava-api3/models"
	"net/http"
	"fmt"
	"log"
	"encoding/json"
)

const (
	DB_HOST = "tcp(localhost:3306)"
	DB_NAME = "nava"
	DB_USER = /*"root"*/ "root"
	DB_PASS = /*""*/ "mypass"
)

var dsn = DB_USER + ":" + DB_PASS + "@" + DB_HOST + "/" + DB_NAME + "?charset=utf8"

func main() {
	db := models.InitDB(dsn)
	defer db.Close()
	log.Println("start InitDB()")
	http.HandleFunc("/", mainIndex)
	log.Println("start HandleFunc('/')")
	http.HandleFunc("/users", userIndex)
	log.Println("start HandleFunc('/users')")
	http.ListenAndServe(":8080", nil)

}

func mainIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%v", "Hello")
}

func userIndex(w http.ResponseWriter, r *http.Request) {
	log.Println("start userIndex()")
	if r.Method != "GET" {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	users, err := models.AllUsers()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	o, _ := json.Marshal(users)
	fmt.Fprintf(w, string(o))
}
