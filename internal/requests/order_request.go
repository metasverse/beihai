package requests

type OrderRequest struct {
	PID     int64 `json:"pid"`
	PayType int   `json:"pay_type" validate:"contains(m=非法的支付方式,v=0,1,2)"`
}
