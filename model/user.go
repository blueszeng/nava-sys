package model

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
	"time"
	"github.com/jmoiron/sqlx"
)

type User struct {
	Base
	Name     string `json:"name"`
	Password string `json:"password"` // just for receive JSON plain-text password but not store in DB
	Secret   []byte
}

type Users []*User

func (u *User) Show(db *sqlx.DB) (User, error) {
	//var updated, deleted mysql.NullTime
	var user User
	err := db.QueryRowx(
		"SELECT id, name, created, updated, deleted FROM user WHERE id = ?",
		u.ID,
	).StructScan(&user)
	//).Scan(
	//	&u.ID,
	//	&u.Name,
	//	&u.Created,
	//	&updated,
	//	&deleted,
	//)
	//if updated.Valid {
	//	u.Updated = updated.Time
	//}
	//if deleted.Valid {
	//	u.Deleted = deleted.Time
	//}
	if err != nil {
		log.Println("Error SELECT in user.Show:", err)
		return user, err
	}
	// Filter only NOT Deleted User
	//if deleted.Valid == true {
	if u.Deleted.Valid == true {
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
		//var i = new(User)
		err := rows.StructScan(&u)
		if err != nil {
			log.Println(">>> rows.Scan() Error= ", err)
			return nil, err
		}
		// Filter Deleted User
		if u.Deleted.Valid == false {
			users = append(users, u)
		}
	}
	log.Println("return users", users)
	return users, nil
}

// Insert New User
func (u *User) New(db *sqlx.DB) (*User, error) {
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
		"SELECT id, name, created FROM user WHERE id = ?",
		lastID,
	).Scan(
		&n.ID,
		&n.Name,
		&n.Created,
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
func (u *User) Update(db *sqlx.DB) (*User, error) {
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
		log.Println("Error db.QueryRow in user.Update()", err)
		return nil, err
	}
	//defer db.Close()
	log.Println("existUser: ", existUser)

	var rs sql.Result
	now := time.Now()
	now.Format(time.RFC3339) // make Time Format fit to MariaDB.DateTime
	//log.Println("Check: t := datetime: ", updateTime)
	if u.Password == "" { // Check if INPUT u.password is BLANK: So, user don't need to change password
		rs, err = db.Exec(
			"UPDATE user SET name= ?, updated=? WHERE id=?",
			u.Name,
			now,
			existUser.ID,
		)
	} else {
		u.SetPass()
		rs, err = db.Exec(
			"UPDATE user SET name= ?, secret= ?, updated=? WHERE id =? ",
			u.Name,
			u.Secret,
			now,
			existUser.ID,
		)
	}
	if err != nil {
		log.Println("Error UPDATE user...", err)
		return nil, err
	}

	// db.QueryRow to check if correct update record
	countRow, _ := rs.RowsAffected()
	log.Println("Number of row updated: ", countRow)

	var updated mysql.NullTime
	n := User{}
	err = db.QueryRowx(
		"SELECT id, name, secret, created, updated FROM user WHERE id =?",
		existUser.ID,
	).Scan(
		&n.ID,
		&n.Name,
		&n.Secret,
		&n.Created,
		&updated,
	)
	//if updated.Valid {
	//	n.Updated = updated.Time
	//}
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

func (u *User) FindByName(db *sqlx.DB) error {
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
func (u *User) Del(db *sqlx.DB) error {
	now := time.Now()
	now.Format(time.RFC3339)
	sql := "UPDATE user SET deleted = ? WHERE id = ?"
	rs, err := db.Exec(sql, now, u.ID)
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
	var deleted mysql.NullTime
	err = db.QueryRow(
		"SELECT id, name, deleted FROM user WHERE id =?",
		u.ID,
	).Scan(
		&u.ID,
		&u.Name,
		&deleted,
	)
	if err != nil {
		log.Println("Error when SELECT updated row??? >>>", err)
		return err
	}
	//if deleted.Valid {
	//	u.Deleted = deleted.Time
	//}
	return nil
}

func (u *User) Undel(db *sqlx.DB) error {
	sql := "UPDATE user SET deleted = ? WHERE id = ?"
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
	//var deleted mysql.NullTime
	err = db.QueryRow(
		"SELECT id, name, deleted FROM user WHERE id =?",
		u.ID,
	).Scan(
		&u.ID,
		&u.Name,
		&u.Deleted,
	)
	if err != nil {
		log.Println("Error when SELECT updated row??? >>>", err)
	}
	//if deleted.Valid {
	//	u.Deleted = deleted.Time
	//}
	return nil
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
