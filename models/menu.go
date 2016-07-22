package models

import (
	"database/sql"
	"log"
)

type Menu struct {
	ID       int    `json:"id"`
	ParentID int    `json:"parent_id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	Path     string `json:"path"`
	Note     string `json:"note"`
}

type Menus []*Menu


func (m *Menu) All(db *sql.DB) ([]*Menu, error) {
	rows, err := db.Query("SELECT id, parent_id, name, code, path, note FROM menu")
	if err != nil {
		log.Println(">>> 1. db.Query Error= ", err)
		return nil, err
	}
	defer rows.Close()

	var menus Menus
	for rows.Next() {
		m := new(Menu)
		err := rows.Scan(
			&m.ID,
			&m.ParentID,
			&m.Name,
			&m.Code,
			&m.Path,
			&m.Note,
		)
		if err != nil {
			log.Println(">>> rows.Scan() Error= ", err)
			return nil, err
		}
		menus = append(menus, m)
	}
	log.Println("Menu:", menus)
	return menus, nil
}

func (m *Menu) Insert(db *sql.DB) error {
	log.Println("Start m.New()")
	sql := "INSERT INTO menu (parent_id, name, code, path, note ) VALUES(?,?,?,?,?)"

	rs, err := db.Exec(sql,
		m.ParentID,
		m.Name,
		m.Code,
		m.Path,
		m.Note,
	)
	if err != nil {
		log.Println(">>>Error cannot exec INSERT menu: >>>", err)
		return err
	}
	log.Println(rs)
	lastID, _ := rs.LastInsertId()

	err = db.QueryRow(
		"SELECT id, parent_id, name, code, path, note FROM menu WHERE id = ?",
		lastID,
	).Scan(
		&m.ID,
		&m.ParentID,
		&m.Name,
		&m.Code,
		&m.Path,
		&m.Note,
	)
	if err != nil {
		return err
	}
	log.Println("Success insert record:", m)
	return nil
}
