package model

import (
	"testing"
	"github.com/jmoiron/sqlx"
	"github.com/mrtomyum/nava-sys/config"
)

var db *sqlx.DB

func init() {
	dsn := config.LoadDSN("config.json")
	db = NewDB(dsn)
	defer db.Close()
}

func TestUser_Get(t *testing.T) {
	var u User
	user, err := u.Get(db)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(user)
}

func TestUser_All(t *testing.T) {

}