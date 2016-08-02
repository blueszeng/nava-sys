package controller

import (
	"net/http"
	"fmt"
	"log"
	"encoding/json"
	m "github.com/mrtomyum/nava-api3/model"
)

func (e *Env) MenuAll(w http.ResponseWriter, r *http.Request) {
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

func (e *Env) MenuInsert(w http.ResponseWriter, r *http.Request) {
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

func (e *Env) MenuTree(w http.ResponseWriter, r *http.Request) {
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
	// Setup root node
	jsonNode := new(m.Node)
	jsonNode.ID = menus[0].ID
	jsonNode.ParentID = menus[0].ParentID
	jsonNode.Text = menus[0].Text
	jsonNode.Icon = menus[0].Icon
	jsonNode.SelectedIcon = menus[0].SelectedIcon
	jsonNode.Path = menus[0].Path
	jsonNode.Note = menus[0].Note
	// Adding child node
	for k, menu := range menus{
		n := new(m.Node)
		if k != 0 {
			n.ID = menu.ID
			n.ParentID = menu.ParentID
			n.Text = menu.Text
			n.Icon = menu.Icon
			n.SelectedIcon = menu.SelectedIcon
			n.Path = menu.Path
			n.Note = menu.Note
			log.Println("n=", n)
			jsonNode.Add(n)
		}
	}
	// Todo: how to remove root node from jsonNode?
	w.WriteHeader(http.StatusOK)
	output, _ := json.Marshal(jsonNode)
	fmt.Fprintf(w, string(output))
}