package chain

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

type BSNChain struct {
	APPId     string
	AppSecret string
	OpenID    string
}

func NewChainClient() *BSNChain {
	return &BSNChain{
		APPId:     "mssc002",
		AppSecret: "341b0713945818bcec4b5a06ff341b2c",
		OpenID:    "XfpEKZKX8y24o2nMfQg26bg4ib1dy",
	}
}

func (c BSNChain) Sign(v interface{}) (sign string, err error) {
	// 已经排序好的字段
	rt := reflect.TypeOf(v).Elem()
	value := reflect.Indirect(reflect.ValueOf(v))
	var keys []string
	valuesMap := make(map[string]string, 0)
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if field.Name == "Sign" {
			continue
		}
		jsonTagName := field.Tag.Get("json")
		keys = append(keys, jsonTagName)
		valuesMap[jsonTagName] = fmt.Sprintf("%v", value.Field(i).Interface())
	}
	sort.Strings(keys)
	var values []string
	for _, key := range keys {
		values = append(values, fmt.Sprintf("%s%s", key, valuesMap[key]))
	}
	text := strings.Join(values, "") + c.AppSecret
	writer := md5.New()
	if _, err := io.WriteString(writer, text); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", writer.Sum(nil)), nil
}

func (c BSNChain) NewAccount(name string) (*AccountResponse, error) {
	req := &AccountRequest{
		AppId:        c.APPId,
		CurrencyType: 3,
		Name:         name,
		Pwd:          "123456",
		Timestamp:    strconv.FormatInt(time.Now().Unix(), 10),
	}
	sign, err := c.Sign(req)
	if err != nil {
		return nil, err
	}
	req.Sign = sign
	var buffer = new(bytes.Buffer)
	if err := json.NewEncoder(buffer).Encode(req); err != nil {
		return nil, err
	}
	resp, err := http.Post("https://api.coltstail.net/v1/ms/wallet/create", "application/json", buffer)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var r AccountResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (c BSNChain) NewProduct(pid, image, wallet, desc, callback string) (*ProductResponse, error) {
	req := &ProductRequest{
		AppId:         c.APPId,
		Pid:           pid,
		OpenID:        "a",
		Name:          uuid.New().String(),
		ImageUrl:      image,
		Description:   desc,
		WalletAddress: wallet,
		NotifyUrl:     callback,
		Timestamp:     strconv.FormatInt(time.Now().Unix(), 10),
	}
	sign, err := c.Sign(req)
	if err != nil {
		return nil, err
	}
	req.Sign = sign
	var buffer = new(bytes.Buffer)
	if err := json.NewEncoder(buffer).Encode(req); err != nil {
		return nil, err
	}
	resp, err := http.Post("https://api.coltstail.net/v1/token/create", "application/json", buffer)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result ProductResponse
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c BSNChain) PreOrder(pid, tokenID, fromAddr, toAddr, url string) (*PreOrderResponse, error) {
	req := &PreOrderRequest{
		AppID:       c.APPId,
		Pid:         pid,
		TokenId:     tokenID,
		FromAddress: fromAddr,
		ToAddress:   toAddr,
		NotifyUrl:   url,
		Timestamp:   strconv.FormatInt(time.Now().Unix(), 10),
	}
	fmt.Printf("%+v", req)
	sign, err := c.Sign(req)
	if err != nil {
		return nil, err
	}
	req.Sign = sign
	var buffer = new(bytes.Buffer)
	if err := json.NewEncoder(buffer).Encode(req); err != nil {
		return nil, err
	}
	resp, err := http.Post("https://api.coltstail.net/v1/trade/pre_order", "application/json", buffer)
	if err != nil {
		return nil, err
	}
	fmt.Println(buffer.String())
	defer resp.Body.Close()
	var result PreOrderResponse
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c BSNChain) Transfer(prepayId string) (*TransferResponse, error) {
	req := &TransferRequest{
		AppID:     c.APPId,
		PrepayId:  prepayId,
		Timestamp: strconv.FormatInt(time.Now().Unix(), 10),
	}
	sign, err := c.Sign(req)
	if err != nil {
		return nil, err
	}
	req.Sign = sign
	var buffer = new(bytes.Buffer)
	if err := json.NewEncoder(buffer).Encode(req); err != nil {
		return nil, err
	}
	resp, err := http.Post("https://api.coltstail.net/v1/token/transfer", "application/json", buffer)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result TransferResponse
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

type Callback struct {
	Appid           string `json:"appid"`
	Openid          string `json:"openid"`
	Pid             string `json:"pid"`
	TxId            string `json:"tx_id"`
	TokenId         string `json:"token_id"`
	ContractAddress string `json:"contract_address"`
	WalletAddress   string `json:"wallet_address"`
	Hash            string `json:"hash"`
}
