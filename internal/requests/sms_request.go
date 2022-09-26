package requests

type SMSSendRequest struct {
	Phone string `json:"phone" validate:"phone(m=请输入合法的手机号码)"`
}

type SMSLoginRequest struct {
	Phone      string `json:"phone" validate:"phone(m=请输入合法的手机号码)"`
	Code       string `json:"code" form:"code" xml:"code"`
	Invitation string `json:"invitation"`
}
