package entity

type RecommendList struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	TokenID    string `json:"token_id"`
	Image      string `json:"image"`
	AuthorID   int    `json:"author_id"`
	AuthorName string `json:"author_name"`
	Price      int    `json:"price"`
	Likes      int    `json:"likes"`
}
