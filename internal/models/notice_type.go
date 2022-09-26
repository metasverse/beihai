package models

type NoticeType struct {
	ID      int64  `json:"id" column:"id"`
	Name    string `json:"name" column:"name"`
	Display bool   `json:"display" column:"display"`
}

func (n NoticeType) TableName() string {
	return "tbl_notice_type"
}
