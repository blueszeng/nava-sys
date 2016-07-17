package models

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"log"
	"errors"
	"strings"
)

type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"` // just for receive JSON plain-text password but not store in DB
	Secret   []byte
	//PeopleID sql.NullInt64 `json:people_id` // TODO: ยังไม่รู้จะรับค่า JSON decode มาใส่ sql.NullXYZ ยังไง เ
	// TODO: เนื่องจาก sql.NullXYZ เป็น struct {???: ???, Valid: boolean}
}

type Users []*User

func (u *User) Show(db *sql.DB) (*User, error) {
	err := db.QueryRow(
		"SELECT * FROM user WHERE id = ?",
		u.ID,
	).Scan(
		&u.ID,
		&u.Name,
		&u.Secret,
	)
	if err != nil {
		log.Fatal("Error SELECT * in user.Show:", err)
		return nil, err
	}
	return u, nil
}

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
		// We do not save plain text password to DB, just secret.
		var i = new(User)
		err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Secret,
		)
		if err != nil {
			log.Println(">>> rows.Scan() Error= ", err)
			return nil, err
		}
		users = append(users, i)
		log.Println("users= ",users, "u= ", i)
	}
	log.Println("return users", users)
	return users, nil
}

// Insert New User
func (u *User) New(db *sql.DB) (*User, error) {
	log.Println(">>start User.New() method")
	rs, err := db.Exec(
		"INSERT INTO user (name, secret) VALUES(?, ?)",
		u.Name,
		u.Secret,
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
		"SELECT id, name, secret FROM user WHERE id = ?",
		lastID,
	).Scan(
		&n.ID,
		&n.Name,
		//&n.Secret,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Fatal("Error not found: ", err)
		} else {
			log.Fatal(err)
		}
	}
	log.Println("Success insert record: ", n)
	return n, nil
}

// UpdateUser by id
func (u *User) Update(db *sql.DB) (*User, error) {
	log.Println(">>start models.user.Update() method")
	// TODO: Check if no exist user.Name // return error and ask to create new user.
	existUser := User{}
	err := db.QueryRow(
		"SELECT id, name, secret FROM user WHERE id = ?",
		u.ID,
	).Scan(
		&existUser.ID,
		&existUser.Name,
		&existUser.Secret,
	)
	if err != nil {
		log.Panic("Error db.QueryRow in user.Update()", err)
	}
	defer db.Close()
	log.Println("existUser: ", existUser)
	if u.Name != existUser.Name {
		err = errors.New("No match exist name: Do you want to create NEW User?")
	}
	// Ok match exist user name...Run command to update data
	var rs sql.Result
	if u.Password == "" { // Check if u.password is BLANK: So, user don't need to change password
		rs, err = db.Exec(
			"UPDATE user SET name= ? WHERE id = ?",
			u.Name,
			existUser.ID,
		)
	} else {
		u.SetPass()
		rs, err = db.Exec(
			"UPDATE user SET name= ?, secret= ? WHERE id =? ",
			u.Name,
			u.Secret,
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
		"SELECT id, name, secret FROM user WHERE id =?",
		existUser.ID,
	).Scan(
		&n.ID,
		&n.Name,
		&n.Secret,
	)
	if err != nil {
		log.Println("Error when SELECT updated row??? >>>", err)
	}
	// return new query update record
	return &n, nil
}

// TODO: Method models.User.Del to delete User (Later we will implement my framework just add delete DateX
func (u *User) Delete(db *sql.DB) error {
	//TODO: Code to DELETE here.
	return nil
}

func (u *User) SetPass() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalln(err)
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

func (u *User) SearchByName(db *sql.DB) error{
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