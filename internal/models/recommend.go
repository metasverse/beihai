package models

import "lihood/internal/enum"

type Recommend struct {
	ID         int64                `json:"id" column:"id"`
	ProductID  int64                `json:"product_id" column:"product_id"`
	Index      int                  `json:"index" column:"index"`
	Status     enum.RecommendStatus `json:"status" column:"status"`
	CreateTime int64                `json:"create_time" column:"create_time"`
}

func (r Recommend) TableName() string {
	return "tbl_recommend"
}
