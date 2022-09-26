package entity

type AuthorEntity struct {
	ID         int64  `json:"id" column:"id"`
	Nickname   string `json:"nickname" column:"nickname"`
	Avatar     string `json:"avatar" column:"avatar"`
	WorksCount int64  `json:"works_count" column:"works_count"`
}
