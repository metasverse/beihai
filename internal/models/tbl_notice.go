package models

type Notice struct {
	ID         int64  `json:"id" column:"id"`
	NID        int64  `json:"nid" column:"nid"`
	Title      string `json:"title" column:"title"`
	Summary    string `json:"summary" column:"summary"`
	Content    string `json:"content" column:"content"`
	Display    bool   `json:"display" column:"display"`
	CreateTime int64  `json:"create_time" column:"create_time"`
}

func (Notice) TableName() string {
	return "tbl_notice"
}
