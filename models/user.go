package models

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
	"time"
	"errors"
)

type User struct {
	Base
	Name     string `json:"name"`
	Password string `json:"password"` // just for receive JSON plain-text password but not store in DB
	Secret   []byte
}

type Users []*User

func (u *User) Show(db *sql.DB) (*User, error) {
	err := db.QueryRow(
		"SELECT id, name, created_at, updated_at, deleted_at FROM user WHERE id = ?",
		u.ID,
	).Scan(
		&u.ID,
		&u.Name,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.DeletedAt,
	)
	if err != nil {
		log.Println("Error SELECT in user.Show:", err)
		return nil, err
	}
	// Filter only NOT Deleted User
	if u.DeletedAt.Valid == true {
		return nil, errors.New("User Deleted. - ผู้ใช้คนนี้ถูกลบแล้ว")
	}
	return u, nil
}

func (u *User) Index(db *sql.DB) ([]*User, error) {
	log.Println(">>> start AllUsers() >> db = ", db)
	rows, err := db.Query(
		"SELECT id, name, created_at, updated_at, deleted_at FROM user")
	if err != nil {
		log.Println(">>> db.Query Error= ", err)
		return nil, err
	}
	defer rows.Close()
	var users Users
	for rows.Next() {
		// We do not save plain text password to DB, just secret.
		var i = new(User)
		err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
		)
		if err != nil {
			log.Println(">>> rows.Scan() Error= ", err)
			return nil, err
		}
		// Filter only NOT Deleted User
		if i.DeletedAt.Valid == false {
			users = append(users, i)
		}
	}
	log.Println("return users", users)
	return users, nil
}

// Insert New User
func (u *User) Insert(db *sql.DB) (*User, error) {
	log.Println(">>start User.New() method")
	datetime := time.Now()
	datetime.Format(time.RFC3339)
	rs, err := db.Exec(
		"INSERT INTO user (name, secret, created_at) VALUES(?, ?, ?)",
		u.Name,
		u.Secret,
		datetime,
	) // no plain text u.Password save to DB
	if err != nil {
		log.Println(">>>Error cannot exec INSERT User: >>>", err)
		return nil, err
	}

	lastID, _ := rs.LastInsertId()
	num, _ := rs.RowsAffected()
	log.Printf("Last insert ID = %d, Number of rows = %d", lastID, num)

	// test query data
	n := new(User)
	err = db.QueryRow(
		"SELECT id, name, created_at FROM user WHERE id = ?",
		lastID,
	).Scan(
		&n.ID,
		&n.Name,
		&n.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Error not found: ", err)
		} else {
			log.Println(err)
		}
	}
	log.Println("Success insert record: ", n)
	return n, nil
}

// UpdateUser by id
func (u *User) Update(db *sql.DB) (*User, error) {
	log.Println(">>start models.user.Update() method")
	existUser := User{}
	err := db.QueryRow(
		"SELECT id, name FROM user WHERE id = ?",
		u.ID,
	).Scan(
		&existUser.ID,
		&existUser.Name,
	)
	if err != nil {
		log.Panic("Error db.QueryRow in user.Update()", err)
	}
	defer db.Close()
	log.Println("existUser: ", existUser)

	var rs sql.Result
	var updateTime = time.Now()
	updateTime.Format(time.RFC3339) // make Time Format fit to MariaDB.DateTime
	log.Println("Check: t := datetime: ", updateTime)
	if u.Password == "" { // Check if INPUT u.password is BLANK: So, user don't need to change password
		rs, err = db.Exec(
			"UPDATE user SET name= ?, updated_at=? WHERE id=?",
			u.Name,
			updateTime,
			existUser.ID,
		)
	} else {
		u.SetPass()
		rs, err = db.Exec(
			"UPDATE user SET name= ?, secret= ?, updated_at=? WHERE id =? ",
			u.Name,
			u.Secret,
			updateTime,
			existUser.ID,
		)
	}
	if err != nil {
		log.Panic("Error UPDATE user...", err)
		return nil, err
	}
	// db.QueryRow to check if correct update record
	countRow, _ := rs.RowsAffected()
	log.Println("Number of row updated: ", countRow)
	n := User{}
	err = db.QueryRow(
		"SELECT id, name, secret, created_at, updated_at FROM user WHERE id =?",
		existUser.ID,
	).Scan(
		&n.ID,
		&n.Name,
		&n.Secret,
		&n.CreatedAt,
		&n.UpdatedAt,
	)
	if err != nil {
		log.Println("Error when SELECT updated row??? >>>", err)
	}
	// return new query update record
	return &n, nil
}

func (u *User) SetPass() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return err
	}
	u.Secret = hash
	log.Println("Got u.Secret: ", u.Secret)
	return nil
}


func (u *User) VerifyPass(p string) error { // not export call from Add() or Update
	err := bcrypt.CompareHashAndPassword(u.Secret, []byte(p))
	if err != nil {
		return err
	}
	return nil
}

func (u *User) FindByName(db *sql.DB) error{
	err := db.QueryRow(
		"SELECT id, name, secret FROM user WHERE name = ?",
		u.Name,
	).Scan(
		&u.ID,
		&u.Name,
		&u.Secret,
	)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
// Method models.User.Del to delete User (Later we will implement my framework just add delete DateX
func (u *User) Delete(db *sql.DB) error {
	delTime := time.Now()
	delTime.Format(time.RFC3339)
	sql := "UPDATE user SET deleted_at = ? WHERE id = ?"
	rs, err := db.Exec(sql, delTime, u.ID)
	if err != nil {
		log.Println(err)
		return err
	}
	rowCnt, err := rs.RowsAffected()
	if err != nil {
		log.Println(err)
	}
	log.Println("Deleted:", rowCnt, "row(s).")

	// TODO return Deleted User
	err = db.QueryRow(
		"SELECT id, name, deleted_at FROM user WHERE id =?",
		u.ID,
	).Scan(
		&u.ID,
		&u.Name,
		&u.DeletedAt,
	)
	if err != nil {
		log.Println("Error when SELECT updated row??? >>>", err)
		return err
	}
	return nil
}

func (u *User) Undelete(db *sql.DB) error {
	sql := "UPDATE user SET deleted_at = ? WHERE id = ?"
	rs, err := db.Exec(sql, nil, u.ID)
	if err != nil {
		log.Println(err)
		return err
	}
	rowCnt, err := rs.RowsAffected()
	if err != nil {
		log.Println(err)
	}
	log.Println("Undeleted:", rowCnt, "row(s).")

	// TODO return Deleted User
	err = db.QueryRow(
		"SELECT id, name, deleted_at FROM user WHERE id =?",
		u.ID,
	).Scan(
		&u.ID,
		&u.Name,
		&u.DeletedAt,
	)
	if err != nil {
		log.Println("Error when SELECT updated row??? >>>", err)
	}
	return nil
}

// function models.User.SearchUsers() here!
func SearchUsers(db *sql.DB, s string) (Users, error) {
	s = "%" + strings.ToLower(s) + "%"
	stmt, err := db.Prepare("SELECT id, name FROM user WHERE LOWER(name) LIKE ?")
	if err != nil {
		log.Println("Error in SearchUsers() - db.Prepare() >>>", err)
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(s)
	if err != nil {
		log.Println("Error in SearchUsers - stmt.Query() >>>", err)
		return nil, err
	}
	defer rows.Close()

	users := Users{}
	for rows.Next() {
		u := new(User)
		err := rows.Scan(&u.ID, &u.Name)
		if err != nil {
			log.Println(">>> rows.Scan() Error= ", err)
			return nil, err
		}
		users = append(users, u)
	}
	log.Println("users = ", users)
	return users, nil
}