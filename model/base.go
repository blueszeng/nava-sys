package model

import (
	"database/sql"
	"encoding/json"
	"github.com/go-sql-driver/mysql"
	"log"
	"time"
)

type Status int

const (
	DRAFT Status = iota // ฉบับร่าง ค่า int = 0
	OPEN // เอกสารบันทึกเข้าระบบงานแล้ว ยังสามารถแก้ไขได้
	HOLD // เอกสารถูกพักรอดำเนินการ
	POST // เอกสารถูกบันทึกบัญชีแล้ว ห้ามแก้ไข
	CANCEL // เอกสารถูกยกเลิกแล้ว ห้ามแก้ไข
)

// Base structure contains fields that are common to objects
// returned by the nava's REST API.
type Base struct {
	ID      uint64       `json:"id"`
	Created JsonNullTime `json:"-"`
	Updated JsonNullTime `json:"-"`
	Deleted JsonNullTime `json:"-"`
}

type Doc struct {
	CreatedBy  User
	UpdatedBy  User
	ApprovedBy User
	CanceledBy User
	DeletedBy  User
}

type JsonNullTime struct {
	mysql.NullTime
}

func (v JsonNullTime) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Time)
	} else {
		return json.Marshal(nil)
	}
}

func (v JsonNullTime) UnmarshalJSON(data []byte) error {
	var x *time.Time
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.Time = *x
	} else {
		v.Valid = false
	}
	return nil
}

type JsonNullInt64 struct {
	sql.NullInt64
}

func (v JsonNullInt64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int64)
	} else {
		return json.Marshal(nil)
	}
}

func (v *JsonNullInt64) UnmarshalJSON(data []byte) error {
	// Unmarshalling into a pointer will let us detect null
	var x *int8
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.Int64 = int64(*x)
	} else {
		v.Valid = false
	}
	return nil
}

type JsonNullDate struct {
	mysql.NullTime
}

func (v JsonNullDate) UnmarshalJSON(data []byte) error {
	//const layout = "02/01/2006"
	var err error
	if data != nil {
		v.Time, err = time.Parse(time.RFC3339, string(data))
		if err != nil {
			return err
		}
		v.Valid = true
		log.Println("data = ", string(data), "v.Time = ", v.Time)
	} else {
		v.Valid = false
	}
	return nil
}
