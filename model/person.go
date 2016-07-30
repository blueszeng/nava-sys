package model

import (
	"github.com/go-sql-driver/mysql"
	"log"
	"time"
	"github.com/jmoiron/sqlx"
)

type Person struct {
	Base
	First     JsonNullString    `json:"first"`
	Last      JsonNullString    `json:"last"`
	Nick      JsonNullString    `json:"nick"`
	Sex       JsonNullString    `json:"sex"`
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
	sql := `INSERT INTO person (
		first,
		last,
		nick,
		sex,
		birth_date
	)
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
	//sql := `SELECT id, created, updated, deleted,
	//	first, last, nick, sex, birth_date
	//	FROM person`
	sql := `SELECT * FROM person`
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
	person := new(Person)
	for rows.Next() {
		err := rows.Scan(
			&person.ID,
			&person.Created,
			&updated,
			&deleted,
			&person.First,
			&person.Last,
			&person.Nick,
			&person.Sex,
			&birthDate,
		)
		if updated.Valid {
			person.Updated = updated.Time
		}
		if deleted.Valid {
			person.Deleted = deleted.Time
		}
		//if nick.Valid {
		//	p.Nick = nick.String
		//}
		if birthDate.Valid {
			person.BirthDate = birthDate.Time
		}
		if err != nil {
			log.Println("Error in rows.Scan():", err)
		}
		persons = append(persons, person)
		log.Println(persons)
	}
	return persons, nil
}
