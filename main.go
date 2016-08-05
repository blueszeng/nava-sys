package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	c "github.com/mrtomyum/nava-sys/controller"
	m "github.com/mrtomyum/nava-sys/model"
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
	// Read configuration file from "config.json"
	config := loadConfig()
	var dsn = config.DBUser + ":" + config.DBPass + "@" + config.DBHost + "/" + config.DBName + "?parseTime=true"
	// Create new DB connection pool
	db, err := m.NewDB(dsn)
	if err != nil {
		log.Panic("NewDB() Error:", err)
	}
	defer db.Close()
	c := &c.Env{DB: db}
	r := SetupRouter(c)
	http.Handle("/", r)
	http.ListenAndServe(":8000", nil)
}

func loadConfig() *Config {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	config := new(Config)
	err := decoder.Decode(&config)
	if err != nil {
		log.Println("error:", err)
	}
	return config
}

func SetupRouter(c *c.Env) *mux.Router{
	// แก้ปัญหา"/" ปิดท้าย URI แล้ว 404 page not found
	// .StrictSlash(true) help ignore last "/" in URI
	r := mux.NewRouter().StrictSlash(true)
	// User
	s := r.PathPrefix("/v1/users").Subrouter()
	s.HandleFunc("/", c.AllUser).Methods("GET"); log.Println("/v1/users GET AllUser");
	s.HandleFunc("/", c.NewUser).Methods("POST"); log.Println("/v1/users POST New User")
	s.HandleFunc("/{id:[0-9]+}", c.ShowUser).Methods("GET"); log.Println("/v1/users/:id GET Show User")
	s.HandleFunc("/{id:[0-9]+}", c.UpdateUser).Methods("POST"); log.Println("/v1/users/:id POST Update User")
	s.HandleFunc("/search", c.SearchUser).Methods("POST"); log.Println("/v1/users/search POST Search User")
	s.HandleFunc("/login", c.LoginUser).Methods("POST"); log.Println("/v1/login POST UserLogin")
	s.HandleFunc("/{id:[0-9]+}", c.DelUser).Methods("DELETE"); log.Println("/v1/users/:id DELETE UserDelete")
	s.HandleFunc("/{id:[0-9]+}/undel", c.UndelUser).Methods("POST"); log.Println("/v1/users/:id/undel PUT UserUndelete")
	// Menu
	s = r.PathPrefix("/v1/menus").Subrouter()
	s.HandleFunc("/", c.AllMenu).Methods("GET"); log.Println("/v1/menus GET AllMenu")
	s.HandleFunc("/", c.NewMenu).Methods("POST"); log.Println("/v1/menus POST NewMenu")
	s.HandleFunc("/tree", c.AllMenuTree).Methods("GET"); log.Println("/v1/menus/tree GET TreeMenu")
	s.HandleFunc("/tree/users/{id:[0-9]+}", c.UserMenuTree).Methods("GET"); log.Println("/v1/menus/tree/users/:id GET FindMenuTreeByUser")
	// Person
	s = r.PathPrefix("/v1/persons").Subrouter()
	s.HandleFunc("/", c.AllPerson).Methods("GET"); log.Println("/v1/persons GET AllPerson")
	s.HandleFunc("/", c.NewPerson).Methods("POST"); log.Println("/v1/persons POST NewPerson")
	s.HandleFunc("/{id:[0-9]+}", c.ShowPerson).Methods("GET"); log.Println("/v1/persons/:id GET ShowPerson")

	return r
}