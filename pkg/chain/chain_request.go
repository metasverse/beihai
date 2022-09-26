package chain

type AccountRequest struct {
	AppId        string `json:"appId"`
	Name         string `json:"name"`
	Pwd          string `json:"pwd"`
	Timestamp    string `json:"timestamp"`
	CurrencyType int    `json:"currency_type"`
	Sign         string `json:"sign"`
}

type AccountResponse struct {
	Success bool   `json:"success"`
	ErrCode int    `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
	Data    struct {
		CurrencyType int    `json:"currency_type"`
		Address      string `json:"address"`
	} `json:"data"`
}

type ProductRequest struct {
	AppId         string `json:"appId"`
	OpenID        string `json:"openid"`
	Pid           string `json:"pid"`
	Name          string `json:"name"`
	ImageUrl      string `json:"imageUrl"`
	Description   string `json:"description"`
	WalletAddress string `json:"wallet_address"`
	NotifyUrl     string `json:"notify_url"`
	Timestamp     string `json:"timestamp"`
	Sign          string `json:"sign"`
}

type ProductResponse struct {
	Success bool   `json:"success"`
	ErrCode int    `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
	Data    struct {
		TxId     string `json:"tx_id"`
		Hash     string `json:"hash"`
		TokenId  string `json:"token_id"`
		MetaData string `json:"meta_data"`
		Status   int    `json:"status"`
	} `json:"data"`
}

type PreOrderRequest struct {
	AppID       string `json:"appId"`
	Pid         string `json:"pid"`
	TokenId     string `json:"token_id"`
	FromAddress string `json:"from_address"`
	ToAddress   string `json:"to_address"`
	NotifyUrl   string `json:"notify_url"`
	Timestamp   string `json:"timestamp"`
	Sign        string `json:"sign"`
}

type PreOrderResponse struct {
	Success bool   `json:"success"`
	ErrCode int    `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
	Data    struct {
		PrepayID string `json:"prepay_id"`
	} `json:"data"`
}

type TransferRequest struct {
	AppID     string `json:"appId"`
	PrepayId  string `json:"prepay_id"`
	Timestamp string `json:"timestamp"`
	Sign      string `json:"sign"`
}

type TransferResponse struct {
	Success bool   `json:"success"`
	ErrCode int    `json:"err_code"`
	Errmsg  string `json:"err_msg"`
}
