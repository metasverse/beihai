package entity

type ProductDetail struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Image        string `json:"image"`
	TokenID      string `json:"token_id"`
	AuthorID     int64  `json:"author_id"`
	AuthorName   string `json:"author_name"`
	AuthorAvatar string `json:"author_avatar"`
	AuthorDesc   string `json:"author_desc"`
	OwnerID      int64  `json:"owner_id"`
	OwnerName    string `json:"owner_name"`
	OwnerAvatar  string `json:"owner_avatar"`
	Price        int64  `json:"price"`
	Likes        int64  `json:"likes"`
	Sales        int64  `json:"sales"`
	Stock        int64  `json:"stock"`
	Description  string `json:"description"`
	Liked        bool   `json:"liked"`
	CreateTime   int64  `json:"create_time"`
	Hash         string `json:"hash"`
	CanBuy       bool   `json:"can_buy"`
}

type ProductList struct {
	ID           int    `json:"id" column:"id"`
	Name         string `json:"name" column:"name"`
	Image        string `json:"image" column:"image"`
	AuthorID     int64  `json:"author_id" column:"author_id"`
	AuthorName   string `json:"author_name" column:"author_name"`
	AuthorAvatar string `json:"author_avatar" column:"author_avatar"`
	Price        int64  `json:"price" column:"price"`
	Likes        int64  `json:"likes" colum:"likes"`
	Sales        int64  `json:"sales" column:"sales"`
	Liked        bool   `json:"liked" column:"is_liked"`
	IsAirDrop    bool   `json:"is_air_drop" column:"is_air_drop"`
	Times        int64  `json:"times" column:"times"`
	Countdown    int64  `json:"countdown" column:"countdown"`
	SaleTime     int64  `json:"sale_time" column:"sale_time"`
}

type SalesRank struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Image  string `json:"image"`
	Amount int64  `json:"amount"`
}
