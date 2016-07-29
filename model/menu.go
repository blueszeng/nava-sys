package model

import (
	"log"
	"github.com/jmoiron/sqlx"
)

type Menu struct {
	ID       int    `json:"id"`
	ParentID int    `json:"parent_id"`
	Text     string `json:"name"`
	Icon     string `json:"icon"`
	SelectedIcon string `json:"selectedIcon" db:"selected_icon"`
	Href string `json:"href" db:"href"`
	Path     string `json:"path"`
	Note     string `json:"note"`
}

type Menus []*Menu


func (m *Menu) All(db *sqlx.DB) ([]*Menu, error) {
	rows, err := db.Query(`
	SELECT
		id,
		parent_id,
		text,
		icon,
		selected_icon,
		href,
		path,
		note
	FROM menu`)
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
			&m.Text,
			&m.Icon,
			&m.SelectedIcon,
			&m.Href,
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

func (m *Menu) Insert(db *sqlx.DB) error {
	log.Println("Start m.New()")
	sql := `INSERT INTO menu (SELECT
		parent_id,
		text,
		icon,
		selected_icon,
		href,
		path,
		note
	VALUES(?,?,?,?,?,?)`

	rs, err := db.Exec(sql,
		m.ParentID,
		m.Text,
		m.Icon,
		m.SelectedIcon,
		m.Href,
		m.Path,
		m.Note,
	)
	if err != nil {
		log.Println(">>>Error cannot exec INSERT menu: >>>", err)
		return err
	}
	log.Println(rs)
	lastID, _ := rs.LastInsertId()
	sql = `SELECT
		id,
		parent_id,
		text,
		icon,
		selected_icon,
		href,
		path,
		note
	FROM menu WHERE id = ?`
	err = db.QueryRow(sql, lastID).Scan(
		&m.ID,
		&m.ParentID,
		&m.Text,
		&m.Icon,
		&m.SelectedIcon,
		&m.Href,
		&m.Path,
		&m.Note,
	)
	if err != nil {
		return err
	}
	log.Println("Success insert record:", m)
	return nil
}
