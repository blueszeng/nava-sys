package model

import (
	"github.com/jmoiron/sqlx"
	"log"
)

type Menu struct {
	Base
	ParentID     uint64  `json:"parent_id" db:"parent_id"`
	Text         string  `json:"text"`
	Icon         string  `json:"icon"`
	SelectedIcon string  `json:"selectedIcon" db:"selected_icon"`
	Href         string  `json:"href"`
	Path         string  `json:"-"`
	Note         string  `json:"-"`
	Child        []*Menu `json:"nodes,omitempty"`
}

type Menus []*Menu

func (m *Menu) All(db *sqlx.DB) ([]*Menu, error) {
	var menus Menus
	err := db.Select(&menus, `SELECT * FROM menu WHERE deleted IS NULL`)
	if err != nil {
		log.Println(">>> 1. db.Query Error= ", err)
		return nil, err
	}
	log.Println("Menu:", menus)
	return menus, nil
}

func (m *Menu) Insert(db *sqlx.DB) (*Menu, error) {
	log.Println("Start m.Insert()")
	sql := `INSERT INTO menu (
		parent_id,
		text,
		icon,
		selected_icon,
		href,
		path,
		note
	)
	VALUES(?,?,?,?,?,?,?)`

	res, err := db.Exec(sql,
		m.ParentID,
		m.Text,
		m.Icon,
		m.SelectedIcon,
		m.Href,
		m.Path,
		m.Note,
	)
	if err != nil {
		log.Println(">>>Error exec INSERT menu: >>>", err)
		return nil, err
	}
	log.Println(res)
	var insertedMenu Menu
	sql = `SELECT * FROM menu WHERE id = ? AND deleted IS NULL`
	lastID, _ := res.LastInsertId()
	err = db.Get(&insertedMenu, sql, lastID)
	if err != nil {
		return nil, err
	}
	log.Println("Success insert record:", insertedMenu)
	return &insertedMenu, nil
}

func (m *Menu) Get(db *sqlx.DB) (*Menu, error) {
	log.Println("run model.Menu.Get()")
	sql := `SELECT * FROM menu WHERE id = ?`
	err := db.Get(&m, sql, m.ID)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (this *Menu) Size() int {
	var size int = len(this.Child)
	for _, c := range this.Child {
		size += c.Size()
	}
	return size
}

func (this *Menu) Add(menus ...*Menu) bool {
	var size = this.Size()
	for _, node := range menus {
		if node.ParentID == this.ID {
			this.Child = append(this.Child, node)
		} else {
			for _, child := range this.Child {
				if child.Add(node) {
					break
				}
			}
		}
	}
	return this.Size() == size+len(menus)
}

//func (menus *Menus) Tree() *Menu {
//	tree := new(Menu)
//	for _, m := range menus{
//		tree.Add(m)
//	}
//	return tree
//}

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
