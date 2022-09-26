package enum

type ProductStatus int

const (
	ProductUnpaid ProductStatus = iota // 未支付
	ProductPaid                        // 已支付
)

type ChainStatus int

const (
	ChainUnpaid ChainStatus = iota // 未上链
	ChainPaid                      // 已上链
)
