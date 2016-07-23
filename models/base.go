package models

import (
	"github.com/go-sql-driver/mysql"
)

// Base structure contains fields that are common to objects
// returned by the nava's REST API.
type Base struct {
	ID      uint64         `json:"id"`
	Created mysql.NullTime `json:"created"` //todo: change datatype to sql.NullTime
	Updated mysql.NullTime `json:"updated"`
	Deleted mysql.NullTime `json:"deleted"`
}

//func (b *Base) Delete(db *sql.DB) error {
//	return nil
//}

// If record deleted Unique field can not be duplicated
// Delete  bool           `json:"deleted"`
// Status Status 	`json:"status"`


type Status int
const (
	ACTIVE Status = 1 + iota
	HOLD
	SUSPEND
)


