package utils

import (
	"fmt"
	"testing"
)

func TestCheckIdCard(t *testing.T) {
	fmt.Println(CheckIdCard("340824199611032233", "张三", "13888888888"))
}

func TestSendMsCode(t *testing.T) {
	fmt.Println(SendMsCode("1234", "15055461510"))
}
