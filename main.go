package main

import (
	_ "github.com/go-sql-driver/mysql"
	c "github.com/mrtomyum/nava-sys/controller"
	m "github.com/mrtomyum/nava-sys/model"
	"log"
	"os"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

//TODO: เมื่อรันจริงต้องเปลี่ยนเป็น Docker Network Bridge IP เช่น 172.17.0.3 เป็นต้น
type Config struct {
	DBHost string `json:"db_host"`
	DBName string `json:"db_name"`
	DBUser string `json:"db_user"`
	DBPass string `json:"db_pass"`
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

func SetupRouter(e *c.Env) *gin.Engine{
	app := gin.New()
	app.Use(gin.Logger())
	app.Use(gin.Recovery())
	userV1 := app.Group("/v1/users")
	{
		userV1.POST("/", e.NewUser)
		userV1.GET("/", e.AllUser)
		userV1.GET("/:id", e.GetUser)
		userV1.PUT("/:id", e.UpdateUser)
		userV1.DELETE("/:id", e.DeleteUser)
		userV1.POST("/undelete/:id", e.UndeleteUser)
		userV1.POST("/login", e.LoginUser)
		userV1.POST("/search", e.SearchUser)
	}

	//personV1 := app.Group("/v1/persons")
	//{
	//	personV1.POST("/", e.NewPerson)
		//personV1.GET("/", e.GetAllPerson)
		//personV1.GET("/:id", e.GetPerson)
	//}

	//roleV1 := app.Group("/v1/roles")
	//{
		//roleV1.GET("/", e.AllRole)
	//}

	menuV1 := app.Group("/v1/menus")
	{
		menuV1.GET("/", e.GetAllMenu)
		menuV1.POST("/", e.PostNewMenu)
		menuV1.GET("/tree", e.GetAllMenuTree)
		menuV1.GET("/tree/users/:id", e.UserMenuTree)
	}
	return app
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
	e := &c.Env{DB: db}
	app := SetupRouter(e)
	app.Run(":8000")
}