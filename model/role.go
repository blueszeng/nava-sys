package model

import (
	"github.com/jmoiron/sqlx"
	"log"
)

type Role struct {
	Base
	TH string `json:"th"`
	EN string `json:"en"`
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
	sql := `SELECT * FROM role`
	err := db.Select(&roles, sql)
	if err != nil {
		log.Println("Error in model.Role.All", err)
		return nil, err
	}
	return roles, nil
}
