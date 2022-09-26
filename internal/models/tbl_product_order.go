package models

import "lihood/internal/enum"

type ProductOrder struct {
	ID         int64              `json:"id"`
	OID        string             `json:"oid"`
	PayType    enum.PayType       `json:"pay_type"`
	UID        int64              `json:"uid"`
	PID        int64              `json:"pid"`
	Status     enum.ProductStatus `json:"status"`
	CreateTime int64              `json:"create_time"`
	PayTime    int64              `json:"pay_time"`
}

func (o ProductOrder) TableName() string {
	return "tbl_product_order"
}
