package pay

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"lihood/conf"
	"net/http"
	"strings"
	"time"
)

func NewAlipay() Payer {
	return &alipay{}
}

type alipay struct{}

func (a alipay) Pay(orderID string, amount int64, cb string) (interface{}, error) {
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
		NotifyUrl:        conf.Instance.Server.Domain + fmt.Sprintf("/api/v1/product/callback/%s", orderID),
	}
	data, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}
	token := signFromRequest(data, conf.Instance.Pay.AppID, conf.Instance.Pay.AppSecret)
	req, err := http.NewRequest(http.MethodPost, "https://test-api-open.chinaums.com/v1/netpay/trade/app-pre-order", bytes.NewBuffer(data))
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

func (a alipay) Query(orderID string) (bool, error) {
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
	fmt.Println(string(data))
	var result struct {
		Status string `json:"status"`
	}
	if err = json.Unmarshal(data, &result); err != nil {
		return false, err
	}
	return result.Status == "TRADE_SUCCESS", nil
}

func signFromRequest(data []byte, appID, appSecret string) string {
	sha256Writer := sha256.New()
	sha256Writer.Write(data)
	bodySign := fmt.Sprintf("%x", sha256Writer.Sum(nil))
	ts := time.Now().Format("20060102150405")
	nonce := strings.ReplaceAll(uuid.New().String(), "-", "")
	key := appID + ts + nonce + bodySign
	writer := hmac.New(sha256.New, []byte(appSecret))
	writer.Write([]byte(key))
	sign := base64.StdEncoding.EncodeToString(writer.Sum(nil))
	return fmt.Sprintf(`OPEN-BODY-SIG AppId="%s", Timestamp="%s", Nonce="%s", Signature="%s"`, appID, ts, nonce, sign)
}
