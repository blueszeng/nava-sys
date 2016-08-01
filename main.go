package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	c "github.com/mrtomyum/nava-api3/controller"
	m "github.com/mrtomyum/nava-api3/model"
	"log"
	"net/http"
	"os"
	"encoding/json"
)

//TODO: เมื่อรันจริงต้องเปลี่ยนเป็น Docker Network Bridge IP เช่น 172.17.0.3 เป็นต้น
type Config struct {
	DBHost string `json:"db_host"`
	DBName string `json:"db_name"`
	DBUser string `json:"db_user"`
	DBPass string `json:"db_pass"`
}

func main() {
	// Read configuration file from "cofig.json"
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	config := Config{}
	err := decoder.Decode(&config)
	if err != nil {
		log.Println("error:", err)
	}

	var dsn = config.DBUser + ":" + config.DBPass + "@" + config.DBHost + "/" + config.DBName + "?parseTime=true"

	// Create new DB connection pool
	db, err := m.NewDB(dsn)
	if err != nil {
		log.Panic("NewDB() Error:", err)
	}

	c := &c.Env{DB: db}
	defer db.Close()

	r := SetupRouter(c)

	http.Handle("/", r)
	http.ListenAndServe(":8000", nil)
}

func SetupRouter(c *c.Env) *mux.Router{
	// แก้ปัญหา"/" ปิดท้าย URI แล้ว 404 page not found
	// .StrictSlash(true) help ignore last "/" in URI
	r := mux.NewRouter().StrictSlash(true)
	// User
	s := r.PathPrefix("/v1/user").Subrouter()
	s.HandleFunc("/", c.AllUser).Methods("GET")
	log.Println("/v1/index GET UserIndex")
	s.HandleFunc("/", c.NewUser).Methods("POST")
	log.Println("/v1/user POST UserInsert")
	s.HandleFunc("/{id:[0-9]+}", c.ShowUser).Methods("GET")
	log.Println("/v1/user/:id GET UserShow")
	s.HandleFunc("/{id:[0-9]+}", c.UpdateUser).Methods("PUT")
	log.Println("/v1/user/:id PUT UserUpdate ")
	s.HandleFunc("/search", c.SearchUser).Methods("POST")
	log.Println("/v1/user/search POST UserSearch")
	s.HandleFunc("/login", c.LoginUser).Methods("POST")
	log.Println("/v1/login POST UserLogin")
	s.HandleFunc("/{id:[0-9]+}", c.DelUser).Methods("DELETE")
	log.Println("/v1/user/:id DELETE UserDelete")
	s.HandleFunc("/{id:[0-9]+}/undelete", c.UndelUser).Methods("PUT")
	log.Println("/v1/user/:id/undelete PUT UserUndelete")
	// Menu
	s = r.PathPrefix("/v1/menu").Subrouter()
	s.HandleFunc("/", c.MenuAll).Methods("GET")
	log.Println("/v1/menu GET AllMenu")
	s.HandleFunc("/", c.MenuInsert).Methods("POST")
	log.Println("/v1/menu POST NewMenu")
	s.HandleFunc("/tree", c.MenuTree).Methods("GET")
	log.Println("/v1/menu/tree GET TreeMenu")
	// Person
	s = r.PathPrefix("/v1/person").Subrouter()
	s.HandleFunc("/", c.AllPerson).Methods("GET")
	log.Println("/v1/person GET AllPerson")
	s.HandleFunc("/", c.NewPerson).Methods("POST")
	log.Println("/v1/person POST NewPerson")
	s.HandleFunc("/{id:[0-9]+}", c.ShowPerson).Methods("GET")
	log.Println("/v1/person/:id GET ShowPerson")

	return r
}