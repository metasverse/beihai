package enum

type OrderStatus int

const (
	OrderStatusUnpaid OrderStatus = iota // 未支付
	OrderStatusPaid                      // 已支付
)
