package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"

	"lihood/g"
	"lihood/internal/models"
	"lihood/internal/repository"
	"lihood/pkg/chain"
	"lihood/pkg/jwt"
	"lihood/utils"
)

type PhoneLoginService interface {
	LoginWithPhone(phone string, code string, invitation string) (string, error)
	Send(phone string) error
	PhoneLogin(phone, invitation string) (string, error)
}

func NewPhoneLoginService() PhoneLoginService {
	return &phoneLoginService{}
}

type phoneLoginService struct{}

func (l phoneLoginService) PhoneLogin(phone, invitation string) (string, error) {
	// 先去account账户中查询是否存在该用户
	dao := repository.NewAccountRepository(g.DB)

	// 判断用户是否存在
	user, err := dao.GetByPhone(phone)

	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			return "", err
		}
	}

	exist := user != nil

	if !exist {
		user = &models.Account{
			Nickname:   utils.HidePhone(phone),
			Avatar:     "https://lihood-1306623008.cos.ap-nanjing.myqcloud.com/3ea6beec64369c2642b92c6726f1epng.png",
			Phone:      phone,
			CreateTime: time.Now().Unix(),
			UID:        utils.GetRandomString(8), // 设置唯一的邀请码字符串
		}

		// 获取用户的邀请码
		for {
			user.UID = utils.GetRandomString(8)
			if _, err = dao.GetByInviteCode(user.UID); err != nil {
				if err == sql.ErrNoRows {
					break
				} else {
					return "", err
				}
			}
		}
		if invitation != "" {
			parent, err := dao.GetByInviteCode(invitation)
			if err != nil && err != sql.ErrNoRows {
				return "", err
			}
			if parent != nil {
				user.PID = parent.ID
			}
		}
		client := chain.NewChainClient()
		username := uuid.New().String()
		resp, err := client.NewAccount(username)
		if err != nil {
			return "", err
		}
		if !resp.Success {
			return "", errors.New(resp.ErrMsg)
		}
		user.BsnAddress = resp.Data.Address
		user.BsnUsername = username
		if err = dao.Create(user); err != nil {
			return "", err
		}
	}
	return jwt.GenToken(user.ID)
}

func (l phoneLoginService) Send(phone string) error {
	// 先去从redis取出最新没有消费的验证码
	key := fmt.Sprintf("sms:login:%s", phone)
	ctx := context.Background()
	var code string
	if cmd := g.Redis.Get(ctx, key); cmd.Err() != nil {
		if cmd.Err() == redis.Nil {
			// 没有验证码，生成一个4位验证码，存入redis，并发送短信
			code = fmt.Sprintf("%04d", rand.Intn(9999))
		} else {
			return cmd.Err()
		}
	} else {
		code = cmd.Val()
	}
	// 存入redis
	if _, err := g.Redis.Set(ctx, key, code, time.Minute*10).Result(); err != nil {
		return err
	}
	return NewSmsService().Send(phone, code)
}

func (l phoneLoginService) LoginWithPhone(phone string, code string, invitation string) (string, error) {
	// 开始限流
	// 从redis里面判断当前的手机号码的验证次数
	checkKey := fmt.Sprintf("sms:login:check:%s", phone)

	ctx := context.Background()

	if err := g.Redis.SetNX(ctx, checkKey, "0", time.Minute*3).Err(); err != nil {
		return "", err
	}

	times, err := g.Redis.Incr(ctx, checkKey).Result()
	if err != nil {
		return "", err
	}

	// 限流，防止恶意攻击
	if times > 5 {
		return "", g.Error("您的验证码已经超过5次验证错误，请稍后再试")
	}

	// 先去从redis取出最新没有消费的验证码
	key := fmt.Sprintf("sms:login:%s", phone)

	var realCode string
	if cmd := g.Redis.Get(ctx, key); cmd.Err() != nil {
		if cmd.Err() == redis.Nil {
			return "", g.Error("验证码错误或已过期")
		} else {
			return "", cmd.Err()
		}
	} else {
		realCode = cmd.Val()
	}
	// 比较验证码
	if code != realCode {
		return "", g.Error("验证码错误或已过期")
	}
	// 删除redis中的checkKey
	if _, err := g.Redis.Del(ctx, checkKey).Result(); err != nil {
		return "", err
	}
	return l.PhoneLogin(phone, invitation)
}
