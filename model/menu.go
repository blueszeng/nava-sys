package model

import (
	"github.com/jmoiron/sqlx"
	"log"
)

type Menu struct {
	ID           uint64  `json:"id"`
	ParentID     uint64  `json:"parent_id" db:"parent_id"`
	Text         string  `json:"text"`
	Icon         string  `json:"icon"`
	SelectedIcon string  `json:"selectedIcon" db:"selected_icon"`
	Href         string  `json:"href"`
	Path         string  `json:"path"`
	Note         string  `json:"note"`
	Child        []*Menu `json:"nodes,omitempty"`
}

type Menus []*Menu

type Role struct {
	ID uint64 `json:"id"`
	TH string `json:"th"`
	EN string `json:"en"`
}

type UserRole struct {
	ID     uint64 `json:"id"`
	UserID uint64 `json:"user_id" db:"user_id"`
	RoleID uint64 `json:"role_id" db:"role_id"`
}

type RoleMenu struct {
	RoleID     uint64 `json:"role_id" db:"role_id"`
	MenuID     uint64 `json:"menu_id" db:"menu_id"`
	CanRead    bool   `json:"can_read" db:"can_read"`
	CanWrite   bool   `json:"can_write" db:"can_write"`
	CanDelete  bool   `json:"can_delete" db:"can_delete"`
	CanApprove bool   `json:"can_approve" db:"can_approve"`
	CanCancel  bool   `json:"can_cancel" db:"can_cancel"`
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

func (m *Menu) New(db *sqlx.DB) error {
	log.Println("Start m.New()")
	sql := `INSERT INTO menu (
		parent_id,
		text,
		icon,
		selected_icon,
		href,
		path,
		note
	)
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
		log.Println(">>>Error exec INSERT menu: >>>", err)
		return err
	}
	log.Println(rs)
	menu := new(Menu)
	sql = `SELECT * FROM menu WHERE id = ?`
	lastID, _ := rs.LastInsertId()
	err = db.Get(&menu, sql, lastID)
	if err != nil {
		return err
	}
	log.Println("Success insert record:", menu)
	return nil
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
