package models

type Bank struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (b Bank) TableName() string {
	return "tbl_bank"
}
