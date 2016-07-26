package models

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log"
	"time"
	"github.com/jmoiron/sqlx"
)

type Person struct {
	Base
	First     string    `json:"first"`
	Last      string    `json:"last"`
	Nick      sql.NullString    `json:"nick"`
	Sex       string    `json:"sex"`
	BirthDate time.Time `json:"birth_date"`
}

type Job struct {
	Base
	Name     string
	ParentID uint64
}

type Org struct {
	Base
	Name     string
	ParentID uint64
}

func (p *Person) New(db *sqlx.DB) error {
	log.Println("run models.Person.New method from:", p)
	birthDate := p.BirthDate.Format(time.RFC3339)
	sql := `INSERT INTO person (first, last, nick, sex, birth_date)
		VALUES (?,?,?,?,?)`
	res, err := db.Exec(sql,
		p.First,
		p.Last,
		p.Nick,
		p.Sex,
		birthDate,
	)
	if err != nil {
		log.Println("Error insert into person...:", err)
		return err
	}
	id, _ := res.LastInsertId()
	var date mysql.NullTime
	sql = `SELECT id, created, first, last, nick, sex, birth_date FROM person WHERE id = ?`
	err = db.QueryRow(sql, id).Scan(
		&p.ID,
		&p.Created,
		&p.First,
		&p.Last,
		&p.Nick,
		&p.Sex,
		&date,
	)
	if err != nil {
		log.Println(err)
		return err
	}
	if date.Valid {
		p.BirthDate = date.Time
	}
	return nil
}

func (p *Person) All(db *sqlx.DB) ([]*Person, error) {
	log.Println("run models.Person.All method from:", p)
	sql := `SELECT id, created, updated, deleted,
		first, last, nick, sex, birth_date
		FROM person`
	rows, err := db.Query(sql)
	if err != nil {
		log.Println(">>> db.Query Error= ", err)
		return nil, err
	}
	defer rows.Close()
	var (
		persons                     []*Person
		//nick                        sql.NullString
		updated, deleted, birthDate mysql.NullTime
	)

	for rows.Next() {
		err := rows.Scan(
			&p.ID,
			&p.Created,
			&updated,
			&deleted,
			&p.First,
			&p.Last,
			&p.Nick,
			&p.Sex,
			&birthDate,
		)
		if updated.Valid {
			p.Updated = updated.Time
		}
		if deleted.Valid {
			p.Deleted = deleted.Time
		}
		//if nick.Valid {
		//	p.Nick = nick.String
		//}
		if birthDate.Valid {
			p.BirthDate = birthDate.Time
		}
		if err != nil {
			log.Println("Error in rows.Scan():", err)
		}
		persons = append(persons, p)
	}
	return persons, nil
}
