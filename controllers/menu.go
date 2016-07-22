package controllers

import (
	"net/http"
	"fmt"
	"log"
	"encoding/json"
	"github.com/mrtomyum/nava-api3/models"
)

func (e *Env) MenuAll(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	m := new(models.Menu)
	menus, err := m.All(e.DB)
	output, err := json.Marshal(menus)
	if err != nil {
		log.Println("Error json.Marshal:", err)
	}
	fmt.Fprintf(w, string(output))
}

func (e *Env) MenuInsert(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(500), 500)
	}
	m := models.Menu{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&m)
	if err != nil {
		log.Println("Error decode.Decode(&m) >>", err)
	}
	err = m.Insert(e.DB)
	if err != nil {
		fmt.Println("Error Insert DB:", err)
	}
	output, _ := json.Marshal(&m)
	fmt.Fprintf(w, string(output))
}

func (e *Env) MenuTree(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET"{
		http.Error(w, http.StatusText(500), 500)
		return
	}

	m := new(models.Menu)
	menus, _ := m.All(e.DB)
	//if err != nil {
	//	log.Println("Error in models.Menu.All: ", err)
	//}

	jsonNode := new(models.Node)
	for _, menu := range menus{
		n := new(models.Node)
		n.ID = menu.ID
		n.ParentID = menu.ParentID
		n.Text = menu.Name
		n.Path = menu.Path
		n.Note = menu.Note
		jsonNode.Add(n)
	}
	output, _ := json.Marshal(jsonNode)
	fmt.Fprintf(w, string(output))
}