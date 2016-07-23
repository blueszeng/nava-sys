package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/mrtomyum/nava-api3/controllers"
	"github.com/mrtomyum/nava-api3/models"
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
	db, err := models.NewDB(dsn)
	if err != nil {
		log.Panic("NewDB() Error:", err)
	}
	c := &controllers.Env{DB: db}
	defer db.Close()

	r := SetupRouter(c)

	http.Handle("/", r)
	http.ListenAndServe(":8000", nil)
}

func SetupRouter(c *controllers.Env) *mux.Router{
	// แก้ปัญหา"/" ปิดท้าย URI แล้ว 404 page not found
	// .StrictSlash(true) help ignore last "/" in URI
	r := mux.NewRouter().StrictSlash(true)
	// User
	s := r.PathPrefix("/api/v1/user").Subrouter()
	s.HandleFunc("/", c.UserIndex).Methods("GET")
	log.Println("/api/v1/index GET UserIndex")
	s.HandleFunc("/", c.UserInsert).Methods("POST")
	log.Println("/api/v1/user POST UserInsert")
	s.HandleFunc("/{id:[0-9]+}", c.UserShow).Methods("GET")
	log.Println("/api/v1/user/:id GET UserShow")
	s.HandleFunc("/{id:[0-9]+}", c.UserUpdate).Methods("PUT")
	log.Println("/api/v1/user/:id PUT UserUpdate ")
	s.HandleFunc("/search", c.UserSearch).Methods("POST")
	log.Println("/api/v1/user/search POST UserSearch")
	s.HandleFunc("/login", c.UserLogin).Methods("POST")
	log.Println("/api/v1/login POST UserLogin")
	s.HandleFunc("/{id:[0-9]+}", c.UserDelete).Methods("DELETE")
	log.Println("start '/api/v1/user/:id' DELETE UserDelete")
	s.HandleFunc("/{id:[0-9]+}/undelete", c.UserUndelete).Methods("PUT")
	log.Println("start '/api/v1/user/:id/undelete' PUT UserUndelete")
	// Menu
	s = r.PathPrefix("/api/v1/menu").Subrouter()
	s.HandleFunc("/", c.MenuAll).Methods("GET")
	log.Println("start Router GET MenuAll")
	s.HandleFunc("/", c.MenuInsert).Methods("POST")
	log.Println("start Router POST MenuNew")
	s.HandleFunc("/tree", c.MenuTree).Methods("GET")
	log.Println("start Router GET MenuTree")

	return r
}