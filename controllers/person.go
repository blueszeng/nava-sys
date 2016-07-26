package controllers

import (
	"net/http"
	"fmt"
	"log"
	"encoding/json"
	"github.com/mrtomyum/nava-api3/models"
	"github.com/mrtomyum/nava-api3/api"
)

func (e *Env) NewPerson(w http.ResponseWriter, r *http.Request) {
	log.Println("call POST New Person()")
	//if r.Method != "POST" {
	//	http.Error(w, http.StatusText(500), 500)
	//	return
	//}
	p := models.Person{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		log.Println("Error decode NewPerson request body>>", err)
	}
	err = p.New(e.DB)
	rs := api.Response{}
	if err != nil {
		rs.Status = "300"
		rs.Message = err.Error()
	} else {
		rs.Status = "201"
		rs.Message = "NEW PERSON CREATED"
		rs.Result = p
	}
	o, _ := json.Marshal(rs)
	fmt.Fprintf(w, "%s", string(o))
}

func (e *Env) AllPerson(w http.ResponseWriter, r *http.Request) {
	log.Println("call GET All Person()")
	p := models.Person{}
	persons, err := p.All(e.DB)
	if err != nil {
		log.Println("Error after call p.All():", err)
	}
	rs := api.Response{}
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		rs.Status = "500"
		rs.Message = err.Error()
	} else {
		rs.Status = "200"
		rs.Message = "OK"
		rs.Result = persons
	}
	output, _ := json.Marshal(rs)
	fmt.Fprintf(w, "%s", string(output))
}

