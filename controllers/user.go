package controllers

import (
	"net/http"
	"fmt"
	"log"
	"encoding/json"
	"database/sql"
	"github.com/mrtomyum/nava-api3/models"
	"github.com/gorilla/mux"
	"strconv"
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

func (e *Env) UserAll(w http.ResponseWriter, r *http.Request) {
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
	output, _ := json.Marshal(users)
	fmt.Fprintf(w, string(output))
}

func (e *Env) UserNew(w http.ResponseWriter, r *http.Request) {
	log.Println("call POST UserAdd()")
	log.Println("Request Body:", r.Body)

	if r.Method != "POST" {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	// POST JSON must in this form:
	// {"name": "xxx", "password": "xxx'}

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
	// call u.New() method from models/user
	newUser, err := u.New(e.DB)
	output, _ := json.Marshal(newUser)
	fmt.Fprintf(w, string(output))
	fmt.Println("Result User inserted to DB: ", newUser)
}

// Todo: write Method UserSearch by match u.id or u.Name.

// Method UserShow to query 1 row of user match u.id
func (e Env) UserShow(w http.ResponseWriter, r *http.Request) {
	log.Println("call GET UserShow()")

	if r.Method != "GET" {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	v := mux.Vars(r)
	userID := v["id"]
	u := new(models.User)
	u.ID, _ = strconv.ParseInt(userID, 10, 64)
	n, err := u.Show(e.DB)
	if err != nil {
		log.Fatal("Error u.Show in c.user.go.Show:", err)
	}
	output, _ := json.Marshal(n)
	fmt.Fprintf(w, string(output))
}

// Todo: Method UserSearch output JSON user.id for client use id as parameter in UserUpdate

func (e Env) UserUpdate(w http.ResponseWriter, r *http.Request) {
	log.Println("call PUT UserUpdate()")
	log.Println("Request Body:", r.Body)

	if r.Method != "PUT" {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	v := mux.Vars(r)
	userID := v["id"]
	u := new(models.User)
	u.ID, _ = strconv.ParseInt(userID, 10, 64)

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
	updateUser, err := u.Update(e.DB)
	fmt.Println("Result User UPDATE to DB: ", updateUser)
	output, _ := json.Marshal(updateUser)
	fmt.Fprintf(w, string(output))
}

// Todo: write Method UserDel to mark deleted by field del_date
