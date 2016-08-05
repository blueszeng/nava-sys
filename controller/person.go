package controller

import (
	"net/http"
	"fmt"
	"log"
	"encoding/json"
	m "github.com/mrtomyum/nava-api3/model"
	"github.com/mrtomyum/nava-api3/api"
	"github.com/gorilla/mux"
	"strconv"
)

func (e *Env) NewPerson(w http.ResponseWriter, r *http.Request) {
	log.Println("call POST New Person()")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") //to allow cross domain AJAX.
	if r.Method != "POST" {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	p := m.Person{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		log.Println("Error decode NewPerson request body>>", err)
	}
	err = p.New(e.DB)
	rs := api.Response{}
	if err != nil {
		rs.Status = api.ERROR
		rs.Message = err.Error()
	} else {
		rs.Status = api.SUCCESS
		rs.Message = "NEW PERSON CREATED"
		rs.Data = p
	}
	w.WriteHeader(http.StatusOK)
	o, _ := json.Marshal(rs)
	fmt.Fprintf(w, "%s", string(o))
}

func (e *Env) AllPerson(w http.ResponseWriter, r *http.Request) {
	log.Println("call GET All Person()")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") //to allow cross domain AJAX.

	p := m.Person{}
	persons, err := p.All(e.DB)
	if err != nil {
		log.Println("Error after call p.All():", err)
	}
	rs := api.Response{}
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		rs.Status = api.ERROR
		rs.Message = err.Error()
	} else {
		rs.Status = api.SUCCESS
		rs.Message = "OK"
		rs.Data = persons
	}
	w.WriteHeader(http.StatusOK)
	output, _ := json.Marshal(rs)
	fmt.Fprintf(w, "%s", string(output))
}

func (e *Env) ShowPerson(w http.ResponseWriter, r *http.Request) {
	log.Println("call GET Show Person(:id)")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") //to allow cross domain AJAX.

	v := mux.Vars(r)
	id := v["id"]
	p := m.Person{}
	p.ID, _ = strconv.ParseUint(id, 10, 64)

	person, err := p.Show(e.DB)
	if err != nil {
		log.Println("Error after call p.All():", err)
	}
	rs := api.Response{}
	if err != nil {
		rs.Status = api.FAIL
		rs.Message = "No Content: " + err.Error()
	} else {
		rs.Status = api.SUCCESS
		rs.Message = "OK"
		rs.Data = person
	}
	w.WriteHeader(http.StatusOK)
	output, _ := json.Marshal(rs)
	fmt.Fprintf(w, "%s", string(output))
}