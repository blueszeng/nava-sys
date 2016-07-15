package controllers

import (
	"net/http"
	"fmt"
	"log"
	"encoding/json"
	"database/sql"
	"github.com/mrtomyum/nava-api3/models"
)

type Env struct{
	DB *sql.DB
}

func MainIndex(w http.ResponseWriter, r *http.Request) {
	log.Println("call MainIndex()")
	fmt.Fprintf(w, "%v", "Hello")
}

func TestIndex(w http.ResponseWriter, r *http.Request) {
	log.Println("call TestIndex()")
	fmt.Fprintf(w, "%v", "Test")
}

func (e *Env) UserIndex(w http.ResponseWriter, r *http.Request) {
	log.Println("call GET UserIndex()")

	if r.Method != "GET" {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	u := models.User{}
	users, err := u.All(e.DB)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	o, _ := json.Marshal(users)
	fmt.Fprintf(w, string(o))
}

func (e *Env) UserAdd(w http.ResponseWriter, r *http.Request) {
	log.Println("call POST UserAdd()")
	log.Println("Request Body:", r.Body)

	if r.Method != "POST" {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	// POST JSON must in this form:
	// {id: 999, name: "xxx", password: "xxx', person_id: 999}

	u := models.User{}
	// retrieve JSON from body request to decoder and decode it to memory address of User{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&u)
	if err != nil {
		log.Fatal("Error decode.Decode(&u) >>", err)
	}
	log.Println("Success decode JSON -> :", u, " Result user decoded -> ", u)
	// hash password to []byte before assign to u.Password with function SetPass
	err = u.SetPass()
	if err != nil {
		log.Fatal("Error u.SetPass(): ", err)
	} else {
		log.Println("Success u.SetPass()")
	}
	newUser, err := u.New(e.DB)
	fmt.Println("Result User inserted to DB: ", newUser)
}

// Todo: write Mehtod u.Search by match u.id or u.Name.
// Todo: write Method u.UserShow to query 1 row of user match u.id

func (e Env) UserEdit(w http.ResponseWriter, r *http.Request) {
	log.Println("call POST UserEdit()")
	log.Println("Request Body:", r.Body)

	if r.Method != "PUT" {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	u := models.User{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&u)
	if err != nil {
		log.Fatal("Error decode.Decode(&u) >>", err)
	}
	log.Println("Success decode JSON -> :", u, " Result user decoded -> ", u)
	// hash u.Password
	err = u.SetPass()
	if err != nil {
		log.Fatal("Error u.SetPass(): ", err)
	} else {
		log.Println("Success u.SetPass()")
	}
	editUser, err := u.New(e.DB)
	fmt.Println("Result User inserted to DB: ", editUser)
}

