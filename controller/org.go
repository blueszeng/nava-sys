package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"github.com/mrtomyum/nava-sys/model"
	"github.com/mrtomyum/nava-sys/api"
	"net/http"
	"github.com/mrtomyum/nava-sys/config"
)

func (e *Env) GetAllOrg(c *gin.Context) {
	log.Println("call GET All Org")
	c.Header("Server", "NAVA SYS")
	c.Header("Host", "api.nava.work:8000")
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")
	var org model.Org
	orgs, err := org.All(e.DB)
	rs := api.Response{Status:api.SUCCESS, Data:orgs}
	if err != nil {
		rs.Status = api.ERROR
		rs.Message = err.Error()
	}
	rs.Link.Self = config.API_HOST + "/v1/orgs"
	c.JSON(http.StatusOK, rs)
}

func (e *Env) GetOrgRole(c *gin.Context) {
	log.Println("call GET All Org")
	c.Header("Server", "NAVA SYS")
	c.Header("Host", "api.nava.work:8000")
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")
	var org model.Org
	orgs, err := org.All(e.DB)
	rs := api.Response{Status:api.SUCCESS, Data:orgs}
	if err != nil {
		rs.Status = api.ERROR
		rs.Message = err.Error()
	}
	rs.Link.Self = config.API_HOST + "/v1/orgs"
	c.JSON(http.StatusOK, rs)
}

func (e *Env) PostOrg(c *gin.Context) {

}
