package models

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"log"
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
	i := new(User)
	err = db.QueryRow("SELECT ID, name FROM user WHERE id = ?", lastID).Scan(&i.ID, &i.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Fatal("Error not found: ", err)
		} else {
			log.Fatal(err)
		}
	}
	log.Println("Success insert record: ", i)
	return i, nil
}

// Edit/UpdateUser
func (u *User) Update(db *sql.DB) (*User, error) {
	log.Println(">>start u.Update() method")

	// TODO: Check if no match user.Name // return error and ask to create new user.

	// TODO: Check if hash password not match maybe user need to change password

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
