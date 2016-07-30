package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	m "github.com/mrtomyum/nava-api3/model"
	"github.com/mrtomyum/nava-api3/api"
)

// Method UserShow to query 1 row of user match u.id
func (e Env) ShowUser(w http.ResponseWriter, r *http.Request) {
	log.Println("call GET UserShow()")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") //to allow cross domain AJAX.

	//
	//if r.Method != "GET" {
	//	http.Error(w, http.StatusText(500), 500)
	//	return
	//}
	v := mux.Vars(r)
	id := v["id"]
	u := new(m.User)
	u.ID, _ = strconv.ParseUint(id, 10, 64)
	log.Println("Print u.ID", id)

	user, err := u.Show(e.DB)

	rs := api.Response{}
	if err != nil {
		rs.Status = "204"
		rs.Message = "No Content: " + err.Error()
	} else {
		rs.Status = "200"
		rs.Message = "OK"
		rs.Result = user
	}
	w.WriteHeader(http.StatusOK)
	o, _ := json.Marshal(rs)
	fmt.Fprintf(w, "%s", string(o))
}

func (e Env) UpdateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("call PUT UserUpdate()")
	log.Println("Request Body:", r.Body)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") //to allow cross domain AJAX.

	if r.Method != "PUT" {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	u := new(m.User)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&u)
	if err != nil {
		log.Println("Error decode.Decode(&u) >>", err)
	}
	log.Println("1. After decoder Check decoder, u= ", decoder, u)
	v := mux.Vars(r)
	userID := v["id"]
	u.ID, _ = strconv.ParseUint(userID, 10, 64)
	log.Println("2. Check v, u= ", v, u)

	updateUser, err := u.Update(e.DB)
	fmt.Println("Result User UPDATE to DB: ", updateUser)

	rs := api.Response{}
	if err != nil {
		rs.Status = "304"
		rs.Message = "Not Modified: " + err.Error()
		//rs.Result = updateUser
	} else {
		rs.Status = "200"
		rs.Message = "Update OK"
		rs.Result = updateUser
	}
	w.WriteHeader(http.StatusOK)
	output, _ := json.Marshal(rs)
	fmt.Fprintf(w, "%s", string(output))
}

func (e *Env) AllUser(w http.ResponseWriter, r *http.Request) {
	log.Println("call GET All User()")
	//if r.Method != "GET" {
	//	http.Error(w, http.StatusText(500), 500)
	//	return
	//}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") //to allow cross domain AJAX.


	u := m.User{}
	users, err := u.All(e.DB)
	rs := api.Response{}
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		rs.Status = "500"
		rs.Message = err.Error()
	} else {
		rs.Status = "200"
		rs.Message = "OK"
		rs.Result = users
	}
	output, err := json.Marshal(rs)
	if err != nil {
		log.Println("Error json.Marshal:", err)
	}
	fmt.Fprintf(w, "%s", string(output))
}

func (e *Env) NewUser(w http.ResponseWriter, r *http.Request) {
	log.Println("call POST UserAdd()")
	log.Println("Request Body:", r.Body)

	if r.Method != "POST" {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	u := m.User{}
	// retrieve JSON from body request to decoder and decode it to memory address of User{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&u)
	if err != nil {
		log.Println("Error decode.Decode(&u) >>", err)
	}
	log.Println("Success decode JSON -> :", u, " Result user decoded -> ", u)
	// hash password to []byte before assign to u.Password with function SetPass
	err = u.SetPass()
	if err != nil {
		log.Println("Error u.SetPass(): ", err)
	} else {
		log.Println("Success u.SetPass()")
	}
	// call u.New() method from m/user
	newUser, err := u.New(e.DB)
	rs := api.Response{}
	if err != nil {
		// reply error message with JSON
		rs.Status = "300"
		rs.Message = err.Error()
	} else {
		rs.Status = "201"
		rs.Message = "CREATED"
		rs.Result = newUser
	}
	w.WriteHeader(http.StatusOK)
	output, _ := json.Marshal(rs)
	fmt.Fprintf(w, "%s" ,string(output))
}

// UserDelete Method to mark deleted by field User.DeletedAt.Valid == true
func (e Env) DelUser(w http.ResponseWriter, r *http.Request){
	log.Println("call GET UserDelete() Method:", r.Method)
	if r.Method != "DELETE" {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	v := mux.Vars(r)
	id := v["id"]
	u := new(m.User)
	u.ID, _ = strconv.ParseUint(id, 10, 64)

	err := u.Del(e.DB)
	rs := api.Response{}
	if err != nil {
		rs.Status = "304"
		rs.Message = "Not Modified: " + err.Error()
	} else {
		rs.Status = "200"
		rs.Message = "DELETED OK"
		rs.Result = u
	}
	output, _ := json.Marshal(rs)
	fmt.Fprintf(w, "%s", string(output))
}
// TODO: UserUndelete Method
func (e Env) UndelUser(w http.ResponseWriter, r *http.Request) {
	log.Println("call GET UserUndelete() Method:", r.Method)
	if r.Method != "PUT" {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	v := mux.Vars(r)
	id := v["id"]
	u := new(m.User)
	u.ID, _ = strconv.ParseUint(id, 10, 64)

	err := u.Undel(e.DB)
	rs := api.Response{}
	if err != nil {
		rs.Status = "304"
		rs.Message = "Not Modified" + err.Error()
	} else {
		rs.Status = "200"
		rs.Message = "UNDELETED OK"
		rs.Result = u
	}
	w.WriteHeader(http.StatusOK)
	output, _ := json.Marshal(rs)
	fmt.Fprintf(w, "%s",string(output))

}
// Login Endpoint
func (e Env) LoginUser(w http.ResponseWriter, r *http.Request) {
	log.Println("call POST Login()")
	log.Println("Request Body:", r.Body)

	if r.Method != "POST" {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	u := m.User{}
	decode := json.NewDecoder(r.Body)
	err := decode.Decode(&u)
	if err != nil {
		log.Println("Error decode.Decode(&u) >>", err)
	}
	log.Println("Success decode JSON -> :", u, " Result user decoded -> ", u)

	// Read User.ID, User.Secret by User.Name from DB
	err = u.FindByName(e.DB)
	if err != nil {
		log.Println(err)
	}
	log.Println("u.SearchByName--> user = ", u)

	// Verify Password
	err = u.VerifyPass(u.Password)
	rs := api.Response{}
	if err != nil {
		log.Println(err)
		rs.Status = "500"
		rs.Message = err.Error()
	} else {
		log.Println("Verify Password PASS!!")
		rs.Status = "200"
		rs.Message = "LOGIN SUCCESS"
	}
	w.WriteHeader(http.StatusOK)
	output, _ := json.Marshal(rs)
	fmt.Fprintf(w, "%s", string(output))
}

// UserSearch Method output JSON user.id for client use id as parameter in UserUpdate
func (e Env) SearchUser(w http.ResponseWriter, r *http.Request) {
	log.Println("call GET UserSearch() Method:", r.Method)

	if r.Method != "POST" {
		http.Error(w, http.StatusText(500), 500)
		return
	}

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
		rs.Status = "404"
		rs.Message = "NOT_FOUND ==>" + err.Error()
	} else {
		rs.Status = "302"
		rs.Message = "FOUND"
		rs.Result = users
	}
	output, _ := json.Marshal(rs)
	fmt.Fprintf(w, "%s", string(output))
}