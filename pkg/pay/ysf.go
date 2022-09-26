package pay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"lihood/conf"
	"net/http"
	"time"
)

func NewCloudPay() Payer {
	return cloudPay{}
}

type cloudPay struct{}

func (a cloudPay) Pay(orderID string, amount int64, cb string) (interface{}, error) {
	var item = struct {
		RequestTimestamp string `json:"requestTimestamp"`
		MerOrderId       string `json:"merOrderId"`
		Mid              string `json:"mid"`
		Tid              string `json:"tid"`
		InstMid          string `json:"instMid"`
		TotalAmount      int64  `json:"totalAmount"`
		NotifyUrl        string `json:"notifyUrl"`
	}{
		RequestTimestamp: time.Now().Format("2006-01-02 15:04:05"),
		MerOrderId:       orderID,
		Mid:              conf.Instance.Pay.MID,
		Tid:              conf.Instance.Pay.TID,
		InstMid:          "APPDEFAULT",
		TotalAmount:      amount,
		NotifyUrl:        cb,
	}
	fmt.Println(item.NotifyUrl)
	data, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(data))
	token := signFromRequest(data, conf.Instance.Pay.AppID, conf.Instance.Pay.AppSecret)
	req, err := http.NewRequest(http.MethodPost, "https://api-mop.chinaums.com/v1/netpay/uac/app-order", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result interface{}
	if err = json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (a cloudPay) Query(orderID string) (bool, error) {
	var item = struct {
		RequestTimestamp string `json:"requestTimestamp"`
		MerOrderId       string `json:"merOrderId"`
		Mid              string `json:"mid"`
		Tid              string `json:"tid"`
		InstMid          string `json:"instMid"`
		MsgId            string `json:"msgId"`
	}{
		RequestTimestamp: time.Now().Format("2006-01-02 15:04:05"),
		MerOrderId:       orderID,
		Mid:              conf.Instance.Pay.MID,
		Tid:              conf.Instance.Pay.TID,
		InstMid:          "APPDEFAULT",
		MsgId:            time.Now().Format("2006-01-02 15:04:05"),
	}
	data, err := json.Marshal(item)
	if err != nil {
		return false, err
	}
	token := signFromRequest(data, conf.Instance.Pay.AppID, conf.Instance.Pay.AppSecret)
	req, err := http.NewRequest(http.MethodPost, conf.Instance.Pay.QueryURL, bytes.NewBuffer(data))
	if err != nil {
		return false, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	var result struct {
		Status string `json:"status"`
	}
	if err = json.Unmarshal(data, &result); err != nil {
		return false, err
	}
	return result.Status == "TRADE_SUCCESS", nil
}
