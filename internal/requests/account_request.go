package requests

type UpdatePhoneRequest struct {
	Phone string `json:"phone" validate:"required(m=手机号不能为空);phone(m=请输入正确的手机号)"`
	Code  string `json:"code" validate:"required(m=验证码不能为空)"`
}

type UpdatePhoneCodeRequest struct {
	Phone string `json:"phone" validate:"required(m=手机号不能为空);phone(m=请输入正确的手机号)"`
}

type AuthenticationRequest struct {
	Name          string `json:"name" validate:"required(m=用户名不能为空)"`
	IDCard        string `json:"id_card" validate:"required(m=身份证号不能为空)"`
	PositiveImage string `json:"positive_image"`
	NegativeImage string `json:"negative_image"`
}

type UpdateAccountInfoRequest struct {
	Nickname string `json:"nickname" validate:"required(m=昵称不能为空)"`
	Avatar   string `json:"avatar" validate:"required(m=头像不能为空);url(m=请上传正确的头像)"`
	Desc     string `json:"desc" validate:"required(m=个人简介不能为空)"`
}
