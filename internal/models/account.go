package models

import "lihood/internal/enum"

type Account struct {
	ID                     int64              `json:"id" column:"id"`
	UID                    string             `json:"uid" column:"uid"`
	PID                    int64              `json:"pid" column:"pid"`
	Nickname               string             `json:"nickname" column:"nickname"`
	Avatar                 string             `json:"avatar" column:"avatar"`
	Status                 enum.AccountStatus `json:"status" column:"status"`
	Phone                  string             `json:"phone" column:"phone"`
	BsnAddress             string             `json:"bsn_address" column:"bsn_address"`
	BsnUsername            string             `json:"bsn_username" column:"bsn_username"`
	Name                   string             `json:"name" column:"name"`
	IDCardNum              string             `json:"id_card_num" column:"id_card_num"`
	IdCardPositiveImageUrl string             `json:"id_card_positive_image_url" column:"id_card_positive_image_url"`
	IdCardNegativeImageUrl string             `json:"id_card_negative_image_url" column:"id_card_negative_image_url"`
	Role                   enum.Role          `json:"role" column:"role"`     // 角色
	Amount                 int64              `json:"amount" column:"amount"` // 账户余额
	CreateTime             int64              `json:"create_time"`
	UpdateTime             int64              `json:"update_time"`
	DelTime                int64              `json:"del_time"`
	Description            string             `json:"description" column:"description"`
}

func (a Account) TableName() string {
	return "tbl_account"
}
