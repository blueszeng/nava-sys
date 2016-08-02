package model

import (
	"github.com/jmoiron/sqlx"
	"log"
)

type Menu struct {
	ID           uint64 `json:"id"`
	ParentID     uint64 `json:"parent_id" db:"parent_id"`
	Text         string `json:"text"`
	Icon         string `json:"icon"`
	SelectedIcon string `json:"selectedIcon" db:"selected_icon"`
	Href         string `json:"href"`
	Path         string `json:"path"`
	Note         string `json:"note"`
}

type Menus []*Menu

type Role struct {
	Base
	TH string
	EN string
}

type UserRole struct {
	UserID uint64
	RoleID uint64
}

type RoleMenu struct {
	RoleID   uint64
	MenuID   uint64
	CanRead  bool
	CanWrite bool
}

func (m *Menu) All(db *sqlx.DB) ([]*Menu, error) {
	var menus Menus
	err := db.Select(&menus, `SELECT * FROM menu`)
	if err != nil {
		log.Println(">>> 1. db.Query Error= ", err)
		return nil, err
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
	menu := new(Menu)
	err = db.Get(&menu, sql, lastID)
	if err != nil {
		return err
	}
	log.Println("Success insert record:", menu)
	return nil
}

func (u *User) FindMenuByUser(db *sqlx.DB) ([]*Menu, error) {
	s := `
	SELECT
		menu.*
	FROM user
	LEFT JOIN user_role ON user.id = user_role.user_id
	LEFT JOIN role ON user_role.role_id = role.id
	LEFT JOIN role_menu ON role.id = role_menu.role_id
	LEFT JOIN menu ON role_menu.menu_id = menu.id
	WHERE user.id = ?
	`
	var menus []*Menu
	err := db.Select(&menus, s, u.ID)
	if err != nil {
		log.Fatal("Error in db.Select(): ", err)
	}
	return menus, nil
}

//func removeDuplicates(a []uint64) []uint64 {
//	result := []uint64{}
//	seen := map[uint64]uint64{}
//	for _, val := range a {
//		if _, ok := seen[val]; !ok {
//			result = append(result, val)
//			seen[val] = val
//		}
//	}
//	return result
//}