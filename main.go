package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/mrtomyum/nava-api3/models"
	"log"
	"net/http"
	"github.com/mrtomyum/nava-api3/controllers"
	"github.com/gorilla/mux"
)

//TODO: Move Config to JSON file and Create Config{} to handle DB const.
const (
	//TODO: เมื่อรันจริงต้องเปลี่ยนเป็น Docker Network Bridge IP เช่น 172.17.0.3 เป็นต้น
	DB_HOST = "tcp(nava.work:3306)"
	DB_NAME = "system"
	DB_USER = "root"
	DB_PASS = "mypass"
)

var dsn = DB_USER + ":" + DB_PASS + "@" + DB_HOST + "/" + DB_NAME + "?parseTime=true"

func main() {
	db, err := models.NewDB(dsn)
	if err != nil {
		log.Panic("NewDB() Error:", err)
	}
	c := &controllers.Env{DB: db}
	defer db.Close()
	log.Println("start NewDB()")

	r := mux.NewRouter()

	// User
	r.HandleFunc("/api/v1/user", c.UserAll).Methods("GET")
	log.Println("start '/api/v1/user' GET UserAll")
	r.HandleFunc("/api/v1/user", c.UserInsert).Methods("POST")
	log.Println("start '/api/v1/user' POST UserNew")
	r.HandleFunc("/api/v1/user/{id:[0-9]+}", c.UserShow).Methods("GET")
	log.Println("start '/api/v1/user/:id' GET UserShow")
	r.HandleFunc("/api/v1/user/{id:[0-9]+}", c.UserUpdate).Methods("PUT")
	log.Println("start'/api/v1/user/:id' PUT UserUpdate ")
	r.HandleFunc("/api/v1/user/search", c.UserSearch).Methods("POST")
	log.Println("start '/api/v1/user/search' POST UserSearch")
	r.HandleFunc("/api/v1/login", c.UserLogin).Methods("POST")
	log.Println("start '/api/v1/login' POST UserLogin")
	r.HandleFunc("/api/v1/user/{id:[0-9]+}", c.UserDelete).Methods("DELETE")
	log.Println("start '/api/v1/user/:id' DELETE UserDelete")
	r.HandleFunc("/api/v1/user/{id:[0-9]+}/undelete", c.UserUndelete).Methods("PUT")
	log.Println("start '/api/v1/user/:id/undelete' PUT UserUndelete")

	// Menu
	r.HandleFunc("/api/v1/menu", c.MenuAll).Methods("GET")
	log.Println("start Router GET MenuAll")
	r.HandleFunc("/api/v1/menu", c.MenuInsert).Methods("POST")
	log.Println("start Router POST MenuNew")
	r.HandleFunc("/api/v1/menu/tree", c.MenuTree).Methods("GET")
	log.Println("start Router GET MenuTree")

	http.Handle("/", r)
	http.ListenAndServe(":8000", nil)
	//http.ListenAndServe(":8000", Handlers(c))
}

// Todo: แก้ปัญหา"/" ปิดท้าย URI แล้ว 404 page not found