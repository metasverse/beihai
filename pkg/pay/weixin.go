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

type weixin struct{}

func (w weixin) Pay(orderID string, amount int64, cb string) (interface{}, error) {
	var item = struct {
		RequestTimestamp string `json:"requestTimestamp"`
		MerOrderId       string `json:"merOrderId"`
		Mid              string `json:"mid"`
		Tid              string `json:"tid"`
		InstMid          string `json:"instMid"`
		TotalAmount      int64  `json:"totalAmount"`
		SubAppId         string `json:"subAppId"`
	}{
		RequestTimestamp: time.Now().Format("2006-01-02 15:04:05"),
		MerOrderId:       orderID,
		Mid:              "89844035732APSD",
		Tid:              conf.Instance.Pay.TID,
		InstMid:          "APPDEFAULT",
		TotalAmount:      amount,
	}
	data, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}
	fmt.Println(conf.Instance.Pay.PayURL)
	token := signFromRequest(data, "8a81c1be818ca2270181f6c390a204c5", "c8a2e2f190414d65a0b2fdce38b90f04")
	req, err := http.NewRequest(http.MethodPost, "https://api-mop.chinaums.com/v1/netpay/wx/app-pre-order", bytes.NewBuffer(data))
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

func (w weixin) Query(orderID string) (bool, error) {
	//TODO implement me
	panic("implement me")
}
