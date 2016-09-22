package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	m "github.com/mrtomyum/nava-sys/model"
	"github.com/mrtomyum/nava-sys/api"
	"github.com/gin-gonic/gin"
	"github.com/mrtomyum/nava-sys/config"
)

// Method UserShow to query 1 row of user match u.id
func (e *Env) GetUser(c *gin.Context) {
	log.Println("call GET UserShow()")
	c.Header("Server", "NAVA SYS")
	c.Header("Host", "api.nava.work:8000")
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")
	id := c.Param("id")
	u := new(m.User)
	u.ID, _ = strconv.ParseUint(id, 10, 64)
	user, err := u.Get(e.DB)
	rs := api.Response{}
	if err != nil {
		rs.Status = api.ERROR
		rs.Message = "No Content: " + err.Error()
		c.JSON(http.StatusNotFound, rs)
	} else {
		rs.Status = api.SUCCESS
		rs.Data = user
		rs.Link.Self = config.API_HOST + "/v1/users/"
		c.JSON(http.StatusOK, rs)
	}
}

func (e *Env) UpdateUser(c *gin.Context) {
	log.Println("call PUT UserUpdate()")
	c.Header("Server", "NAVA SYS")
	c.Header("Host", "api.nava.work:8000")
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")
	u := new(m.User)
	rs := api.Response{}
	if err := c.BindJSON(&u); err != nil {
		log.Println("Error c.BindJSON(&u) >>", err)
	} else {
		id := c.Param("id")
		rs.Link.Self = config.API_HOST + "/v1/users/" + id
		u.ID, err = strconv.ParseUint(id, 10, 64)
		if err != nil {
			rs.Status = api.FAIL
			rs.Message = err.Error()
		}
		updatedUser, err := u.Update(e.DB)
		fmt.Println("Result User UPDATE to DB: ", updatedUser)
		rs := api.Response{}
		if err != nil {
			rs.Status = api.ERROR
			rs.Message = err.Error()
		} else {
			rs.Status = api.SUCCESS
			rs.Data = updatedUser
		}
	}
	c.JSON(http.StatusOK, rs)
}

func (e *Env) AllUser(c *gin.Context) {
	log.Println("call GET All User()")
	c.Header("Server", "NAVA SYS")
	c.Header("Host", "api.nava.work:8000")
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")
	u := m.User{}
	users, err := u.All(e.DB)
	rs := api.Response{}
	if err != nil {
		// reply error message with JSON
		rs.Status = api.ERROR
		rs.Message = err.Error()
	} else {
		rs.Status = api.SUCCESS
		rs.Data = users
	}
	rs.Link.Self = config.API_HOST + "/v1/users"
	c.JSON(http.StatusOK, rs)
}

func (e *Env) NewUser(c *gin.Context) {
	log.Println("call POST UserAdd()")
	c.Header("Server", "NAVA SYS")
	c.Header("Host", "api.nava.work:8000")
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")
	var u m.User
	rs := api.Response{}
	rs.Link.Self = "api.nava.work:8000/v1/users"
	// retrieve JSON from body request to decoder and decode it to memory address of User{}
	if err := c.BindJSON(&u); err != nil {
		rs.Status = api.FAIL
		rs.Message = err.Error()
		c.JSON(http.StatusBadRequest, rs)
		return
	} else {
		// hash password to []byte before assign to u.Password with function SetPass
		err = u.SetPass()
		if err != nil {
			log.Println("Error u.SetPass(): ", err)
		} else {
			log.Println("Success u.SetPass()")
		}
		// call u.New() method from m.user
		newUser, err := u.Insert(e.DB)
		if err != nil {
			// reply error message with JSON
			rs.Status = api.ERROR
			rs.Message = err.Error()
			c.JSON(http.StatusConflict, rs)
		} else {
			rs.Status = api.SUCCESS
			rs.Link.Related = config.API_HOST + "/v1/users/"+ string(newUser.ID)
			c.JSON(http.StatusOK, rs)
		}
		return
	}
}
// UserDelete Method to mark deleted by field User.DeletedAt.Valid == true
func (e *Env) DeleteUser(c *gin.Context){
	log.Println("call GET UserDelete()")
	c.Header("Server", "NAVA SYS")
	c.Header("Host", "api.nava.work:8000")
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")
	rs := api.Response{}
	id := c.Param("id")
	rs.Link.Self = "api.nava.work:8000/v1/users/" + id
	u := new(m.User)
	u.ID, _ = strconv.ParseUint(id, 10, 64)
	err := u.Delete(e.DB)
	if err != nil {
		// reply error message with JSON
		rs.Status = api.ERROR
		rs.Message = err.Error()
		c.JSON(http.StatusConflict, rs)
	} else {
		rs.Status = api.SUCCESS
		rs.Data = u
		c.JSON(http.StatusOK, rs)
	}
}
//------------------------
//  User UnDelete Method
//------------------------
func (e *Env) UndeleteUser(c *gin.Context) {
	log.Println("call GET UserUndelete()")
	c.Header("Server", "NAVA SYS")
	c.Header("Host", "api.nava.work:8000")
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")

	id := c.Param("id")
	u := new(m.User)
	u.ID, _ = strconv.ParseUint(id, 10, 64)

	rs := api.Response{}
	rs.Link.Self = "api.nava.work:8000/v1/users/undelete/" + id
	u, err := u.Undel(e.DB)
	if err != nil {
		rs.Status = api.ERROR
		rs.Message = "Not Modified" + err.Error()
		c.JSON(http.StatusNotModified, rs)
	} else {
		rs.Status = api.SUCCESS
		rs.Data = u
		c.JSON(http.StatusOK, rs)
	}
}
// Login Endpoint
func (e *Env) Login(c *gin.Context) {
	log.Println("call POST Login()")
	c.Header("Server", "NAVA SYS")
	c.Header("Host", "api.nava.work:8000")
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")
	var loginUser m.User
	rs := api.Response{}
	rs.Link.Self = "api.nava.work:8000/v1/users/login"
	if err := c.BindJSON(&loginUser); err != nil {
		log.Println("Error c.BindJSON(&loginUser) >>", err)
		rs.Status = api.FAIL
		rs.Message = err.Error()
		c.JSON(http.StatusBadRequest, rs)
	} else {
		log.Println("Success decode JSON -> :", loginUser)
		foundUser, err := loginUser.FindByName(e.DB)
		if err != nil {
			rs.Status = api.ERROR
			rs.Message = "Cannot find user or user deleted: " + err.Error()
			c.JSON(http.StatusUnauthorized, rs)
			return
		}
		// Verify Password
		err = foundUser.VerifyPass(loginUser.Password)
		log.Println ("Check Password:", loginUser.Password)
		if err != nil {
			log.Println("VerifyPass Fail:", err)
			rs.Status = api.ERROR
			rs.Message = "Wrong username or password:" + err.Error()
			c.JSON(http.StatusUnauthorized, rs)
			c.AbortWithStatus(401)
		} else {
			// Make UserPermission for response
			p, err := foundUser.Permission(e.DB)
			if err != nil {
				log.Println("Error from call user.Permission()", err)
			}
			log.Println("User Permission: ", p)
			rs.Status = api.SUCCESS
			rs.Data = p
			id := strconv.FormatUint(foundUser.ID, 10)
			rs.Link.Related = "http://api.nava.work:8000/v1/menus/tree/users/" + id
			c.JSON(http.StatusOK, rs)
		}
	}
}

// UserSearch Method output JSON user.id for client use id as parameter in UserUpdate
func (e *Env) SearchUser(c *gin.Context) {
	log.Println("call GET UserSearch()")
	c.Header("Server", "NAVA SYS")
	c.Header("Host", "api.nava.work:8000")
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")

	var s api.Search
	if err := c.BindJSON(&s); err != nil {
		log.Println("Error decode.Decode(&u) >>", err)

	} else {
		log.Println("Success decode JSON -> :", s)
		log.Println("Search string 'name'=", s.Name)
		users, err := m.SearchUsers(e.DB, s.Name)
		if err != nil {
			log.Println("Error in Query:", err)
		}
		rs := api.Response{}
		rs.Link.Self = "api.nava.work:8000/v1/users/search"
		if users == nil {
			rs.Status = api.ERROR
			rs.Message = "NOT_FOUND ==>" + err.Error()
			c.JSON(http.StatusNotFound, rs)
		} else {
			rs.Status = api.SUCCESS
			rs.Data = users
			c.JSON(http.StatusOK, rs)
		}
	}
}

func (e *Env) GetUserOrg(c *gin.Context) {
	log.Println("call GetUserOg()")
	c.Header("Server", "NAVA SYS")
	c.Header("Host", "api.nava.work:8000")
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")

	id := c.Param("id")
	u := new(m.User)
	u.ID, _ = strconv.ParseUint(id, 10, 64)
	org, err := u.GetOrg(e.DB)
	rs := api.Response{}
	if err != nil {
		rs.Status = api.ERROR
		rs.Message = "No Content: " + err.Error()
		c.JSON(http.StatusNotFound, rs)
	} else {
		rs.Status = api.SUCCESS
		rs.Data = org
		rs.Link.Self = config.API_HOST + "/v1/users/org"
		c.JSON(http.StatusOK, rs)
	}
}