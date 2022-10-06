package pay

import (
	"testing"
)

func TestHfPay(t *testing.T) {
	data, err := HfPay("商品", "421122198812084932", "黄勇", "10.00", "", "1")
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(string(data))
	}
}
