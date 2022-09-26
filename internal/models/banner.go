package models

type Banner struct {
	Index      int    `json:"index" column:"index"`
	ID         int64  `json:"id" column:"id"`
	CreateTime int64  `json:"create_time" column:"create_time"`
	UpdateTime int64  `json:"update_time" column:"update_time"`
	DelTime    int64  `json:"del_time" column:"del_time"`
	ProductID  int64  `json:"product_id" column:"product_id"`
	Image      string `json:"image" column:"image"`
	Name       string `json:"name" column:"name"`
	Link       string `json:"link" column:"link"`
	Status     uint   `json:"status" column:"status"`
	BannerType uint   `json:"banner_type" column:"banner_type"`
}

func (b Banner) TableName() string {
	return "tbl_banner"
}
