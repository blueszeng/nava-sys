package model

type org struct {
	Base
	NameTh string `json:"name_th" db:"name_th"`
	NameEn string `json:"name_en" db:"name_en"`
	ParentId uint64 `json:"parent_id" db:"parent_id"`

}