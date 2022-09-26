package requests

type BankCreateRequest struct {
	Name     string `json:"name" validate:"required(m=银行名称不能为空)"`
	BankID   int64  `json:"bank_id" validate:"required(m=银行ID不能为空)"`
	BankName string `json:"bank_name" validate:"required(m=银行支行不能为空)"`
	BankNum  string `json:"bank_num" validate:"required(m=银行卡号不能为空)"`
}
