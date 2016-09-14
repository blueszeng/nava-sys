package controller

import (
	"net/http"
	"fmt"
	"log"
	"encoding/json"
	m "github.com/mrtomyum/nava-sys/model"
	"strconv"
	"github.com/gorilla/mux"
)

func (e *Env) AllMenu(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") //to allow cross domain AJAX.

	m := new(m.Menu)
	menus, err := m.All(e.DB)
	if err != nil {
		log.Println("Error m.All():", err)
		w.WriteHeader(http.StatusNotFound)
	}
	w.WriteHeader(http.StatusOK)
	output, err := json.Marshal(menus)
	if err != nil {
		log.Println("Error json.Marshal:", err)
	}
	fmt.Fprintf(w, string(output))
}

func (e *Env) NewMenu(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") //to allow cross domain AJAX.
	m := m.Menu{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&m)
	if err != nil {
		log.Println("Error decode.Decode(&m) >>", err)
	}
	err = m.Insert(e.DB)
	if err != nil {
		fmt.Println("Error Insert DB:", err)
		w.WriteHeader(http.StatusNotImplemented)
	}
	w.WriteHeader(http.StatusOK)
	output, _ := json.Marshal(&m)
	fmt.Fprintf(w, string(output))
}

func (e *Env) AllMenuTree(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET"{
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") //to allow cross domain AJAX.

	menu := new(m.Menu)
	menus, err := menu.All(e.DB)
	if err != nil {
		log.Println("Error in m.Menu.All: ", err)
	}
	tree := CreateMenuTree(menus)
	w.WriteHeader(http.StatusOK)
	output, _ := json.Marshal(tree.Child) // remove root node from tree
	fmt.Fprintf(w, string(output))
}

func (e *Env) UserMenuTree(w http.ResponseWriter, r *http.Request) {
	log.Println("UserMenuTree()...")
	//if r.Method != "GET"{
	//	http.Error(w, http.StatusText(500), 500)
	//	return
	//}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	v := mux.Vars(r)
	id := v["id"]
	u := new(m.User)
	u.ID, _ = strconv.ParseUint(id, 10, 64)

	menus, err := u.Menus(e.DB)
	if err != nil {
		log.Fatal("Error call u.FindMenuByUser()", err)
		w.WriteHeader(http.StatusNotFound)
	}

	log.Println(menus)
	tree := CreateMenuTree(menus)
	log.Println(tree)
	w.WriteHeader(http.StatusOK)
	output, _ := json.Marshal(tree.Child)
	fmt.Fprintf(w, string(output))
}

func CreateMenuTree(menus []*m.Menu) *m.Menu {
	// Setup root node
	tree := new(m.Menu)
	// Adding child node
	for _, m := range menus {
		tree.Add(m)
	}
	return tree
}