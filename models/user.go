package models

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"log"
	"errors"
)

type User struct {
	ID       int64  `json:id`
	Name     string `json:name`
	Password string `json:password` // just for receive JSON plain-text password but not store in DB
	Secret   []byte
	//PeopleID sql.NullInt64 `json:people_id` // TODO: ยังไม่รู้จะรับค่า JSON decode มาใส่ sql.NullXYZ ยังไง เ
	// TODO: เนื่องจาก sql.NullXYZ เป็น struct {???: ???, Valid: boolean}
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
		// we not save plain text password in database just secret
		err := rows.Scan(&u.ID, &u.Name, &u.Secret)
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

// Insert New User
func (u *User) New(db *sql.DB) (*User, error) {
	log.Println(">>start User.New() method")

	// TODO: check if exist same user.Name
	// TODO: IF user has permission, ask to edit or ...reject and return error

	res, err := db.Exec("INSERT INTO user VALUE(?, ?, ?)",u.ID, u.Name, u.Secret) // no plain text u.Password save to DB
	if err != nil {
		log.Fatal(">>>Error Exec INSERT...User: >>>", err)
	}

	lastID, _ := res.LastInsertId()
	num, _ := res.RowsAffected()
	log.Printf("Last insert ID = %d, Number of rows = %d", lastID, num)

	// test query data
	n := new(User)
	err = db.QueryRow("SELECT id, name, secret FROM user WHERE id = ?", lastID).Scan(&n.ID, &n.Name, &n.Secret)
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

// Edit/UpdateUser
func (u *User) Update(db *sql.DB) (*User, error) {
	log.Println(">>start u.Update() method")

	// TODO: Check if no exist user.Name // return error and ask to create new user.
	existUser := User{}
	err := db.QueryRow("SELECT name FROM user WHERE name = ?", u.Name).Scan(&existUser.Name)
	if err != nil {
		log.Panic("Error db.QueryRow in user.Update()", err)
	}
	if u.Name == existUser.Name {
		// Ok match exist user name...Run command to update data
		var res sql.Result
		var err error
		// TODO: Check if u.password == nil: then user don't need to change password
		// db.Exec to update record
		if u.Password == "" {
			res, err = db.Exec("UPDATE user SET name= ? WHERE id = ?", u.Name, existUser.ID)
		} else {
			res, err = db.Exec("UPDATE user SET name= ?, secret= ? WHERE id =? ", u.Name, u.Secret, existUser.ID)
		}
		if err != nil {
			log.Panic("Error UPDATE user...", err)
			return nil, err
		}
		// db.QueryRow to check if correct update record
		countrow, _ := res.RowsAffected()
		log.Println("Number of row updated: ", countrow)
		n := User{}
		err = db.QueryRow("SELECT name, secret FROM user WHERE id =?", existUser.ID).Scan(&n.Name, &n.Secret)
		if err != nil {
			log.Println("Error when SELECT updated row??? >>>", err)
		}
		// return new query update record
		return &n, nil
	}
	err = errors.New("No match exist name: Do you want to create NEW User?")


	return u, nil
}

// Delete User (Later we will implement my framework just add delete DateX
func (u *User) Delete(db *sql.DB) error {
	return nil
}

func (u *User) SetPass() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	u.Secret = hash
	return nil
}

func (u *User) verifyPass(p string) error { // not export call from Add() or Update
	err := bcrypt.CompareHashAndPassword(u.Secret, []byte(p))
	if err != nil {
		return err
	}
	return nil
}
