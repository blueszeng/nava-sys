package models

import (
	"time"
)

// Base structure contains fields that are common to objects
// returned by the nava's REST API.
type Base struct {
	ID      uint64         `json:"id"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	Deleted time.Time `json:"deleted"`
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


