package models

type AccountWallet struct {
	Kind       int    `json:"kind"`
	SourceType int    `json:"source_type"`
	ID         int64  `json:"id"`
	UID        int64  `json:"uid"`
	Amount     int64  `json:"amount"`
	SourceID   int64  `json:"source_id"`
	CreateTime int64  `json:"create_time"`
	Remark     string `json:"remark"`
}
