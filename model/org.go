package model

import (
	"github.com/jmoiron/sqlx"
	"log"
	"github.com/mrtomyum/nava-sys/api"
	"github.com/guregu/null"
)

type Org struct {
	Base
	NameTh   null.String  `json:"name_th" db:"name_th"`
	NameEn   null.String  `json:"name_en" db:"name_en"`
	ParentId uint64  `json:"parent_id" db:"parent_id"`
	Roles    []*Role `json:"roles,omitempty"`
	Related api.Link `json:"related,omitempty"`
	//Permission []*MenuPermission
}

func (o *Org) All(db *sqlx.DB) (orgs []*Org, err error) {
	sql := `SELECT * FROM org`
	err = db.Select(&orgs, sql)
	if err != nil {
		log.Println("Error db.Select(&orgs)", err.Error())
		return nil, err
	}
	for _, org := range orgs {
		sql = `
			SELECT * FROM role
			WHERE org_id = ?
		`
		var roles []*Role
		err = db.Select(&roles, sql, org.ID)
		if err != nil {
			log.Println("Error db.Select(&role)", err.Error())
			return nil, err
		}
		for _, r := range roles {
			org.Roles = append(org.Roles, r)
		}
	}

	return orgs, nil
}
