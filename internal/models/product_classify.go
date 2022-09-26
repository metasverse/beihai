package models

type ProductClassify struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	CreateTime int64  `json:"create_time"`
	UpdateTime int64  `json:"update_time"`
	DelTime    int64  `json:"del_time"`
}

func (p ProductClassify) TableName() string {
	return "tbl_product_classify"
}
