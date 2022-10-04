package models

import "lihood/internal/enum"

type Product struct {
	ID          int64              `json:"id,omitempty" column:"id"`
	Name        string             `json:"name,omitempty" column:"name"`
	Description string             `json:"description,omitempty" column:"description"`
	Image       string             `json:"image,omitempty" column:"image"`
	Status      enum.ProductStatus `json:"status,omitempty" column:"status"`
	AuthorID    int64              `json:"author_id,omitempty" column:"author_id"`
	Classify    int64              `json:"classify,omitempty" column:"classify"`
	Price       int64              `json:"price,omitempty" column:"price"`
	Stock       int64              `json:"stock,omitempty" column:"stock"`
	Index       int                `json:"-" column:"index"`
	OrderNo     string             `json:"-" column:"order_no"`
	TxID        string             `json:"tx_id,omitempty" column:"tx_id"`
	TokenID     string             `json:"token_id" column:"token_id"`
	PayType     int                `json:"pay_type" column:"pay_type"`
	Cname       string             `json:"cname" column:"c_name"`
	CreateTime  int64              `json:"create_time,omitempty" column:"create_time"`
	UpdateTime  int64              `json:"-" column:"update_time"`
	DelTime     int64              `json:"-" column:"del_time"`
	SaleTime    int64              `json:"sale_time" column:"sale_time"`
	AdvanceHour int                `json:"advance_hour" column:"advance_hour"`
}

func (p Product) TableName() string {
	return "tbl_product"
}
