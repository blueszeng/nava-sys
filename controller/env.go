package controller

import (
	"github.com/jmoiron/sqlx"
	"github.com/gin-gonic/gin"
)

type Env struct {
	DB *sqlx.DB
}

func NotImplemented(c *gin.Context) {
	c.Header("Server", "NAVA SYS")
	c.Header("Host", "api.nava.work:8000")
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(200, "Not Implemented")
}