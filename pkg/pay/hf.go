package pay

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
)

// 加密参数
type PfxData struct {
	PfxFileName string `json:"pfx_file_name"`
	PfxFilePwd  string `json:"pfx_file_pwd"`
}

//加签请求参数
type MakeSignFrom struct {
	Data   string `json:"data"`
	Params string `json:"params"`
}

//解签请求参数
type VerifySignFrom struct {
	CheckValue string `json:"check_value" comment:"加密参数"`
	CertFile   string `json:"cert_file" comment:"解密数据证书服务器路径"`
}

//加签返回加密参数
type MakeSignData struct {
	CheckValue string `json:"check_value" comment:"加密参数"`
}

type HfPayData struct {
	MerCustId   string `json:"mer_cust_id" comment:"商户号"`
	Version     string `json:"version" comment:"版本号"`
	OrderDate   string `json:"order_date" comment:"订单日期"`
	OrderId     string `json:"order_id" comment:"订单编号，唯一"`
	TransAmt    string `json:"trans_amt" comment:"订单金额，最多两位小数"`
	GoodsDesc   string `json:"goods_desc" comment:"商品描述"`
	UserName    string `json:"user_name" comment:"真实姓名"`
	IdCardType  string `json:"id_card_type" comment:"证件类型"`
	IdCard      string `json:"id_card" comment:"身份证号"`
	InCustId    string `json:"in_cust_id" comment:"入账客户号"`
	RetUrl      string `json:"ret_url" comment:"支付成功，前端转跳页面"`
	BgRetUrl    string `json:"bg_ret_url" comment:"支付异步回调链接"`
	DevInfoJson string `json:"dev_info_json"`
	ObjectInfo  string `json:"object_info"`
	CheckValue  string `json:"check_value" comment:"加密参数"`
}

var Pfx = PfxData{PfxFileName: "/root/tomcat/apache-tomcat-8.5.64/webapps/HF0331.pfx", PfxFilePwd: "123456"}

func HfPay(desc, idcard, name, amount, retURL, typ string) ([]byte, error) {
	//统一下单参数
	var hfPayData = HfPayData{
		MerCustId:   "6666000100035059",
		Version:     "10",
		OrderDate:   time.Now().Format("20060102"),
		OrderId:     GenerateOrderId(1),
		TransAmt:    amount,
		GoodsDesc:   desc,
		UserName:    name,
		IdCardType:  "10",
		IdCard:      idcard,
		InCustId:    "6666000100035059",
		RetUrl:      "http://124.223.104.243/test.html?type=" + typ,
		BgRetUrl:    retURL, //回调post接收到的是个加密字符串，用解签方法进行解密获取内容
		DevInfoJson: `{"devType": "1","ipAddr": "124.223.104.243","IMEI": "011472001976595"}`,
		ObjectInfo:  `{"marketType":"1"}`,
	}

	//将下单参数加密
	data, _ := json.Marshal(&Pfx) //秘钥
	params, _ := json.Marshal(&hfPayData)
	payload := fmt.Sprintf("data=%s&params=%s", string(data), string(params))
	body, err := Send("http://124.223.104.243:8080/hfpcfca/cfca/makeSign", payload)
	if err != nil {
		return nil, err
	}
	var ret MakeSignData
	err = json.Unmarshal(body, &ret)
	if err != nil {
		return nil, err
	}
	//t.Log(ret.CheckValue)

	//开始调用下单接口
	payload = fmt.Sprintf("mer_cust_id=%s&version=%s&order_date=%s&order_id=%s&trans_amt=%s&goods_desc=%s&dev_info_json=%s&object_info=%s&user_name=%s&id_card_type=%s&id_card=%s&in_cust_id=%s&ret_url=%s&bg_ret_url=%s&check_value=%s",
		hfPayData.MerCustId,
		hfPayData.Version,
		hfPayData.OrderDate,
		hfPayData.OrderId,
		hfPayData.TransAmt,
		hfPayData.GoodsDesc,
		hfPayData.DevInfoJson,
		hfPayData.ObjectInfo,
		hfPayData.UserName,
		hfPayData.IdCardType,
		hfPayData.IdCard,
		hfPayData.InCustId,
		hfPayData.RetUrl,
		hfPayData.BgRetUrl,
		ret.CheckValue,
	)

	body, err = Send("https://hfpay.testpnr.com/api/hfpwallet/pay033", payload)
	if err != nil {
		return nil, err
	}
	//var ret MakeSignData
	err = json.Unmarshal(body, &ret)
	if err != nil {
		return nil, err
	}

	//解签，解密下单返回的参数
	var vs VerifySignFrom
	vs.CheckValue = ret.CheckValue
	vs.CertFile = "/root/tomcat/apache-tomcat-8.5.64/webapps/CFCA_ACS_TEST_CA.cer"
	p, _ := json.Marshal(&vs)
	payload = fmt.Sprintf("params=%s", string(p))

	return Send("http://124.223.104.243:8080/hfpcfca/cfca/verifySign", payload)
}

// GenerateOrderId 生成分布式订单编号
func GenerateOrderId(node int64) string {
	n, _ := snowflake.NewNode(node)
	s := n.Generate()
	return "BH" + s.String()
}

// Send 支付接口请求封装
func Send(url, payload string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return io.ReadAll(res.Body)
}
