package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func CheckIdCard(idcard, name, mobile string) (bool, error) {
	const path = "http://sjsys.market.alicloudapi.com/communication/personal/1979"
	values := url.Values{}
	values.Add("idcard", idcard)
	values.Add("name", name)
	values.Add("mobile", mobile)
	reader := strings.NewReader(values.Encode())
	req, _ := http.NewRequest(http.MethodPost, path, reader)
	req.Header.Add("Authorization", "APPCODE "+"2bbefa357a14442db4ef157f9421e6c4")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	var item struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		return false, err
	}
	if item.Code != "10000" {
		return false, fmt.Errorf("check idcard failed: %s", item.Message)
	}
	return true, nil
}

func SendMsCode(code, mobile string) (bool, error) {
	const path = "https://dfsns.market.alicloudapi.com/data/send_sms"
	values := url.Values{}
	values.Add("content", fmt.Sprintf("code:%s", code))
	values.Add("phone_number", mobile)
	values.Add("template_id", "TPL_09828")
	reader := strings.NewReader(values.Encode())
	req, _ := http.NewRequest(http.MethodPost, path, reader)
	req.Header.Add("Authorization", "APPCODE "+"2bbefa357a14442db4ef157f9421e6c4")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	var item struct {
		RequestId string `json:"request_id"`
		Status    string `json:"status"`
		Reason    string `json:"reason"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		return false, err
	}
	if item.Status != "OK" {
		return false, fmt.Errorf("send mscode failed: %s", item.Reason)
	}
	return true, nil
}
