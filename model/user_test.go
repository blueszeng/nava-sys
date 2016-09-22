package model

import (
	"testing"
	"github.com/jmoiron/sqlx"
	"github.com/mrtomyum/nava-sys/config"
)

var (
	db *sqlx.DB
	u User
)

func init() {
	dsn := config.LoadDSN("../config.json")
	db = NewDB(dsn)
}

func TestUser_GetById(t *testing.T) {
	u.ID = 1
	user, err := u.Get(db)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(user)
}

func TestUser_All(t *testing.T) {
	users, err := u.All(db)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(users)
}

func TestUser_InsertSuccessAndDelete(t *testing.T){
	// Setup. Arrange
	users := MockUsers()
	for _, u := range users {
		newUser, err := u.Insert(db)
		if err != nil {
			t.Error(err.Error())
		}
		t.Log("Success Inserted User", newUser)
	}

	// Tear down. ToDo: use script delete Table User
	// Loop for range users
	//t.Logf("Tear Down: User %s  deleted", newUser.Name)
}

