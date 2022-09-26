package chain

import (
	"fmt"
	"lihood/boot"
	"lihood/conf"
	"testing"
)

func TestBSNChain_NewAccount(t *testing.T) {
	boot.MustSetup()
	client := BSNChain{
		APPId:     "mssc002",
		AppSecret: "341b0713945818bcec4b5a06ff341b2c",
		OpenID:    "XfpEKZKX8y24o2nMfQg26bg4ib1dy",
	}
	resp, err := client.NewProduct("https://www.baidu.com/img/flexible/logo/pc/peak-result.png", "iaa1wnvq6yyp7x2a6ec6uc4h2wlyzvxmhr5alqg2hp", "desc", conf.Instance.Server.Domain+"/api/v1/product/callback")
	if err != nil {
		t.Error(err)
	} else {
		fmt.Printf("%+v\n", resp)
	}
}
