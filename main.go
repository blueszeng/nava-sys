package main

import (
	_ "github.com/go-sql-driver/mysql"
	c "github.com/mrtomyum/nava-sys/controller"
	m "github.com/mrtomyum/nava-sys/model"
	"github.com/gin-gonic/gin"
	"github.com/mrtomyum/nava-sys/config"
)

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
		userV1.GET("/:id/orgs", e.GetUserOrg)
	}

	//personV1 := app.Group("/v1/persons")
	//{
	//	personV1.POST("/", e.NewPerson)
		//personV1.GET("/", e.GetAllPerson)
		//personV1.GET("/:id", e.GetPerson)
	//}

	roleV1 := app.Group("/v1/roles")
	{
		roleV1.GET("/", e.GetAllRole)
		roleV1.POST("/", c.NotImplemented)
		roleV1.DELETE("/:id", c.NotImplemented)
	}

	menuV1 := app.Group("/v1/menus")
	{
		menuV1.GET("/", e.GetAllMenu)
		menuV1.POST("/", e.PostNewMenu)
		menuV1.GET("/tree", e.GetAllMenuTree)
		menuV1.GET("/tree/users/:id", e.UserMenuTree)
		menuV1.DELETE("/:id", c.NotImplemented)
	}

	orgV1 := app.Group("/v1/orgs")
	{
		orgV1.GET("/", e.GetAllOrg)
	}
	return app
}

func main() {
	// Read configuration file from "config.json"
	dsn := config.LoadDSN("config.json")
	// Create new DB connection pool
	db := m.NewDB(dsn)
	defer db.Close()
	e := &c.Env{DB: db}
	app := SetupRouter(e)
	app.Run(":8000")
}