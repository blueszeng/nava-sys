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

func (e *Env) MainIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%v", "Hello")

}

func (e *Env) UserIndex(w http.ResponseWriter, r *http.Request) {
	log.Println("start userIndex()")
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

