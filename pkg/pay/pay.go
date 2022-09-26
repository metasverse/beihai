package pay

import "lihood/internal/enum"

type Payer interface {
	Pay(orderID string, amount int64, cb string) (interface{}, error)
	Query(orderID string) (bool, error)
}

func PayerFactory(payType enum.PayType) Payer {
	switch payType {
	case enum.Alipay:
		return NewAlipay()
	case enum.CloudPay:
		return NewCloudPay()
	default:
		return weixin{}
	}
}
