package services

import (
	"context"
	"fmt"
	"lihood/g"
	"math/rand"
	"time"
)

type SmsService interface {
	Send(phone string, code string) error
	SendWithDrawMessage(phone string) error
	CheckWithDrawCode(phone, code string) error
}

func NewSmsService() SmsService {
	return &smsService{}
}

type smsService struct{}

func (s smsService) CheckWithDrawCode(phone, code string) error {
	ctx := context.Background()
	key := fmt.Sprintf("withdraw:%s", phone)
	result, err := g.Redis.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	if result != code {
		return g.Error("验证码错误")
	}
	return nil
}

func (s smsService) SendWithDrawMessage(phone string) error {
	// 先限流
	ctx := context.Background()
	// 生成4为验证码
	code := fmt.Sprintf("%04d", rand.Intn(9999))
	code = "1234"
	// 存入redis
	if _, err := g.Redis.Set(ctx, fmt.Sprintf("withdraw:%s", phone), code, time.Minute*10).Result(); err != nil {
		return err
	}
	// 发送短信
	return s.Send(phone, code)
}

func (s smsService) Send(phone string, code string) error {
	return nil
}
