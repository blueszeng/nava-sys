package models

import (
	"log"
	"golang.org/x/crypto/bcrypt"
	"database/sql"
)

type User struct {
	ID       int64	`json:id`
	Name     string	`json:name`
	Password []byte	`json:password`
	PeopleID sql.NullInt64
}

type Users []*User

func (u *User) All(db *sql.DB) ([]*User, error) {
	log.Println(">>> start AllUsers() >> db = ", db)
	rows, err := db.Query("SELECT * FROM user")
	if err != nil {
		log.Println(">>> db.Query Error= ", err)
		return nil, err
	}
	defer rows.Close()
	var users Users
	for rows.Next() {
		//u := new(User)
		err := rows.Scan(&u.ID, &u.Name, &u.Password, &u.PeopleID)
		if err != nil {
			log.Println(">>> rows.Scan() Error= ", err)
			rows.Close()
			return nil, err
		}
		users = append(users, u)
	}
	if err = rows.Err(); err != nil {
		log.Println(">>> rows.Err()= ", err)
		rows.Close()
		return nil, err
	}
	rows.Close()
	log.Println("return users", users)
	return users, nil
}

func (u *User) Add(db *sql.DB) (*User, error){
	err := error()
	return u, err
}

func (u *User) SetPass(p string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	u.Password = hash
	return nil
}

func (u *User) VerifyPass(p string) error {
	err := bcrypt.CompareHashAndPassword(u.Password, []byte(p))
	if err != nil {
		return err
	}
	return nil
}