package model

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
	"time"
)

type User struct {
	Base
	Name     string `json:"name"`
	Password string `json:"password,omitempty"` // just for receive JSON plain-text password but not store in DB
	Secret   []byte `json:"-"`
}

type Users []*User

func (u *User) Get(db *sqlx.DB) (User, error) {
	sql := `
	SELECT
		id,
		name,
		created,
		updated,
		deleted
	FROM user
	WHERE id = ?`
	var user User
	err := db.Get(&user, sql, u.ID)
	if err != nil {
		log.Println("Error SELECT() in User.Show:", err)
		return user, err
	}
	// Filter only NOT Deleted User
	//if deleted.Valid == true {
	if user.Deleted.Valid == true {
		user = User{}
		return user, errors.New("User Deleted. - ผู้ใช้คนนี้ถูกลบแล้ว")
	}
	return user, nil
}

func (u *User) All(db *sqlx.DB) ([]*User, error) {
	log.Println(">>> start AllUsers() >> db = ", db)
	err := db.Ping()
	if err != nil {
		log.Println("Ping Error", err)
	}

	var users Users
	sql := `SELECT id, name, created, updated, deleted FROM user`
	rows, err := db.Queryx(sql)
	if err != nil {
		log.Println(">>> db.Query Error= ", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		// We do not save plain text password to DB, just secret.
		var i = new(User)
		err := rows.StructScan(&i)
		if err != nil {
			log.Println(">>> rows.Scan() Error= ", err)
			return nil, err
		}
		// Filter Deleted User
		if i.Deleted.Valid == false {
			users = append(users, i)
		}
	}
	log.Println("return users", users)
	return users, nil
}

// Insert New User
func (u *User) New(db *sqlx.DB) (*User, error) {
	log.Println(">>start User.New() method")
	log.Println("Test User receiver:", u.Name)
	rsp, err := db.Exec(
		"INSERT INTO user (name, secret) VALUES(?, ?)",
		u.Name,
		u.Secret,
	) // no plain text u.Password save to DB
	if err != nil {
		log.Println(">>>Error cannot exec INSERT User: >>>", err)
		return nil, err
	}
	lastID, err := rsp.LastInsertId()
	if err != nil {
		log.Println("Error in rsp.LastInsertId", err)
	}
	log.Println("LastInsertId=", lastID)
	// test query data
	newUser := User{}
	err = db.Get(&newUser, "SELECT * FROM user WHERE id = ?", lastID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Error not found: ", err)
		} else {
			log.Println(err)
		}
		return nil, err
	}
	log.Println("Success insert record: ", newUser)
	return &newUser, nil
}

// UpdateUser by id
func (u *User) Update(db *sqlx.DB) (*User, error) {
	log.Println(">>start models.user.Update() method")

	existUser := User{}
	s := `SELECT *
		FROM user
		WHERE id = ?`
	err := db.Get(&existUser, s, u.ID)
	if err != nil {
		log.Println("Error db.QueryRow in user.Update()", err)
		return nil, err
	}
	if existUser.Deleted.Valid == true {
		return nil, errors.New("User Deleted")
	}
	log.Println("existUser: ", existUser)

	if u.Password == "" { // Check if INPUT u.password is BLANK: So, user don't need to change password
		s = `UPDATE user SET
				name= ?
			WHERE id=?`
		_, err = db.Exec(s,
			u.Name,
			existUser.ID,
		)
	} else {
		u.SetPass()
		s = `UPDATE user SET
				name= ?,
				secret= ?
			WHERE id =?`
		_, err = db.Exec(s,
			u.Name,
			u.Secret,
			existUser.ID,
		)
	}
	if err != nil {
		log.Println("Error UPDATE user...", err)
		return nil, err
	}

	// db.QueryRow to check if correct update record
	n := User{}
	s = `SELECT *
		FROM user
		WHERE id =?`
	err = db.Get(&n, s, existUser.ID)
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
	log.Println("bcrypt...", u.Secret, p)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) FindByName(db *sqlx.DB) (*User, error) {
	sql := `
		SELECT *
		FROM user
		WHERE name = ?`
	var user User
	err := db.Get(&user, sql, u.Name)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &user, nil
}

// Method models.User.Del to delete User (Later we will implement my framework just add delete DateX
func (u *User) Del(db *sqlx.DB) (*User, error) {
	now := time.Now()
	now.Format(time.RFC3339)
	sql := "UPDATE user SET deleted = ? WHERE id = ?"
	_, err := db.Exec(sql, now, u.ID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var user User
	sql = `SELECT * FROM user WHERE id =?`
	err = db.Get(&user, sql, u.ID)
	if err != nil {
		log.Println("Error when SELECT updated row??? >>>", err)
		return nil, err
	}
	user.Secret = nil
	return &user, nil
}

func (u *User) Undel(db *sqlx.DB) (*User, error) {
	sql := "UPDATE user SET deleted = ? WHERE id = ?"
	rs, err := db.Exec(sql, nil, u.ID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("u.ID=", u.ID)
	rowCnt, err := rs.RowsAffected()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("Undeleted:", rowCnt, "row(s).")
	var user User
	sql = `SELECT * FROM user WHERE id =?`
	err = db.Get(&user, sql, u.ID)
	if err != nil {
		log.Println("Error when SELECT updated row??? >>>", err)
	}
	user.Secret = nil
	return &user, nil
}

func (u *User) Menus(db *sqlx.DB) ([]*Menu, error) {
	log.Println("call model.User.Menus")
	s := `
	SELECT
		menu.*
	FROM user
	LEFT JOIN user_role ON user.id = user_role.user_id
	LEFT JOIN role ON user_role.role_id = role.id
	LEFT JOIN role_menu ON role.id = role_menu.role_id
	LEFT JOIN menu ON role_menu.menu_id = menu.id
	WHERE user.id = ? AND role_menu.can_read = true
	`
	var menus []*Menu
	err := db.Select(&menus, s, u.ID)
	if err != nil {
		log.Fatal("Error in db.Select(): ", err)
	}
	return menus, nil
}

// function models.User.SearchUsers() here!
func SearchUsers(db *sqlx.DB, s string) (Users, error) {
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
