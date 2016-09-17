package model

import (
	"github.com/jmoiron/sqlx"
	"log"
)

type Role struct {
	Base
	NameTh string `json:"name_th" db:"name_th"`
	NameEn string `json:"name_en" db:"name_en"`
	OrgId  uint64 `json:"org_id" db:"org_id"`
}

type UserRole struct {
	Base
	UserID uint64 `json:"user_id" db:"user_id"`
	RoleID uint64 `json:"role_id" db:"role_id"`
}

type RoleMenu struct {
	Base
	RoleID     uint64 `json:"role_id" db:"role_id"`
	MenuID     uint64 `json:"menu_id" db:"menu_id"`
	CanRead    bool   `json:"can_read" db:"can_read"`
	CanWrite   bool   `json:"can_write" db:"can_write"`
	CanDelete  bool   `json:"can_delete" db:"can_delete"`
	CanRun     bool   `json:"can_run" db:"can_run"`
	CanApprove bool   `json:"can_approve" db:"can_approve"`
	CanCancel  bool   `json:"can_cancel" db:"can_cancel"`
}

func (r *Role) All(db *sqlx.DB) ([]*Role, error) {
	log.Println("call AllRole")

	var roles []*Role
	sql := `SELECT * FROM role WHERE deleted IS NULL`
	err := db.Select(&roles, sql)
	if err != nil {
		log.Println("Error in model.Role.All", err)
		return nil, err
	}
	return roles, nil
}

// -----------------------
// User Permission struct
// -----------------------
type UserPermission struct {
	UserID     uint64 `json:"user_id"`
	UserName   string `json:"user_name"`
	Permission []*MenuPermission
}
type MenuPermission struct {
	OrgId      uint64 `json:"org_id" db:"org_id"`
	OrgNameTh  string `json:"org_name_th" db:"org_name_th"`
	ID         uint64 `json:"id"`
	Text       string `json:"text" db:"text"`
	ParentID   uint64 `json:"parent_id" db:"parent_id"`
	CanRead    bool   `json:"can_read" db:"can_read"`
	CanWrite   bool   `json:"can_write" db:"can_write"`
	CanDelete  bool   `json:"can_delete" db:"can_delete"`
	CanRun     bool   `json:"can_run" db:"can_run"`
	CanApprove bool   `json:"can_approve" db:"can_approve"`
	CanCancel  bool   `json:"can_cancel" db:"can_cancel"`
}

func (u *User) Permission(db *sqlx.DB) (UserPermission, error) {
	sql := `
	SELECT
		org.id as org_id,
		org.name_th as org_name_th,
		menu.id,
		menu.text,
		menu.parent_id,
		role_menu.can_read,
		role_menu.can_write,
		role_menu.can_delete,
		role_menu.can_run,
		role_menu.can_approve,
		role_menu.can_cancel
	FROM user
	LEFT JOIN user_role ON user.id = user_role.user_id
	LEFT JOIN role ON user_role.role_id = role.id
	LEFT JOIN role_menu ON role.id = role_menu.role_id
	LEFT JOIN menu ON role_menu.menu_id = menu.id
	LEFT JOIN org ON role.org_id = org.id
	WHERE user.id = ?
	AND menu.id <> ISNULL(menu.id)
	AND menu.deleted IS NULL
	`
	up := UserPermission{}
	perms := []*MenuPermission{}
	err := db.Select(&perms, sql, u.ID)
	if err != nil {
		return up, err
	}
	up.UserID = u.ID
	up.UserName = u.Name
	up.Permission = perms
	return up, nil
}
