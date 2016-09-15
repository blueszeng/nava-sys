package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	m "github.com/mrtomyum/nava-sys/model"
	"github.com/mrtomyum/nava-sys/api"
	"github.com/gin-gonic/gin"
)

// Method UserShow to query 1 row of user match u.id
func (e Env) GetUser(c *gin.Context) {
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
		rs.Link.Self = "host + version + /users"
		c.JSON(http.StatusOK, rs)
	}
}

func (e Env) UpdateUser(c *gin.Context) {
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
		u.ID, err = strconv.ParseUint(id, 10, 64)
		if err != nil {
			rs.Status = api.ERROR
			rs.Message = err.Error()
		}
		updatedUser, err := u.Update(e.DB)
		fmt.Println("Result User UPDATE to DB: ", updatedUser)
		rs := api.Response{}
		if err != nil {
			rs.Status = api.ERROR
			rs.Message = err.Error()
			c.JSON(http.StatusNotModified, rs)
		} else {
			rs.Status = api.SUCCESS
			rs.Data = updatedUser
			c.JSON(http.StatusOK, rs)
		}
	}
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
		c.JSON(http.StatusNotFound, rs)
	} else {
		rs.Status = api.SUCCESS
		rs.Data = users
		c.JSON(http.StatusOK, rs)
	}
}

func (e *Env) NewUser(c *gin.Context) {
	log.Println("call POST UserAdd()")
	c.Header("Server", "NAVA SYS")
	c.Header("Host", "api.nava.work:8000")
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")
	var u m.User
	rs := api.Response{}
	// retrieve JSON from body request to decoder and decode it to memory address of User{}
	if err := c.BindJSON(&u); err != nil {
		rs.Status = api.ERROR
		rs.Message = err.Error()
		c.JSON(http.StatusBadRequest, rs)
	} else {
		// hash password to []byte before assign to u.Password with function SetPass
		err = u.SetPass()
		if err != nil {
			log.Println("Error u.SetPass(): ", err)
		} else {
			log.Println("Success u.SetPass()")
		}
		// call u.New() method from m.user
		newUser, err := u.New(e.DB)
		if err != nil {
			// reply error message with JSON
			rs.Status = api.ERROR
			rs.Message = err.Error()
			c.JSON(http.StatusConflict, rs)
		} else {
			rs.Status = api.SUCCESS
			rs.Data = newUser
			c.JSON(http.StatusOK, rs)
		}
	}
	return
}
// UserDelete Method to mark deleted by field User.DeletedAt.Valid == true
func (e Env) DeleteUser(c *gin.Context){
	log.Println("call GET UserDelete()")
	c.Header("Server", "NAVA SYS")
	c.Header("Host", "api.nava.work:8000")
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")

	id := c.Param("id")
	u := new(m.User)
	u.ID, _ = strconv.ParseUint(id, 10, 64)
	u, err := u.Del(e.DB)
	rs := api.Response{}
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
//  User Undelete Method
func (e Env) UndeleteUser(c *gin.Context) {
	log.Println("call GET UserUndelete()")
	c.Header("Server", "NAVA SYS")
	c.Header("Host", "api.nava.work:8000")
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")

	id := c.Param("id")
	u := new(m.User)
	u.ID, _ = strconv.ParseUint(id, 10, 64)

	u, err := u.Undel(e.DB)
	rs := api.Response{}
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
func (e Env) LoginUser(c *gin.Context) {
	log.Println("call POST Login()")
	c.Header("Server", "NAVA SYS")
	c.Header("Host", "api.nava.work:8000")
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")
	var loginUser m.User
	rs := api.Response{}
	if err := c.BindJSON(&loginUser); err != nil {
		log.Println("Error c.BindJSON(&loginUser) >>", err)
		rs.Status = api.ERROR
		rs.Message = err.Error()
		c.JSON(http.StatusBadRequest, rs)
	} else {
		log.Println("Success decode JSON -> :", loginUser, " Result user decoded -> ", loginUser)
		foundUser, err := loginUser.FindByName(e.DB)
		if err != nil {
			rs.Status = api.ERROR
			rs.Message = err.Error()
			c.JSON(http.StatusUnauthorized, rs)
		}
		// Verify Password
		err = foundUser.VerifyPass(loginUser.Password)
		log.Println ("Check Password:", loginUser.Password)
		if err != nil {
			log.Println("VerifyPass Fail:", err)
			rs.Status = api.ERROR
			rs.Message = err.Error()
			c.JSON(http.StatusUnauthorized, rs)
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
func (e Env) SearchUser(w http.ResponseWriter, r *http.Request) {
	log.Println("call GET UserSearch() Method:", r.Method)

	if r.Method != "POST" {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") //to allow cross domain AJAX.

	// get search string from r.Body
	var s api.Search
	decode := json.NewDecoder(r.Body)
	err := decode.Decode(&s)
	log.Println("Search string 'name'=", s.Name)
	if err != nil {
		log.Println("Error decode.Decode(&u) >>", err)
	}
	log.Println("Success decode JSON -> :", s)

	users, err := m.SearchUsers(e.DB, s.Name)
	if err != nil {
		log.Println("Error in Query:", err)
	}
	rs := api.Response{}
	if users == nil {
		rs.Status = api.ERROR
		rs.Message = "NOT_FOUND ==>" + err.Error()
		w.WriteHeader(http.StatusNotFound)
	} else {
		rs.Status = api.SUCCESS
		rs.Data = users
		w.WriteHeader(http.StatusOK)
	}
	output, _ := json.Marshal(rs)
	fmt.Fprintf(w, "%s", string(output))
}
