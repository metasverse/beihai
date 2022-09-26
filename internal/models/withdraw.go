package models

type Withdraw struct {
	ID           int64 `json:"id"`
	UID          int64 `json:"uid"`
	BankID       int64 `json:"bank_id"`
	Amount       int64 `json:"amount"`
	Status       int   `json:"status"`
	CreateTime   int64 `json:"created_at"`
	WithdrawTime int64 `json:"updated_at"`
}

func (Withdraw) TableName() string {
	return "tbl_withdraw"
}
