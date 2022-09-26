package models

import "lihood/internal/enum"

type AccountBank struct {
	ID         int64                  `json:"id" column:"id"`
	UID        int64                  `json:"uid" column:"uid"`
	BankID     int64                  `json:"bank_id" column:"bank_id"`
	Name       string                 `json:"name" column:"name"`
	BankName   string                 `json:"bank_name" column:"bank_name"`
	BankNum    string                 `json:"bank_num" column:"bank_num"`
	Status     enum.AccountBankStatus `json:"status" column:"status"`
	CreateTime int64                  `json:"create_time" column:"create_time"`
	UpdateTime int64                  `json:"update_time" column:"update_time"`
	DelTime    int64                  `json:"del_time" column:"del_time"`
}

func (a AccountBank) TableName() string {
	return "tbl_account_bank"
}
