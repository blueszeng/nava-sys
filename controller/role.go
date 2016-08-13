package controller

import (
	"net/http"
	"fmt"
	"log"
	"encoding/json"
	"github.com/mrtomyum/nava-sys/model"
)

func (e *Env) AllRole(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") //to allow cross domain AJAX.


	role := new(model.Role)
	roles, err := role.All(e.DB)
	w.WriteHeader(http.StatusOK)
	output, err := json.Marshal(roles)
	if err != nil {
		log.Println("Error json.Marshal:", err)
	}
	fmt.Fprintf(w, string(output))
}
