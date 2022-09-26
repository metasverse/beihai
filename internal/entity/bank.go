package entity

import "lihood/internal/enum"

type AccountBank struct {
	ID             int64                  `json:"id"`
	UID            int64                  `json:"uid"`
	BankID         int64                  `json:"bank_id"`
	Name           string                 `json:"name"`
	BankName       string                 `json:"bank_name"`
	BankNum        string                 `json:"bank_num"`
	Status         enum.AccountBankStatus `json:"status"`
	CreateTime     int64                  `json:"create_time"`
	UpdateTime     int64                  `json:"update_time"`
	DelTime        int64                  `json:"del_time"`
	BankOfficeName string                 `json:"bank_office_name"`
}
