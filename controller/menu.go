package controller

import (
	"net/http"
	"fmt"
	"log"
	m "github.com/mrtomyum/nava-sys/model"
	"strconv"
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

func (e *Env) GetAllMenuTree(c *gin.Context) {
	log.Println("call GetAllMenuTree()")
	c.Header("Server", "NAVA SYS")
	c.Header("Host", "api.nava.work:8000")
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")
	menu := new(m.Menu)
	rs := api.Response{}
	menus, err := menu.All(e.DB)
	if err != nil {
		fmt.Println("Error Insert DB:", err)
		rs.Status = api.ERROR
		rs.Message = err.Error()
	} else {
		tree := CreateMenuTree(menus)
		rs.Status = api.SUCCESS
		rs.Data = tree.Child
	}
	rs.Link.Self = "api.nava.work:8000/v1/menus/tree"
	c.JSON(http.StatusOK, rs)

}

func (e *Env) UserMenuTree(c *gin.Context) {
	log.Println("UserMenuTree()...")
	c.Header("Server", "NAVA SYS")
	c.Header("Host", "api.nava.work:8000")
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")

	u := new(m.User)
	rs := api.Response{}
	id := c.Param("id")
	var err error
	u.ID, err = strconv.ParseUint(id, 10, 64)
	if err != nil {
		rs.Status = api.FAIL
		rs.Message = err.Error()
	}
	menus, err := u.Menus(e.DB)
	if err != nil {
		fmt.Println("Error Insert DB:", err)
		rs.Status = api.ERROR
		rs.Message = err.Error()
	} else {
		tree := CreateMenuTree(menus)
		rs.Status = api.SUCCESS
		rs.Data = tree.Child
		log.Println(menus)
		log.Println(tree.Child)
	}
	rs.Link.Self = "api.nava.work:8000/v1/menus/tree/users"
	c.JSON(http.StatusOK, rs)
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