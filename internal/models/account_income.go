package models

import "lihood/internal/enum"

type AccountIncome struct {
	ID         int64           `json:"id"`
	UID        int64           `json:"uid"`
	Type       enum.IncomeType `json:"type"`
	Amount     int64           `json:"amount"`
	Remark     string          `json:"remark"`
	CreateTime int64           `json:"create_time"`
}

func (i AccountIncome) TableName() string {
	return "tbl_account_income"
}
