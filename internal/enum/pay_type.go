package enum

type PayType int

const (
	Alipay   PayType = iota // 支付宝
	Wechat                  // 微信
	CloudPay                // 云闪付
)
