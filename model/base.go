package model

import (
	"time"
	"encoding/json"
	"database/sql"
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

type JsonNullString struct {
	sql.NullString
}

func (v JsonNullString) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.String)
	} else {
		return json.Marshal(nil)
	}
}

func (v JsonNullString) UnmarshalJSON(data []byte) error {
	var x *string
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.String = *x
	} else {
		v.Valid = false
	}
	return nil
}

