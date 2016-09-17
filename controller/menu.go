package controller

import (
	"net/http"
	"fmt"
	"log"
	"encoding/json"
	m "github.com/mrtomyum/nava-sys/model"
	"strconv"
	"github.com/gorilla/mux"
	"github.com/gin-gonic/gin"
	"github.com/mrtomyum/nava-sys/api"
)

func (e *Env) GetAllMenu(c *gin.Context) {
	log.Println("call GetAllMenu()")
	c.Header("Server", "NAVA SYS")
	c.Header("Host", "api.nava.work:8000")
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")
	m := new(m.Menu)
	rs := api.Response{}
	menus, err := m.All(e.DB)
	if err != nil {
		log.Println("Error m.All():", err)
		rs.Status = api.ERROR
		rs.Message = err.Error()
		c.JSON(http.StatusNotFound, rs)
	}
	rs.Status = api.SUCCESS
	rs.Link.Self = "api.nava.work:8000/v1/menus"
	rs.Data = menus
	c.JSON(http.StatusOK, rs)
}

func (e *Env) PostNewMenu(c *gin.Context) {
	log.Println("call GetAllMenu()")
	c.Header("Server", "NAVA SYS")
	c.Header("Host", "api.nava.work:8000")
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")
	m := m.Menu{}
	rs := api.Response{}
	if err := c.BindJSON(&m); err != nil{
		log.Println("Error decode.Decode(&m) >>", err)
		rs.Status = api.FAIL
		rs.Message = err.Error()
	} else {
		newMenu, err := m.Insert(e.DB)
		if err != nil {
			fmt.Println("Error Insert DB:", err)
			rs.Status = api.ERROR
			rs.Message = err.Error()
		} else {
			rs.Status = api.SUCCESS
			rs.Data = newMenu
		}
	}
	rs.Link.Self = "api.nava.work:8000/v1/menus"
	c.JSON(http.StatusOK, rs)
}

func (e *Env) GetAllMenuTree(w http.ResponseWriter, r *http.Request) {
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