package pay

import (
	"fmt"
	"github.com/google/uuid"
	"lihood/boot"
	"lihood/g"
	"lihood/internal/enum"
	"testing"
)

func TestClient_DoRequest(t *testing.T) {
	boot.MustSetup()
	client := PayerFactory(enum.CloudPay)
	fmt.Println("32FY" + uuid.New().String()[:8])
	resp, err := client.Pay("32FY"+uuid.New().String()[:8], 1, g.ChainCallback("asd"))
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", resp)
}

func TestClientQuery(t *testing.T) {
	boot.MustSetup()
	client := PayerFactory(enum.Alipay)
	ok, err := client.Query("103A402d4fcf")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", ok)
}
