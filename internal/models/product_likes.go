package models

type ProductLikes struct {
	ID         int64 `json:"id"`
	UID        int64 `json:"uid"`
	PID        int64 `json:"pid"`
	CreateTime int64 `json:"create_time"`
}

func (p ProductLikes) TableName() string {
	return "tbl_product_likes"
}
