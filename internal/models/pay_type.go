package models

type PayType struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status int    `json:"status"`
}

func (p PayType) TableName() string {
	return "tbl_pay_type"
}
