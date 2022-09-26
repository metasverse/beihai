package models

type UserProduct struct {
	ID         int64  `json:"id,omitempty" column:"id"`
	PID        int64  `json:"pid" column:"pid"`
	UID        int64  `json:"uid" column:"uid"`
	Times      int64  `json:"times" column:"times"`
	TxID       string `json:"tx_id" column:"tx_id"`
	TokenID    string `json:"token_id" column:"token_id"`
	Hash       string `json:"hash" column:"hash"`
	CID        string `json:"cid" column:"cid"`
	Status     bool   `json:"status" column:"status"`           // 是否上链
	Display    bool   `json:"display" column:"display"`         // 是否展示
	IsAirDrop  bool   `json:"is_air_drop" column:"is_air_drop"` // 是否为空投
	Reason     string `json:"reason" column:"reason"`
	CName      string `json:"c_name" column:"c_name"`
	CreateTime int64  `json:"create_time" column:"create_time"`
	SaleTime   int64  `json:"sale_time" column:"sale_time"`
}

func (p UserProduct) TableName() string {
	return "tbl_user_product"
}
