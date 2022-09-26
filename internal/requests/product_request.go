package requests

type ProductRequest struct {
	Name        string `json:"name" validate:"required(m=作品名称不能为空)"`
	Price       int64  `json:"price" validate:"required(m=作品价格不能为空);gt(m=作品价格必须大于0,v=0)"`
	Image       string `json:"image" validate:"url(m=请输入正确的图片地址)"`
	Count       uint64 `json:"count" validate:"required(m=作品数量不能为空);gt(m=作品数量必须大于0,v=0)"`
	Description string `json:"description" validate:"required(m=作品描述不能为空)"`
	PayType     int    `json:"pay_type"`
	SaleTime    uint8  `json:"sale_time"`
}

type PublicProductRequest struct {
	Uid         int64  `json:"uid"`
	Name        string `json:"name" validate:"required(m=作品名称不能为空)"`
	Price       int64  `json:"price" validate:"required(m=作品价格不能为空);gt(m=作品价格必须大于0,v=0)"`
	Image       string `json:"image" validate:"url(m=请输入正确的图片地址)"`
	Count       uint64 `json:"count" validate:"required(m=作品数量不能为空);gt(m=作品数量必须大于0,v=0)"`
	Description string `json:"description" validate:"required(m=作品描述不能为空)"`
	IsAirDrop   bool   `json:"is_air_drop"`
	Password    string `json:"password"`
}
