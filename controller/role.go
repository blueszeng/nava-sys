package controller

import (
	"net/http"
	"github.com/mrtomyum/nava-sys/model"
	"github.com/gin-gonic/gin"
	"github.com/mrtomyum/nava-sys/api"
)

func (e *Env) GetAllRole(c *gin.Context) {
	c.Header("Server", "NAVA SYS")
	c.Header("Host", "api.nava.work:8000")
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")
	var rs api.Response
	role := new(model.Role)
	roles, err := role.All(e.DB)
	if err != nil {
		rs.Status = api.ERROR
		rs.Message = err.Error()
	} else {
		rs.Status = api.SUCCESS
		rs.Self = "http://api.nava.work:8000/v1/roles"
		rs.Data = roles
	}
	c.JSON(http.StatusOK, rs)
}
